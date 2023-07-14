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

// CreateAppStatefulSet 创建statefulSet类型的整个应用app
// 包括 configmap、pvc、statefulSet、service、TODO ingress
func CreateAppStatefulSet(form forms.StatefulSetAddForm) (*responses.Response, error) {
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
	limitsGpu, err := service.SplitRSC(form.Gpu, m)
	if err != nil {
		return nil, err
	}
	// 创建uuid，以便筛选出属于同一组的deploy、pod、service等
	newUuid := string(uuid.NewUUID())
	label := map[string]string{
		"uuid": newUuid,
	}

	// 创建PVC，持久存储
	var pvcTemplate []corev1.PersistentVolumeClaim
	var volumeMounts []corev1.VolumeMount
	if form.PvcStorage != "" {
		pvcTemplate = make([]corev1.PersistentVolumeClaim, 1)
		//volumeMounts := make([]corev1.VolumeMount, len(form.PvcPath))
		volumeMounts = make([]corev1.VolumeMount, 1)
		pvcName := form.Name + "-pvc"
		if form.StorageClassName == "" {
			return nil, errors.New("已填写PvcStorage,StorageClassName不能为空")
		}
		pvcTemplate[0] = corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{Name: pvcName, Namespace: form.Namespace},
			Spec: corev1.PersistentVolumeClaimSpec{
				StorageClassName: &form.StorageClassName,
				AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.PersistentVolumeAccessMode(readWriteOnce)},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(form.PvcStorage),
					},
				},
			},
		}
		// 写死为/data目录
		volumeMounts[0] = corev1.VolumeMount{
			Name:      pvcName,
			MountPath: "/data",
		}
	}

	// 0.创建configMap，存储环境变量
	configName := form.Name + "-configMap"
	_, err = service.CreateConfigMap(configName, form.Namespace, label, form.Env)
	if err != nil {
		DeleteAppSetfulset(form.Name, form.Namespace)
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

	// 1. 创建sts
	// 端口
	num := len(form.Ports)
	ports := make([]corev1.ContainerPort, num)
	for i, port := range form.Ports {
		ports[i] = corev1.ContainerPort{
			ContainerPort: port,
		}
	}
	spec := appsv1.StatefulSetSpec{
		Replicas: &form.Replicas,
		Selector: &metav1.LabelSelector{MatchLabels: label},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: label},
			Spec: corev1.PodSpec{
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
							},
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:              resource.MustParse(limitsCpu),
								corev1.ResourceMemory:           resource.MustParse(limitsMemory),
								corev1.ResourceEphemeralStorage: resource.MustParse(limitsStorage),
								service.GpuShare:                resource.MustParse(limitsGpu),
							},
						},
						VolumeMounts: volumeMounts,
					},
				},
			},
		},
		VolumeClaimTemplates: pvcTemplate,
	}
	_, err = service.CreateStatefulSet(form.Name, form.Namespace, strForm, label, spec)
	if err != nil {
		DeleteAppSetfulset(form.Name, form.Namespace)
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
		DeleteAppSetfulset(form.Name, form.Namespace)
		return nil, err
	}

	// TODO Nginx

	return &responses.OK, nil
}

// DeleteAppSetfulset 删除
func DeleteAppSetfulset(name, ns string) (*responses.Response, error) {
	var err1 error
	if _, err := service.DeleteConfigMap(name+"-configMap", ns); err != nil {
		err1 = err
	}
	if _, err := service.DeletePVC(ns, name+"-pvc"); err != nil {
		err1 = err
	}
	if _, err := service.DeleteStatefulSet(name, ns); err != nil {
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

// ListAppStatesulSet 列出ns下的所有 AppStatesulSet
func ListAppStatesulSet(ns string) (*responses.AppStatefulSetList, error) {
	list, err := dao.ClientSet.AppsV1().StatefulSets(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	num := len(list.Items)
	stsList := make([]responses.AppStatefulSet, num)
	for i, sts := range list.Items {
		// 获取对应service
		serviceName := sts.Name + "-service"
		svc, err := service.GetService(serviceName, ns)
		if err != nil {
			zap.S().Errorln("service/application/appStatefulSet:", err)
			svc = &corev1.Service{}
		}
		// 获取对应pvc
		label := map[string]string{
			"uuid": sts.Labels["uuid"],
		}
		selector := labels.SelectorFromSet(label).String()
		pvcList, err := service.ListPVC(ns, selector)
		if err != nil {
			return nil, err
		}
		// 获取挂载的路径
		volumeMounts := sts.Spec.Template.Spec.Containers[0].VolumeMounts
		pathNum := len(volumeMounts)
		pvcPath := make([]string, pathNum)
		for i2, path := range volumeMounts {
			pvcPath[i2] = path.MountPath
		}
		// 获取资源信息
		limitCpu := sts.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceCPU]
		limitMemory := sts.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceMemory]
		limitStorage := sts.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceEphemeralStorage]
		limitGpu := sts.Spec.Template.Spec.Containers[0].Resources.Limits[service.GpuShare]
		// 获取对应pod
		podList, err := service.ListStatefulSetPod(ns, selector)
		if err != nil {
			return nil, err
		}
		tmp := responses.AppStatefulSet{
			Name:              sts.Name,
			Namespace:         sts.Namespace,
			CreatedAt:         sts.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Replicas:          sts.Status.Replicas,
			Ports:             svc.Spec.Ports,
			Image:             sts.Spec.Template.Spec.Containers[0].Image,
			UpdatedReplicas:   sts.Status.UpdatedReplicas,
			ReadyReplicas:     sts.Status.ReadyReplicas,
			AvailableReplicas: sts.Status.AvailableReplicas,
			Resources: responses.Resources{
				Cpu:     limitCpu.String(),
				Memory:  limitMemory.String(),
				Storage: limitStorage.String(),
				GPU:     limitGpu.String(),
			},
			PvcPath: pvcPath,
			PodList: podList,
			PvcList: pvcList.PvcList,
		}
		stsList[i] = tmp
	}
	return &responses.AppStatefulSetList{Response: responses.OK, Length: num, StatefulSetList: stsList}, nil
}

// GetAppStatefulSet 更新之前先获取deployApp的信息
func GetAppStatefulSet(name, ns string) (*forms.StatefulSetAddForm, error) {
	form := forms.StatefulSetAddForm{}
	sts, err := service.GetStatefulSet(name, ns)
	if err != nil {
		return nil, err
	}
	strForm := sts.Annotations["form"]
	err = json.Unmarshal([]byte(strForm), &form)
	if err != nil {
		return nil, err
	}
	return &form, nil
}

func UpdateAppStatefulSet(form forms.StatefulSetAddForm) (*responses.Response, error) {
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
	limitsGpu, err := service.SplitRSC(form.Gpu, m)
	if err != nil {
		return nil, err
	}
	configName := form.Name + "-configMap"
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
	// 更新statefulSet
	num := len(form.Ports)
	ports := make([]corev1.ContainerPort, num)
	for i, port := range form.Ports {
		ports[i] = corev1.ContainerPort{
			ContainerPort: port,
		}
	}
	sts, err := service.GetStatefulSet(form.Name, form.Namespace)
	if err != nil {
		return nil, err
	}
	label := sts.Labels
	// 创建PVC，持久存储
	var pvcTemplate []corev1.PersistentVolumeClaim
	var volumeMounts []corev1.VolumeMount
	if form.PvcStorage != "" {
		pvcTemplate = make([]corev1.PersistentVolumeClaim, 1)
		//volumeMounts := make([]corev1.VolumeMount, len(form.PvcPath))
		volumeMounts = make([]corev1.VolumeMount, 1)
		pvcName := form.Name + "-pvc"
		if form.StorageClassName == "" {
			return nil, errors.New("已填写PvcStorage,StorageClassName不能为空")
		}
		pvcTemplate[0] = corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{Name: pvcName, Namespace: form.Namespace},
			Spec: corev1.PersistentVolumeClaimSpec{
				StorageClassName: &form.StorageClassName,
				AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.PersistentVolumeAccessMode(readWriteOnce)},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(form.PvcStorage),
					},
				},
			},
		}
		// 写死为/data目录
		volumeMounts[0] = corev1.VolumeMount{
			Name:      pvcName,
			MountPath: "/data",
		}
	}
	spec := appsv1.StatefulSetSpec{
		Replicas: &form.Replicas,
		Selector: &metav1.LabelSelector{MatchLabels: label},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: label},
			Spec: corev1.PodSpec{
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
							},
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:              resource.MustParse(limitsCpu),
								corev1.ResourceMemory:           resource.MustParse(limitsMemory),
								corev1.ResourceEphemeralStorage: resource.MustParse(limitsStorage),
								service.GpuShare:                resource.MustParse(limitsGpu),
							},
						},
						VolumeMounts: volumeMounts,
					},
				},
			},
		},
		VolumeClaimTemplates: pvcTemplate,
	}
	if _, err := service.UpdateStatefulSet(form.Name, ns, strForm, spec); err != nil {
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
