package application

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"Kube-CC/service"
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/uuid"
	"strconv"
)

var (
	n           = 10 // 使request为limit的1/10
	accessModes = "ReadWriteOnce"
)

// CreateAppDeploy 创建deploy类型的整个应用app
// 包括 configmap、pvc、deploy、service、TODO ingress
func CreateAppDeploy(form forms.DeployAddForm) (*responses.Response, error) {
	// 将form序列化为string，存入deploy的注释
	jsonBytes, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}
	strForm := string(jsonBytes)
	// 准备工作
	// 分割申请资源
	m := int(form.Replicas)
	requestCpu, err := service.SplitRSC(form.Cpu, n*m)
	if err != nil {
		return nil, err
	}
	requestMemory, err := service.SplitRSC(form.Memory, n*m)
	if err != nil {
		return nil, err
	}
	requestStorage, err := service.SplitRSC(form.Storage, n*m)
	if err != nil {
		return nil, err
	}
	limitsCpu, err := service.SplitRSC(form.Cpu, m)
	if err != nil {
		return nil, err
	}
	limitsMemory, err := service.SplitRSC(form.Memory, m)
	if err != nil {
		return nil, err
	}
	limitsStorage, err := service.SplitRSC(form.Storage, m)
	if err != nil {
		return nil, err
	}

	// 创建uuid，以便筛选出属于同一组的deploy、pod、service等
	newUuid := string(uuid.NewUUID())
	label := map[string]string{
		"uuid": newUuid,
	}

	// 创建PVC，持久存储
	var volumes []corev1.Volume
	var volumeMounts []corev1.VolumeMount
	if form.PvcStorage != "" {
		volumes = make([]corev1.Volume, 1)
		//volumeMounts = make([]corev1.VolumeMount, len(form.PvcPath))
		volumeMounts = make([]corev1.VolumeMount, 1)
		if form.StorageClassName == "" {
			return nil, errors.New("已填写PvcStorage,StorageClassName不能为空")
		}
		pvcName := form.Name + "-pvc"
		_, err = service.CreatePVC(form.Namespace, pvcName, form.StorageClassName, form.PvcStorage, accessModes)
		if err != nil {
			return nil, err
		}
		volumes[0] = corev1.Volume{
			Name: pvcName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvcName,
				},
			},
		}
		// 写死为/data目录
		volumeMounts[0] = corev1.VolumeMount{
			Name:      pvcName,
			MountPath: "/data",
		}
		//for i, path := range form.PvcPath {
		//	volumeMounts[i] = corev1.VolumeMount{
		//		Name:      pvcName,
		//		MountPath: path,
		//	}
		//}
	}

	// 0.创建configMap，存储环境变量
	configName := form.Name + "-configMap"
	if len(form.Env) > 0 {
		_, err = service.CreateConfigMap(configName, form.Namespace, label, form.Env)
		if err != nil {
			DeleteAppDeploy(form.Name, form.Namespace)
			return nil, err
		}
	}
	env := make([]corev1.EnvVar, len(form.Env))
	j := 0
	for k, _ := range form.Env {
		env[j] = corev1.EnvVar{
			Name: k,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: configName,
					},
					Key: k,
				},
			},
		}
		j++
	}

	// 1. 创建deployment
	// 端口
	num := len(form.Ports)
	ports := make([]corev1.ContainerPort, num)
	for i, port := range form.Ports {
		ports[i] = corev1.ContainerPort{
			ContainerPort: port,
		}
	}
	spec := appsv1.DeploymentSpec{
		Replicas: &form.Replicas,
		Selector: &metav1.LabelSelector{MatchLabels: label},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: label},
			Spec: corev1.PodSpec{
				Volumes:       volumes,
				RestartPolicy: corev1.RestartPolicyAlways,
				Containers: []corev1.Container{
					{
						Name:            form.Name,
						Image:           form.Image,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Command:         form.Command,
						Args:            form.Args,
						Ports:           ports,
						Env:             env,
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:              resource.MustParse(requestCpu),
								corev1.ResourceMemory:           resource.MustParse(requestMemory),
								corev1.ResourceEphemeralStorage: resource.MustParse(requestStorage),
								//TODO GPU
							},
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:              resource.MustParse(limitsCpu),
								corev1.ResourceMemory:           resource.MustParse(limitsMemory),
								corev1.ResourceEphemeralStorage: resource.MustParse(limitsStorage),
							},
						},
						VolumeMounts: volumeMounts,
					},
				},
			},
		},
	}
	_, err = service.CreateDeploy(form.Name, form.Namespace, strForm, label, spec)
	if err != nil {
		DeleteAppDeploy(form.Name, form.Namespace)
		return nil, err
	}

	// 2 创建service
	if num > 0 {
		servicePorts := make([]corev1.ServicePort, num)
		for i, port := range form.Ports {
			servicePorts[i] = corev1.ServicePort{
				Name:       strconv.Itoa(int(port)),
				Port:       port,
				TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: port},
			}
		}
		serviceName := form.Name + "-service"
		serviceSpec := corev1.ServiceSpec{
			Type:     corev1.ServiceTypeNodePort,
			Selector: label,
			Ports:    servicePorts,
		}
		_, err = service.CreateService(serviceName, form.Namespace, label, serviceSpec)
		if err != nil {
			// 删除上面的资源
			DeleteAppDeploy(form.Name, form.Namespace)
			return nil, err
		}
	}
	// TODO Nginx

	return &responses.OK, nil
}

// DeleteAppDeploy 删除
func DeleteAppDeploy(name, ns string) (*responses.Response, error) {
	var err1 error
	if _, err := service.DeleteConfigMap(name+"-configMap", ns); err != nil {
		err1 = err
	}
	if _, err := service.DeletePVC(ns, name+"-pvc"); err != nil {
		err1 = err
	}
	if _, err := service.DeleteDeploy(name, ns); err != nil {
		err1 = err
	}
	if _, err := service.DeleteService(name+"-service", ns); err != nil {
		err1 = err
	}
	if err1 != nil {
		return nil, err1
	}
	return &responses.OK, nil
}

// ListAppDeploy 列出ns下的所有appDeploy
func ListAppDeploy(ns string, label string) (*responses.AppDeployList, error) {
	list, err := dao.ClientSet.AppsV1().Deployments(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	num := len(list.Items)
	deployList := make([]responses.AppDeploy, num)
	for i, deploy := range list.Items {
		// 获取对应service
		serviceName := deploy.Name + "-service"
		// 针对spark定制化
		if deploy.Name == sparkMasterDeployName {
			serviceName = sparkMasterServiceName
		}
		svc, err := service.GetService(serviceName, ns)
		if err != nil {
			zap.S().Errorln("service/application/appDeploy:", err)
			svc = &corev1.Service{}
		}
		// 获取对应pvc
		pvcName := deploy.Name + "-pvc"
		pvc, err := service.GetPVC(ns, pvcName)
		if err != nil {
			zap.S().Errorln("service/application/appDeploy:", err)
			pvc = &corev1.PersistentVolumeClaim{}
		}
		// 获取挂载的路径
		volumeMounts := deploy.Spec.Template.Spec.Containers[0].VolumeMounts
		pathNum := len(volumeMounts)
		pvcPath := make([]string, pathNum)
		for i2, path := range volumeMounts {
			pvcPath[i2] = path.MountPath
		}
		// 获取资源信息
		limitCpu := deploy.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceCPU]
		limitMemory := deploy.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceMemory]
		limitStorage := deploy.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceEphemeralStorage]
		// 获取对应pod
		label1 := map[string]string{
			"uuid": deploy.Labels["uuid"],
		}
		selector := labels.SelectorFromSet(label1).String()
		podList, err := service.ListDeployPod(ns, selector)
		if err != nil {
			return nil, err
		}
		tmp := responses.AppDeploy{
			Name:              deploy.Name,
			Namespace:         deploy.Namespace,
			CreatedAt:         deploy.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Replicas:          deploy.Status.Replicas,
			Ports:             svc.Spec.Ports,
			Image:             deploy.Spec.Template.Spec.Containers[0].Image,
			UpdatedReplicas:   deploy.Status.UpdatedReplicas,
			ReadyReplicas:     deploy.Status.ReadyReplicas,
			AvailableReplicas: deploy.Status.AvailableReplicas,
			Resources: responses.Resources{
				Cpu:     limitCpu.String(),
				Memory:  limitMemory.String(),
				Storage: limitStorage.String(),
				PVC:     pvc.Spec.Resources.Requests.Storage().String(),
				// TODO GPU
			},
			Volume:  pvc.Spec.VolumeName,
			PvcPath: pvcPath,
			PodList: podList,
		}
		deployList[i] = tmp
	}
	return &responses.AppDeployList{Response: responses.OK, Length: num, DeployList: deployList}, nil
}

// GetAppDeploy 更新之前先获取deployApp的信息
func GetAppDeploy(name, ns string) (*forms.DeployAddForm, error) {
	form := forms.DeployAddForm{}
	deploy, err := service.GetDeploy(name, ns)
	if err != nil {
		return nil, err
	}
	strForm := deploy.Annotations["form"]
	err = json.Unmarshal([]byte(strForm), &form)
	if err != nil {
		return nil, err
	}
	return &form, nil
}

func UpdateAppDeploy(form forms.DeployAddForm) (*responses.Response, error) {
	// 将form序列化为string，存入deploy的注释
	jsonBytes, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}
	strForm := string(jsonBytes)
	// 准备工作
	// 分割申请资源
	m := int(form.Replicas)
	requestCpu, err := service.SplitRSC(form.Cpu, n*m)
	if err != nil {
		return nil, err
	}
	requestMemory, err := service.SplitRSC(form.Memory, n*m)
	if err != nil {
		return nil, err
	}
	requestStorage, err := service.SplitRSC(form.Storage, n*m)
	if err != nil {
		return nil, err
	}
	limitsCpu, err := service.SplitRSC(form.Cpu, m)
	if err != nil {
		return nil, err
	}
	limitsMemory, err := service.SplitRSC(form.Memory, m)
	if err != nil {
		return nil, err
	}
	limitsStorage, err := service.SplitRSC(form.Storage, m)
	if err != nil {
		return nil, err
	}
	configName := form.Name + "-configMap"
	pvcName := form.Name + "-pvc"
	serviceName := form.Name + "-service"
	// 更新configmap
	ns := form.Namespace
	if _, err := service.UpdateConfigMap(configName, ns, form.Env); err != nil {
		return nil, err
	}
	env := make([]corev1.EnvVar, len(form.Env))
	j := 0
	for k, _ := range form.Env {
		env[j] = corev1.EnvVar{
			Name: k,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: configName,
					},
					Key: k,
				},
			},
		}
		j++
	}
	// 更新deploy
	num := len(form.Ports)
	ports := make([]corev1.ContainerPort, num)
	for i, port := range form.Ports {
		ports[i] = corev1.ContainerPort{
			ContainerPort: port,
		}
	}
	deploy, err := service.GetDeploy(form.Name, form.Namespace)
	if err != nil {
		return nil, err
	}
	label := deploy.Labels
	// 创建PVC，持久存储
	var volumes []corev1.Volume
	var volumeMounts []corev1.VolumeMount
	if form.PvcStorage != "" {
		volumes = make([]corev1.Volume, 1)
		volumeMounts = make([]corev1.VolumeMount, 1)
		_, err = service.UpdateOrCreatePvc(form.Namespace, pvcName, form.StorageClassName, form.PvcStorage, accessModes)
		if err != nil {
			return nil, err
		}
		volumes[0] = corev1.Volume{
			Name: pvcName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvcName,
				},
			},
		}
		// 写死为/data目录
		volumeMounts[0] = corev1.VolumeMount{
			Name:      pvcName,
			MountPath: "/data",
		}
		//for i, path := range form.PvcPath {
		//	volumeMounts[i] = corev1.VolumeMount{
		//		Name:      pvcName,
		//		MountPath: path,
		//	}
		//}
	}
	spec := appsv1.DeploymentSpec{
		Replicas: &form.Replicas,
		Selector: &metav1.LabelSelector{MatchLabels: label},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: label},
			Spec: corev1.PodSpec{
				Volumes:       volumes,
				RestartPolicy: corev1.RestartPolicyAlways,
				Containers: []corev1.Container{
					{
						Image:           form.Image,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Command:         form.Command,
						Args:            form.Args,
						Ports:           ports,
						Env:             env,
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:              resource.MustParse(requestCpu),
								corev1.ResourceMemory:           resource.MustParse(requestMemory),
								corev1.ResourceEphemeralStorage: resource.MustParse(requestStorage),
								//TODO GPU
							},
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:              resource.MustParse(limitsCpu),
								corev1.ResourceMemory:           resource.MustParse(limitsMemory),
								corev1.ResourceEphemeralStorage: resource.MustParse(limitsStorage),
							},
						},
						VolumeMounts: volumeMounts,
					},
				},
			},
		},
	}
	if _, err := service.UpdateDeploy(form.Name, ns, strForm, spec); err != nil {
		return nil, err
	}

	// 更新service
	servicePorts := make([]corev1.ServicePort, num)
	for i, port := range form.Ports {
		servicePorts[i] = corev1.ServicePort{
			Name:       strconv.Itoa(int(port)),
			Port:       port,
			TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: port},
		}
	}
	serviceSpec := corev1.ServiceSpec{
		Type:     corev1.ServiceTypeNodePort,
		Selector: label,
		Ports:    servicePorts,
	}
	if _, err := service.UpdateService(serviceName, ns, serviceSpec); err != nil {
		return nil, err
	}

	// TODO niginx
	return &responses.OK, nil
}
