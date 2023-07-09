package application

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/conf"
	"Kube-CC/service"
	"encoding/json"
	"errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/uuid"
	"time"
)

const (
	sparkMasterDeployName  = "spark-master"
	sparkWorkerDeployName  = "spark-worker"
	sparkMasterServiceName = "spark-master"         // 特定名字，不能改，与spark镜像关联
	sparkWorkerServiceName = "spark-worker-service" // 特定名字，不能改，与spark镜像关联
	sparkIngressName       = "spark-ingress"
)

// CreateSpark 为uid创建spark，masterReplicas默认1， masterReplicas默认2
func CreateSpark(name, u_id string, masterReplicas int32, workerReplicas int32, expiredTime *time.Time, resources forms.ApplyResources) (*responses.Response, error) {
	// uuid
	newUuid := string(uuid.NewUUID())
	ns := name + "-" + newUuid
	label := map[string]string{
		"image": "spark",
		"u_id":  u_id,
	}
	masterLabel := map[string]string{
		"component": "spark-master",
		"uuid":      newUuid + "1",
	}
	workerLabel := map[string]string{
		"component": "spark-worker",
		"uuid":      newUuid + "2",
	}
	rsc := forms.Resources{
		Cpu:        resources.Cpu,
		Memory:     resources.Memory,
		Storage:    resources.Storage,
		PvcStorage: resources.PvcStorage,
		Gpu:        resources.Gpu,
	}
	// 将form序列化为string，存入deploy的注释
	form := forms.SparkUpdateForm{
		Name:           ns,
		MasterReplicas: masterReplicas,
		WorkerReplicas: workerReplicas,
		ExpiredTime:    expiredTime,
		ApplyResources: resources,
	}
	jsonBytes, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}
	strForm := string(jsonBytes)
	// 准备工作
	// 分割申请资源
	m := int(masterReplicas + workerReplicas)
	requestCpu, err := service.SplitRSC(resources.Cpu, n*m)
	if err != nil {
		return nil, err
	}
	requestMemory, err := service.SplitRSC(resources.Memory, n*m)
	if err != nil {
		return nil, err
	}
	requestStorage, err := service.SplitRSC(resources.Storage, n*m)
	if err != nil {
		return nil, err
	}
	limitsCpu, err := service.SplitRSC(resources.Cpu, m)
	if err != nil {
		return nil, err
	}
	limitsMemory, err := service.SplitRSC(resources.Memory, m)
	if err != nil {
		return nil, err
	}
	limitsStorage, err := service.SplitRSC(resources.Storage, m)
	if err != nil {
		return nil, err
	}
	// 创建namespace
	_, err = service.CreateNs(ns, strForm, expiredTime, label, rsc)
	if err != nil {
		DeleteSpark(ns)
		return nil, err
	}
	// 创建PVC，持久存储
	var masterVolumes, workerVolumes []corev1.Volume
	var masterVolumeMounts, workerVolumeMounts []corev1.VolumeMount
	if resources.PvcStorage != "" {
		masterVolumes = make([]corev1.Volume, 1)
		workerVolumes = make([]corev1.Volume, 1)
		masterVolumeMounts = make([]corev1.VolumeMount, 1)
		workerVolumeMounts = make([]corev1.VolumeMount, 1)
		// 分割资源
		pvcStorage, err := service.SplitRSC(resources.PvcStorage, 2)
		if err != nil {
			DeleteSpark(ns)
			return nil, err
		}

		if resources.StorageClassName == "" {
			DeleteSpark(ns)
			return nil, errors.New("已填写PvcStorage,StorageClassName不能为空")
		}
		masterPvcName := sparkMasterDeployName + "-pvc"
		workerPvcName := sparkWorkerDeployName + "-pvc"
		_, err = service.CreatePVC(ns, masterPvcName, resources.StorageClassName, pvcStorage, accessModes)
		if err != nil {
			DeleteSpark(ns)
			return nil, err
		}
		_, err = service.CreatePVC(ns, workerPvcName, resources.StorageClassName, pvcStorage, accessModes)
		if err != nil {
			DeleteSpark(ns)
			return nil, err
		}
		masterVolumes[0] = corev1.Volume{
			Name: masterPvcName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: masterPvcName,
				},
			},
		}
		workerVolumes[0] = corev1.Volume{
			Name: workerPvcName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: workerPvcName,
				},
			},
		}
		// 写死为/data目录
		masterVolumeMounts[0] = corev1.VolumeMount{
			Name:      masterPvcName,
			MountPath: "/data",
		}
		workerVolumeMounts[0] = corev1.VolumeMount{
			Name:      workerPvcName,
			MountPath: "/data",
		}
	}
	// spark的master控制器
	masterSpec := appsv1.DeploymentSpec{
		Replicas: &masterReplicas,
		Selector: &metav1.LabelSelector{MatchLabels: masterLabel},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: masterLabel},
			Spec: corev1.PodSpec{
				Volumes: masterVolumes,
				Containers: []corev1.Container{
					{
						Name:            "spark-master",
						Image:           conf.SparkImage,
						ImagePullPolicy: corev1.PullIfNotPresent, // 镜像拉取策略
						Command:         []string{"/start-master"},
						Args:            []string{conf.SshPwd},
						Ports: []corev1.ContainerPort{
							{ContainerPort: 7077},
							{ContainerPort: 8080},
							{ContainerPort: 22},
						},
						VolumeMounts: masterVolumeMounts,
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
					},
				},
				RestartPolicy: corev1.RestartPolicyAlways,
			},
		},
	}
	_, err = service.CreateDeploy(sparkMasterDeployName, ns, "", masterLabel, masterSpec)
	if err != nil {
		DeleteSpark(ns)
		return nil, err
	}

	// spark的master的service
	masterServiceSpec := corev1.ServiceSpec{
		Type:     corev1.ServiceTypeNodePort,
		Selector: masterLabel,
		Ports: []corev1.ServicePort{ // 默认生成nodePort
			{Name: "spark", Port: 7077, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7077}},
			{Name: "http", Port: 8080, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8080}},
			{Name: "ssh", Port: 22, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 22}},
		},
	}
	_, err = service.CreateService(sparkMasterServiceName, ns, masterLabel, masterServiceSpec)
	if err != nil {
		DeleteSpark(ns)
		return nil, err
	}

	// spark的worker的控制器
	workerSpec := appsv1.DeploymentSpec{
		Replicas: &workerReplicas,
		Selector: &metav1.LabelSelector{MatchLabels: workerLabel},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: workerLabel},
			Spec: corev1.PodSpec{
				Volumes: workerVolumes,
				Containers: []corev1.Container{
					{Name: "spark-worker",
						Image:           conf.SparkImage,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Command:         []string{"/start-worker"},
						Args:            []string{conf.SshPwd},
						Ports: []corev1.ContainerPort{
							{ContainerPort: 8081},
							{ContainerPort: 22},
						},
						VolumeMounts: workerVolumeMounts,
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
					},
				},
				RestartPolicy: corev1.RestartPolicyAlways,
			},
		},
	}
	_, err = service.CreateDeploy(sparkWorkerDeployName, ns, "", workerLabel, workerSpec)
	if err != nil {
		DeleteSpark(ns)
		return nil, err
	}

	// spark的worker的service
	workerServiceSpec := corev1.ServiceSpec{
		Type:     corev1.ServiceTypeNodePort,
		Selector: workerLabel,
		Ports: []corev1.ServicePort{
			{Name: "http", Port: 8081, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8081}},
			{Name: "ssh", Port: 22, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 22}},
		},
	}
	_, err = service.CreateService(sparkWorkerServiceName, ns, workerLabel, workerServiceSpec)
	if err != nil {
		DeleteSpark(ns)
		return nil, err
	}

	// TODO spark的ingress
	//ingressSpec := v1beta1.IngressSpec{
	//	Rules: []v1beta1.IngressRule{
	//		// 代理master服务
	//		{
	//			Host: fmt.Sprintf("spark.%s", conf.ProjectName),
	//			IngressRuleValue: v1beta1.IngressRuleValue{
	//				HTTP: &v1beta1.HTTPIngressRuleValue{
	//					Paths: []v1beta1.HTTPIngressPath{
	//						{
	//							Path:    "/master-" + s,
	//							Backend: v1beta1.IngressBackend{ServiceName: sparkMasterServiceName, ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 8080}},
	//						},
	//					},
	//				},
	//			},
	//		},
	//		// 代理worker服务
	//		{
	//			Host: fmt.Sprintf("spark.%s", conf.ProjectName),
	//			IngressRuleValue: v1beta1.IngressRuleValue{
	//				HTTP: &v1beta1.HTTPIngressRuleValue{
	//					Paths: []v1beta1.HTTPIngressPath{
	//						{
	//							Path:    "/worker-" + s,
	//							Backend: v1beta1.IngressBackend{ServiceName: sparkWorkerServiceName, ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 8081}},
	//						},
	//					},
	//				},
	//			},
	//		},
	//	},
	//}
	//_, err = service.CreateIngress(sparkIngressName, ns, map[string]string{}, ingressSpec)
	//if err != nil {
	//	return nil, err
	//}

	return &responses.OK, nil
}

// ListSpark  获取uid用户下的所有spark
func ListSpark(u_id string) (*responses.SparkListResponse, error) {
	label := map[string]string{
		"image": "spark",
	}
	if u_id != "" {
		label["u_id"] = u_id
	}
	// 将map标签转换为string
	selector := labels.SelectorFromSet(label).String()
	sparks, err := service.ListNs(selector)
	if err != nil {
		return nil, err
	}
	sparkList := make([]responses.Spark, sparks.Length)
	for i, spark := range sparks.NsList {
		// 获取deploy
		deploy, err := ListAppDeploy(spark.Name, "")
		if err != nil {
			return nil, err
		}
		sparkList[i] = responses.Spark{
			Ns:         spark,
			DeployList: deploy.DeployList,
		}
	}

	return &responses.SparkListResponse{
		Response:  responses.OK,
		Length:    sparks.Length,
		SparkList: sparkList,
	}, nil
}

// DeleteSpark 删除spark
func DeleteSpark(ns string) (*responses.Response, error) {
	var err1 error
	if _, err := service.DeleteIngress(sparkIngressName, ns); err != nil {
		err1 = err
	}
	if _, err := service.DeleteService(sparkWorkerServiceName, ns); err != nil {
		err1 = err
	}
	if _, err := service.DeleteService(sparkMasterServiceName, ns); err != nil {
		err1 = err
	}
	if _, err := service.DeleteDeploy(sparkWorkerDeployName, ns); err != nil {
		err1 = err
	}
	if _, err := service.DeleteDeploy(sparkMasterDeployName, ns); err != nil {
		err1 = err
	}
	if _, err := service.DeletePVC(ns, sparkMasterDeployName+"-pvc"); err != nil {
		err1 = err
	}
	if _, err := service.DeletePVC(ns, sparkWorkerDeployName+"-pvc"); err != nil {
		err1 = err
	}
	if _, err := service.DeleteNs(ns); err != nil {
		err1 = err
	}
	if err1 != nil {
		return nil, err1
	}
	return &responses.OK, nil
}

// UpdateSpark 更新spark以及replicas
func UpdateSpark(name string, masterReplicas int32, workerReplicas int32, expiredTime *time.Time, resources forms.ApplyResources) (*responses.Response, error) {
	rsc := forms.Resources{
		Cpu:        resources.Cpu,
		Memory:     resources.Memory,
		Storage:    resources.Storage,
		PvcStorage: resources.PvcStorage,
		Gpu:        resources.Gpu,
	}
	// 将form序列化为string，存入deploy的注释
	form := forms.SparkUpdateForm{
		Name:           name,
		MasterReplicas: masterReplicas,
		WorkerReplicas: workerReplicas,
		ExpiredTime:    expiredTime,
		ApplyResources: resources,
	}
	jsonBytes, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}
	strForm := string(jsonBytes)
	if _, err := service.UpdateNs(name, strForm, expiredTime, rsc); err != nil {
		return nil, err
	}
	// 准备工作
	// 分割申请资源
	m := int(masterReplicas + workerReplicas)
	requestCpu, err := service.SplitRSC(resources.Cpu, n*m)
	if err != nil {
		return nil, err
	}
	requestMemory, err := service.SplitRSC(resources.Memory, n*m)
	if err != nil {
		return nil, err
	}
	requestStorage, err := service.SplitRSC(resources.Storage, n*m)
	if err != nil {
		return nil, err
	}
	limitsCpu, err := service.SplitRSC(resources.Cpu, m)
	if err != nil {
		return nil, err
	}
	limitsMemory, err := service.SplitRSC(resources.Memory, m)
	if err != nil {
		return nil, err
	}
	limitsStorage, err := service.SplitRSC(resources.Storage, m)
	if err != nil {
		return nil, err
	}
	// 创建PVC，持久存储
	var masterVolumes, workerVolumes []corev1.Volume
	var masterVolumeMounts, workerVolumeMounts []corev1.VolumeMount
	// pvc不为空，更新或创建,为空则删除，但是删不掉，哈哈哈
	if resources.PvcStorage != "" {
		masterVolumes = make([]corev1.Volume, 1)
		workerVolumes = make([]corev1.Volume, 1)
		masterVolumeMounts = make([]corev1.VolumeMount, 1)
		workerVolumeMounts = make([]corev1.VolumeMount, 1)
		// 分割资源
		pvcStorage, err := service.SplitRSC(resources.PvcStorage, 2)
		if err != nil {
			return nil, err
		}

		if resources.StorageClassName == "" {
			return nil, errors.New("已填写PvcStorage,StorageClassName不能为空")
		}
		masterPvcName := sparkMasterDeployName + "-pvc"
		workerPvcName := sparkWorkerDeployName + "-pvc"
		_, err = service.UpdateOrCreatePvc(name, masterPvcName, resources.StorageClassName, pvcStorage, accessModes)
		if err != nil {
			return nil, err
		}
		_, err = service.UpdateOrCreatePvc(name, workerPvcName, resources.StorageClassName, pvcStorage, accessModes)
		if err != nil {
			return nil, err
		}
		masterVolumes[0] = corev1.Volume{
			Name: masterPvcName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: masterPvcName,
				},
			},
		}
		workerVolumes[0] = corev1.Volume{
			Name: workerPvcName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: workerPvcName,
				},
			},
		}
		// 写死为/data目录
		masterVolumeMounts[0] = corev1.VolumeMount{
			Name:      masterPvcName,
			MountPath: "/data",
		}
		workerVolumeMounts[0] = corev1.VolumeMount{
			Name:      workerPvcName,
			MountPath: "/data",
		}
	}
	// 更新master的Replicas
	master, err := service.GetDeploy(sparkMasterDeployName, name)
	if err != nil {
		return nil, err
	}
	master.Spec.Replicas = &masterReplicas
	master.Spec.Template.Spec.Containers[0].Resources = corev1.ResourceRequirements{
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
	}
	master.Spec.Template.Spec.Volumes = masterVolumes
	master.Spec.Template.Spec.Containers[0].VolumeMounts = masterVolumeMounts
	if _, err := service.UpdateDeploy(sparkMasterDeployName, name, "", master.Spec); err != nil {
		return nil, err
	}

	// 更新worker的Replicas
	worker, err := service.GetDeploy(sparkWorkerDeployName, name)
	if err != nil {
		return nil, err
	}
	worker.Spec.Replicas = &workerReplicas
	worker.Spec.Template.Spec.Containers[0].Resources = corev1.ResourceRequirements{
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
	}
	worker.Spec.Template.Spec.Volumes = workerVolumes
	worker.Spec.Template.Spec.Containers[0].VolumeMounts = workerVolumeMounts
	if _, err := service.UpdateDeploy(sparkWorkerDeployName, name, "", worker.Spec); err != nil {
		return nil, err
	}

	return &responses.OK, nil
}

// GetSpark 更新之前先获取信息
func GetSpark(name string) (*forms.SparkUpdateForm, error) {
	form := forms.SparkUpdateForm{}
	ns, err := service.GetNs(name)
	if err != nil {
		return nil, err
	}
	strForm := ns.Annotations["form"]
	err = json.Unmarshal([]byte(strForm), &form)
	if err != nil {
		return nil, err
	}
	return &form, nil
}
