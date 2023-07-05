package application

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"Kube-CC/service"
	"context"
	"errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/uuid"
	"strconv"
)

// CreateAppStatefulSet创建statefulSet类型的整个应用app
// 包括 configmap、pvc、deploy、service、TODO ingress
func CreateAppStatefulSet(form forms.StatefulSetAddForm) (*responses.Response, error) {
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
			return nil, err
		}
		// 获取对应pvc
		pvcName := sts.Name + "-pvc"
		pvc, err := service.GetPVC(ns, pvcName)
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
		limitCpu := sts.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceLimitsCPU]
		limitMemory := sts.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceLimitsMemory]
		limitStorage := sts.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceLimitsEphemeralStorage]
		// 获取对应pod
		label := map[string]string{
			"uuid": sts.Labels["uuid"],
		}
		selector := labels.SelectorFromSet(label).String()
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
				PVC:     pvc.Spec.Resources.Requests.Storage().String(),
				// TODO GPU
			},
			Volume:  pvc.Spec.VolumeName,
			PvcPath: pvcPath,
			PodList: podList,
		}
		stsList[i] = tmp
	}
	return &responses.AppStatefulSetList{Response: responses.OK, Length: num, StatefulSetList: stsList}, nil
}

// GetAppStatefulSet 更新之前先获取deployApp的信息
func GetAppStatefulSet(name, ns string) (*forms.StatefulSetAddForm, error) {
	configName := name + "-configMap"
	pvcName := name + "-pvc"
	configMap, err := service.GetConfigMap(configName, ns)
	if err != nil {
		return nil, err
	}
	pvc, err := service.GetPVC(ns, pvcName)
	if err != nil {
		return nil, err
	}
	sts, err := service.GetStatefulSet(name, ns)

	// 取出ports参数
	portList := sts.Spec.Template.Spec.Containers[0].Ports
	ports := make([]int32, len(portList))
	for i, port := range portList {
		ports[i] = port.ContainerPort
	}

	// 取出资源参数
	cpu := sts.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceLimitsCPU]
	memory := sts.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceLimitsMemory]
	storage := sts.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceLimitsEphemeralStorage]
	pvcStorage := pvc.Spec.Resources.Requests[corev1.ResourceStorage]
	// 取出挂载路径
	mounts := sts.Spec.Template.Spec.Containers[0].VolumeMounts
	paths := make([]string, len(mounts))
	for i, mount := range mounts {
		paths[i] = mount.MountPath
	}
	stsApp := forms.StatefulSetAddForm{
		Name:      name,
		Namespace: ns,
		Replicas:  *sts.Spec.Replicas,
		Image:     sts.Spec.Template.Spec.Containers[0].Image,
		Command:   sts.Spec.Template.Spec.Containers[0].Command,
		Args:      sts.Spec.Template.Spec.Containers[0].Args,
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
	return &stsApp, nil
}

func UpdateAppStatefulSet(form forms.StatefulSetAddForm) (*responses.Response, error) {
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
	spec := appsv1.StatefulSetSpec{
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
	if _, err := service.UpdateStatefulSet(form.Name, ns, spec); err != nil {
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
