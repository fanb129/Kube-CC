package application

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"Kube-CC/service"
	"context"
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
	// 准备工作
	// 分割申请资源
	requestCpu, err := service.SplitRSC(form.Cpu, n)
	if err != nil {
		return nil, err
	}
	requestMemory, err := service.SplitRSC(form.Memory, n)
	if err != nil {
		return nil, err
	}
	requestStorage, err := service.SplitRSC(form.Storage, n)
	if err != nil {
		return nil, err
	}

	// 创建uuid，以便筛选出属于同一组的deploy、pod、service等
	newUuid := string(uuid.NewUUID())
	label := map[string]string{
		"uuid": newUuid,
	}

	// 创建PVC，持久存储
	volumes := make([]corev1.Volume, 0)
	volumeMounts := make([]corev1.VolumeMount, len(form.PvcPath))
	if form.PvcStorage != "" {
		if form.StorageClassName == "" {
			return nil, errors.New("已填写PvcStorage,StorageClassName不能为空")
		}
		pvcName := form.Name + "-pvc"
		_, err = service.CreatePVC(form.Namespace, pvcName, form.StorageClassName, form.PvcStorage, accessModes)
		if err != nil {
			return nil, err
		}
		volumes = append(volumes, corev1.Volume{
			Name: pvcName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvcName,
				},
			},
		})
		for i, path := range form.PvcPath {
			volumeMounts[i] = corev1.VolumeMount{
				Name:      pvcName,
				MountPath: path,
			}
		}
	}

	// 0.创建configMap，存储环境变量
	configName := form.Name + "-configMap"
	_, err = service.CreateConfigMap(configName, form.Namespace, label, form.Env)
	if err != nil {
		DeleteAppDeploy(form.Name, form.Namespace)
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
						Image:           form.Image,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Command:         form.Command,
						Args:            form.Args,
						Ports:           ports,
						Env:             env,
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceRequestsCPU:              resource.MustParse(requestCpu),
								corev1.ResourceRequestsMemory:           resource.MustParse(requestMemory),
								corev1.ResourceRequestsEphemeralStorage: resource.MustParse(requestStorage),
								//TODO GPU
							},
							Limits: corev1.ResourceList{
								corev1.ResourceLimitsCPU:              resource.MustParse(form.Cpu),
								corev1.ResourceLimitsMemory:           resource.MustParse(form.Memory),
								corev1.ResourceLimitsEphemeralStorage: resource.MustParse(form.Storage),
							},
						},
						VolumeMounts: volumeMounts,
					},
				},
			},
		},
	}
	_, err = service.CreateDeploy(form.Name, form.Namespace, label, spec)
	if err != nil {
		DeleteAppDeploy(form.Name, form.Namespace)
		return nil, err
	}

	// 2 创建service
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
		limitCpu := deploy.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceLimitsCPU]
		limitMemory := deploy.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceLimitsMemory]
		limitStorage := deploy.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceLimitsEphemeralStorage]
		// 获取对应pod
		label := map[string]string{
			"uuid": deploy.Labels["uuid"],
		}
		selector := labels.SelectorFromSet(label).String()
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
	configName := name + "-configMap"
	pvcName := name + "-pvc"
	configMap, err := service.GetConfigMap(configName, ns)
	if err != nil {
		zap.S().Errorln("service/application/appDeploy:", err)
		configMap = &corev1.ConfigMap{}
	}
	pvc, err := service.GetPVC(ns, pvcName)
	if err != nil {
		zap.S().Errorln("service/application/appDeploy:", err)
		pvc = &corev1.PersistentVolumeClaim{}
	}
	deploy, err := service.GetDeploy(name, ns)
	if err != nil {
		return nil, err
	}

	// 取出ports参数
	portList := deploy.Spec.Template.Spec.Containers[0].Ports
	ports := make([]int32, len(portList))
	for i, port := range portList {
		ports[i] = port.ContainerPort
	}

	// 取出资源参数
	cpu := deploy.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceLimitsCPU]
	memory := deploy.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceLimitsMemory]
	storage := deploy.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceLimitsEphemeralStorage]
	pvcStorage := pvc.Spec.Resources.Requests[corev1.ResourceStorage]
	// 取出挂载路径
	mounts := deploy.Spec.Template.Spec.Containers[0].VolumeMounts
	paths := make([]string, len(mounts))
	for i, mount := range mounts {
		paths[i] = mount.MountPath
	}
	deployApp := forms.DeployAddForm{
		Name:      name,
		Namespace: ns,
		Replicas:  *deploy.Spec.Replicas,
		Image:     deploy.Spec.Template.Spec.Containers[0].Image,
		Command:   deploy.Spec.Template.Spec.Containers[0].Command,
		Args:      deploy.Spec.Template.Spec.Containers[0].Args,
		Ports:     ports,
		Env:       configMap.Data,
		ApplyResources: forms.ApplyResources{
			Cpu:              cpu.String(),
			Memory:           memory.String(),
			Storage:          storage.String(),
			PvcStorage:       pvcStorage.String(),
			StorageClassName: *pvc.Spec.StorageClassName,
			PvcPath:          paths,
			// TODO GPU
		},
	}
	return &deployApp, nil
}

func UpdateAppDeploy(form forms.DeployAddForm) (*responses.Response, error) {
	// 准备工作
	// 分割申请资源
	requestCpu, err := service.SplitRSC(form.Cpu, n)
	if err != nil {
		return nil, err
	}
	requestMemory, err := service.SplitRSC(form.Memory, n)
	if err != nil {
		return nil, err
	}
	requestStorage, err := service.SplitRSC(form.Storage, n)
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
	volumes := make([]corev1.Volume, 0)
	volumeMounts := make([]corev1.VolumeMount, len(form.PvcPath))
	if form.PvcStorage != "" {
		err = service.UpdatePVC(form.Namespace, pvcName, form.PvcStorage)
		if err != nil {
			return nil, err
		}
		volumes = append(volumes, corev1.Volume{
			Name: pvcName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvcName,
				},
			},
		})
		for i, path := range form.PvcPath {
			volumeMounts[i] = corev1.VolumeMount{
				Name:      pvcName,
				MountPath: path,
			}
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
						Image:           form.Image,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Command:         form.Command,
						Args:            form.Args,
						Ports:           ports,
						Env:             env,
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceRequestsCPU:              resource.MustParse(requestCpu),
								corev1.ResourceRequestsMemory:           resource.MustParse(requestMemory),
								corev1.ResourceRequestsEphemeralStorage: resource.MustParse(requestStorage),
								//TODO GPU
							},
							Limits: corev1.ResourceList{
								corev1.ResourceLimitsCPU:              resource.MustParse(form.Cpu),
								corev1.ResourceLimitsMemory:           resource.MustParse(form.Memory),
								corev1.ResourceLimitsEphemeralStorage: resource.MustParse(form.Storage),
							},
						},
						VolumeMounts: volumeMounts,
					},
				},
			},
		},
	}
	if _, err := service.UpdateDeploy(form.Name, ns, spec); err != nil {
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
