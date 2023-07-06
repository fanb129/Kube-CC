package application

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/conf"
	"Kube-CC/service"
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
	HDFS_MASTER_SERVICE         = "HDFS_MASTER_SERVICE"
	HDOOP_YARN_MASTER           = "HDOOP_YARN_MASTER"
	HADOOP_NODE_TYPE            = "HADOOP_NODE_TYPE"
	hadoopConfigMapName         = "hadoop-configmap"
	hadoopHdfsMasterDeployName  = "hadoop-hdfs-master"
	hadoopHdfsMasterServiceName = "hadoop-hdfs-master-service"
	datanodeDeployName          = "hadoop-datanode"
	hadoopYarnMasterDeployName  = "hadoop-yarn-master"
	hadoopYarnMasterServiceName = "hadoop-yarn-master-service"
	hadoopYarnNodeDeployName    = "hadoop-yarn-node"
	hadoopYarnNodeServiceName   = "hadoop-yarn-node-service"
)

// CreateHadoop 创建hadoop  hdfsMasterReplicas,datanodeReplicas,yarnMasterReplicas,yarnNodeReplicas 默认1，3，1，3
func CreateHadoop(u_id string, hdfsMasterReplicas, datanodeReplicas, yarnMasterReplicas, yarnNodeReplicas int32, expiredTime *time.Time, resources forms.ApplyResources) (*responses.Response, error) {

	newUuid := string(uuid.NewUUID())
	ns := "hadoop-" + newUuid
	label := map[string]string{
		"image": "hadoop",
		"uuid":  newUuid,
	}
	label["u_id"] = u_id
	hdfsMasterLabel := map[string]string{
		"name": "hdfs-master",
		"uuid": newUuid,
	}
	datanodeLabel := map[string]string{
		"name": "hadoop-datanode",
		"uuid": newUuid,
	}
	yarnMasterLabel := map[string]string{
		"name": "yarn-master",
		"uuid": newUuid,
	}
	yarnNodeLabel := map[string]string{
		"name": "yarn-node",
		"uuid": newUuid,
	}
	// 创建namespace
	rsc := forms.Resources{
		Cpu:        resources.Cpu,
		Memory:     resources.Memory,
		Storage:    resources.Storage,
		PvcStorage: resources.PvcStorage,
		Gpu:        resources.Gpu,
	}
	// 准备工作
	// 分割申请资源
	m := int(hdfsMasterReplicas + yarnMasterReplicas + yarnNodeReplicas + datanodeReplicas)
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

	_, err = service.CreateNs(ns, expiredTime, label, rsc)
	if err != nil {
		return nil, err
	}
	// 创建PVC，持久存储
	volumes := make([]corev1.Volume, 0)
	volumeMounts := make([]corev1.VolumeMount, len(resources.PvcPath))
	if resources.PvcStorage != "" {
		if resources.StorageClassName == "" {
			return nil, errors.New("已填写PvcStorage,StorageClassName不能为空")
		}
		pvcName := ns + "-pvc"
		_, err = service.CreatePVC(ns, pvcName, resources.StorageClassName, resources.PvcStorage, accessModes)
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
		for i, path := range resources.PvcPath {
			volumeMounts[i] = corev1.VolumeMount{
				Name:      pvcName,
				MountPath: path,
			}
		}
	}
	// 创建configMap
	_, err = service.CreateConfigMap(hadoopConfigMapName, ns, map[string]string{}, map[string]string{
		HDFS_MASTER_SERVICE: "hadoop-hdfs-master",
		HDOOP_YARN_MASTER:   "hadoop-yarn-master",
	})
	if err != nil {
		return nil, err
	}

	// 创建hdfs-master的deploy
	//var hdfsMasterReplicas int32 = 1
	hdfsMasterSpec := appsv1.DeploymentSpec{
		Replicas: &hdfsMasterReplicas,
		Selector: &metav1.LabelSelector{MatchLabels: hdfsMasterLabel},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: hdfsMasterLabel,
			},
			Spec: corev1.PodSpec{
				Volumes: volumes,
				Containers: []corev1.Container{
					{
						Name:            "hdfs-master",
						Image:           conf.HadoopImage,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports: []corev1.ContainerPort{
							{ContainerPort: 9000},
							{ContainerPort: 50070},
						},
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
						Env: []corev1.EnvVar{
							{Name: HADOOP_NODE_TYPE, Value: "namenode"},
							{Name: HDFS_MASTER_SERVICE, ValueFrom: &corev1.EnvVarSource{
								ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: hadoopConfigMapName,
									},
									Key: HDFS_MASTER_SERVICE,
								},
							}},
							{Name: HDOOP_YARN_MASTER, ValueFrom: &corev1.EnvVarSource{
								ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: hadoopConfigMapName,
									},
									Key: HDOOP_YARN_MASTER,
								},
							}},
						},
					},
				},
				RestartPolicy: corev1.RestartPolicyAlways,
			},
		},
	}
	_, err = service.CreateDeploy(hadoopHdfsMasterDeployName, ns, "", hdfsMasterLabel, hdfsMasterSpec)
	if err != nil {
		return nil, err
	}

	// 创建hdfs-master的service
	hdfsMasterServiceSpec := corev1.ServiceSpec{
		Type:     corev1.ServiceTypeNodePort,
		Selector: hdfsMasterLabel,
		Ports: []corev1.ServicePort{
			{Name: "rpc", Port: 9000, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9000}},
			{Name: "http", Port: 50070, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 50070}},
		},
	}
	_, err = service.CreateService(hadoopHdfsMasterServiceName, ns, hdfsMasterLabel, hdfsMasterServiceSpec)
	if err != nil {
		return nil, err
	}

	// 创建datanode的deploy
	//var datanodeReplicas int32 = 3
	datanodeSpec := appsv1.DeploymentSpec{
		Replicas: &datanodeReplicas,
		Selector: &metav1.LabelSelector{MatchLabels: datanodeLabel},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: datanodeLabel,
			},
			Spec: corev1.PodSpec{
				Volumes: volumes,
				Containers: []corev1.Container{
					{
						Name:            "hadoop-datanode",
						Image:           conf.HadoopImage,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports: []corev1.ContainerPort{
							{ContainerPort: 9000},
							{ContainerPort: 50070},
						},
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
						Env: []corev1.EnvVar{
							{
								Name:  HADOOP_NODE_TYPE,
								Value: "datanode",
							},
							{
								Name: HDFS_MASTER_SERVICE,
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: hadoopConfigMapName,
										},
										Key: HDFS_MASTER_SERVICE,
									},
								},
							},
							{
								Name: HDOOP_YARN_MASTER,
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: hadoopConfigMapName,
										},
										Key: HDOOP_YARN_MASTER,
									},
								},
							},
						},
					},
				},
				RestartPolicy: corev1.RestartPolicyAlways,
			},
		},
	}
	_, err = service.CreateDeploy(datanodeDeployName, ns, "", datanodeLabel, datanodeSpec)
	if err != nil {
		return nil, err
	}

	// 创建yarn-master的deploy
	//var yarnMasterReplicas int32 = 1
	yarnMasterSpec := appsv1.DeploymentSpec{
		Replicas: &yarnMasterReplicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: yarnMasterLabel,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: yarnMasterLabel,
			},
			Spec: corev1.PodSpec{
				Volumes: volumes,
				Containers: []corev1.Container{
					{
						Name:            "yarn-master",
						Image:           conf.HadoopImage,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports: []corev1.ContainerPort{
							{ContainerPort: 9000},
							{ContainerPort: 50070},
						},
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
						Env: []corev1.EnvVar{
							{
								Name:  HADOOP_NODE_TYPE,
								Value: "resourceman",
							},
							{
								Name: HDFS_MASTER_SERVICE,
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: hadoopConfigMapName,
										},
										Key: HDFS_MASTER_SERVICE,
									},
								},
							},
							{
								Name: HDOOP_YARN_MASTER,
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: hadoopConfigMapName,
										},
										Key: HDOOP_YARN_MASTER,
									},
								},
							},
						},
					},
				},
				RestartPolicy: corev1.RestartPolicyAlways,
			},
		},
	}
	_, err = service.CreateDeploy(hadoopYarnMasterDeployName, ns, "", yarnMasterLabel, yarnMasterSpec)
	if err != nil {
		return nil, err
	}

	// yarn-master的service
	yarnMasterServiceSpec := corev1.ServiceSpec{
		Type:     corev1.ServiceTypeNodePort,
		Selector: yarnMasterLabel,
		Ports: []corev1.ServicePort{
			{Name: "8030", Port: 8030},
			{Name: "8031", Port: 8031},
			{Name: "8032", Port: 8032},
			{Name: "http", Port: 8088, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8088}},
		},
	}
	_, err = service.CreateService(hadoopYarnMasterServiceName, ns, yarnMasterLabel, yarnMasterServiceSpec)
	if err != nil {
		return nil, err
	}

	// yarn-node的deploy
	//var yarnNodeReplicas int32 = 3
	yarnNodeSpec := appsv1.DeploymentSpec{
		Replicas: &yarnNodeReplicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: yarnNodeLabel,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: yarnNodeLabel,
			},
			Spec: corev1.PodSpec{
				Volumes: volumes,
				Containers: []corev1.Container{
					{
						Name:            "yarn-node",
						Image:           conf.HadoopImage,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports: []corev1.ContainerPort{
							{ContainerPort: 8040},
							{ContainerPort: 8041},
							{ContainerPort: 8042},
						},
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
						Env: []corev1.EnvVar{
							{Name: HADOOP_NODE_TYPE, Value: "yarnnode"},
							{
								Name: HDFS_MASTER_SERVICE,
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: hadoopConfigMapName,
										},
										Key: HDFS_MASTER_SERVICE,
									},
								},
							},
							{
								Name: HDOOP_YARN_MASTER,
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: hadoopConfigMapName,
										},
										Key: HDOOP_YARN_MASTER,
									},
								},
							},
						},
					},
				},
				RestartPolicy: corev1.RestartPolicyAlways,
			},
		},
	}
	_, err = service.CreateDeploy(hadoopYarnNodeDeployName, ns, "", yarnNodeLabel, yarnNodeSpec)
	if err != nil {
		return nil, err
	}

	// yarn-node的service
	yarnNodeServiceSpec := corev1.ServiceSpec{
		Selector: yarnNodeLabel,
		Ports: []corev1.ServicePort{
			{Port: 8040},
		},
	}
	_, err = service.CreateService(hadoopYarnNodeServiceName, ns, yarnNodeLabel, yarnNodeServiceSpec)
	if err != nil {
		return nil, err
	}

	return &responses.OK, nil
}

// ListHadoop 获取uid下的所有hadoop
func ListHadoop(u_id string) (*responses.HadoopListResponse, error) {
	label := map[string]string{
		"image": "hadoop",
	}
	if u_id != "" {
		label["u_id"] = u_id
	}
	// 将map标签转换为string
	selector := labels.SelectorFromSet(label).String()
	hadoops, err := service.ListNs(selector)
	if err != nil {
		return nil, err
	}
	hadoopList := make([]responses.Hadoop, hadoops.Length)
	for i, hadoop := range hadoops.NsList {
		// 获取deploy
		deployList, err := service.ListDeploy(hadoop.Name, "")
		if err != nil {
			return nil, err
		}
		// 获取service
		serviceList, err := service.ListService(hadoop.Name, "")
		if err != nil {
			return nil, err
		}
		// 获取pod
		podList, err := service.ListPod(hadoop.Name, "")
		hadoopList[i] = responses.Hadoop{
			Ns:          hadoop,
			DeployList:  deployList,
			ServiceList: serviceList,
			PodList:     podList,
		}
	}
	return &responses.HadoopListResponse{
		Response:   responses.OK,
		Length:     hadoops.Length,
		HadoopList: hadoopList,
	}, nil
}

// DeleteHadoop 删除指定hadoop
func DeleteHadoop(ns string) (*responses.Response, error) {
	var err1 error
	if _, err := service.DeleteService(hadoopYarnNodeServiceName, ns); err != nil {
		err1 = err
	}
	if _, err := service.DeleteDeploy(hadoopYarnNodeDeployName, ns); err != nil {
		err1 = err
	}
	if _, err := service.DeleteService(hadoopYarnMasterServiceName, ns); err != nil {
		err1 = err
	}
	if _, err := service.DeleteDeploy(hadoopYarnMasterDeployName, ns); err != nil {
		err1 = err
	}
	if _, err := service.DeleteDeploy(datanodeDeployName, ns); err != nil {
		err1 = err
	}
	if _, err := service.DeleteService(hadoopHdfsMasterServiceName, ns); err != nil {
		err1 = err
	}
	if _, err := service.DeleteDeploy(hadoopHdfsMasterDeployName, ns); err != nil {
		err1 = err
	}
	if _, err := service.DeleteConfigMap(hadoopConfigMapName, ns); err != nil {
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

// UpdateHadoop 更新hadoop的uid，以及replicas
func UpdateHadoop(name string, hdfsMasterReplicas, datanodeReplicas, yarnMasterReplicas, yarnNodeReplicas int32, expiredTime *time.Time, resources forms.ApplyResources) (*responses.Response, error) {
	rsc := forms.Resources{
		Cpu:        resources.Cpu,
		Memory:     resources.Memory,
		Storage:    resources.Storage,
		PvcStorage: resources.PvcStorage,
		Gpu:        resources.Gpu,
	}
	if _, err := service.UpdateNs(name, expiredTime, rsc); err != nil {
		return nil, err
	}

	// 准备工作
	// 分割申请资源
	m := int(hdfsMasterReplicas + yarnMasterReplicas + yarnNodeReplicas + datanodeReplicas)
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

	// 更新hdfsMaster的Replicas
	hdfsMaster, err := service.GetDeploy(hadoopHdfsMasterDeployName, name)
	if err != nil {
		return nil, err
	}
	hdfsMaster.Spec.Replicas = &hdfsMasterReplicas
	hdfsMaster.Spec.Template.Spec.Containers[0].Resources = corev1.ResourceRequirements{
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
	if _, err := service.UpdateDeploy(hadoopHdfsMasterDeployName, name, "", hdfsMaster.Spec); err != nil {
		return nil, err
	}

	// 更新datanode的Replicas
	datanode, err := service.GetDeploy(datanodeDeployName, name)
	if err != nil {
		return nil, err
	}
	datanode.Spec.Replicas = &datanodeReplicas
	datanode.Spec.Template.Spec.Containers[0].Resources = corev1.ResourceRequirements{
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
	if _, err := service.UpdateDeploy(datanodeDeployName, name, "", datanode.Spec); err != nil {
		return nil, err
	}

	// 更新yarnMaster的Replicas
	yarnMaster, err := service.GetDeploy(hadoopYarnMasterDeployName, name)
	if err != nil {
		return nil, err
	}
	yarnMaster.Spec.Replicas = &yarnMasterReplicas
	yarnMaster.Spec.Template.Spec.Containers[0].Resources = corev1.ResourceRequirements{
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
	if _, err := service.UpdateDeploy(hadoopYarnMasterDeployName, name, "", yarnMaster.Spec); err != nil {
		return nil, err
	}

	// 更新yarnNode的Replicas
	yarnNode, err := service.GetDeploy(hadoopYarnNodeDeployName, name)
	if err != nil {
		return nil, err
	}
	yarnNode.Spec.Replicas = &yarnNodeReplicas
	yarnNode.Spec.Template.Spec.Containers[0].Resources = corev1.ResourceRequirements{
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
	if _, err := service.UpdateDeploy(hadoopYarnNodeDeployName, name, "", yarnNode.Spec); err != nil {
		return nil, err
	}

	return &responses.OK, nil
}
