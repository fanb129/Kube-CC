package service

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/conf"
	"strconv"
	"time"
)

const (
	HDFS_MASTER_SERVICE         = "HDFS_MASTER_SERVICE"
	HDOOP_YARN_MASTER           = "HDOOP_YARN_MASTER"
	HADOOP_NODE_TYPE            = "HADOOP_NODE_TYPE"
	hadoopConfigMapName         = "hadoop-configmap"
	hadoopHdfsMasterDeployName  = "hadoop-hdfs-master-deploy"
	hadoopHdfsMasterServiceName = "hadoop-hdfs-master-service"
	datanodeDeployName          = "hadoop-datanode-deploy"
	hadoopYarnMasterDeployName  = "hadoop-yarn-master-deploy"
	hadoopYarnMasterServiceName = "hadoop-yarn-master-service"
	hadoopYarnNodeDeployName    = "hadoop-yarn-node-deploy"
	hadoopYarnNodeServiceName   = "hadoop-yarn-node-service"
)

// CreateHadoop 创建hadoop  hdfsMasterReplicas,datanodeReplicas,yarnMasterReplicas,yarnNodeReplicas 默认1，3，1，3
func CreateHadoop(u_id uint, hdfsMasterReplicas, datanodeReplicas, yarnMasterReplicas, yarnNodeReplicas int32) (*common.Response, error) {
	// 获取当前时间戳，纳秒
	s := strconv.FormatInt(time.Now().UnixNano(), 10)
	ns := "hadoop-" + s
	label := map[string]string{
		"image": "hadoop",
	}
	hdfsMasterSelector := map[string]string{
		"name": "hdfs-master",
	}
	hdfsMasterLabel := map[string]string{
		"name": "hdfs-master",
	}
	datanodeSelector := map[string]string{
		"name": "hadoop-datanode",
	}
	datanodeLabel := map[string]string{
		"name": "hadoop-datanode",
	}
	yarnMasterSelector := map[string]string{
		"name": "yarn-master",
	}
	yarnMasterLabel := map[string]string{
		"name": "yarn-master",
	}
	yarnNodeSelector := map[string]string{
		"name": "yarn-node",
	}
	yarnNodeLabel := map[string]string{
		"name": "yarn-node",
	}
	if u_id != 0 {
		uid := strconv.Itoa(int(u_id))
		label["u_id"] = uid
		hdfsMasterLabel["u_id"] = uid
		datanodeLabel["u_id"] = uid
		yarnMasterLabel["u_id"] = uid
		yarnNodeLabel["u_id"] = uid
	}
	// 创建namespace
	_, err := CreateNs(ns, label)
	if err != nil {
		return nil, err
	}

	// 创建configMap
	_, err = CreateConfigMap(hadoopConfigMapName, ns, label, map[string]string{
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
		Selector: &metav1.LabelSelector{MatchLabels: hdfsMasterSelector},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: hdfsMasterLabel,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:            "hdfs-master",
						Image:           conf.HadoopImage,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports: []corev1.ContainerPort{
							{ContainerPort: 9000},
							{ContainerPort: 50070},
						},
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
	_, err = CreateDeploy(hadoopHdfsMasterDeployName, ns, label, hdfsMasterSpec)
	if err != nil {
		return nil, err
	}

	// 创建hdfs-master的service
	hdfsMasterServiceSpec := corev1.ServiceSpec{
		Type:     corev1.ServiceTypeNodePort,
		Selector: hdfsMasterSelector,
		Ports: []corev1.ServicePort{
			{Name: "rpc", Port: 9000, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9000}},
			{Name: "http", Port: 50070, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 50070}},
		},
	}
	_, err = CreateService(hadoopHdfsMasterServiceName, ns, label, hdfsMasterServiceSpec)
	if err != nil {
		return nil, err
	}

	// 创建datanode的deploy
	//var datanodeReplicas int32 = 3
	datanodeSpec := appsv1.DeploymentSpec{
		Replicas: &datanodeReplicas,
		Selector: &metav1.LabelSelector{MatchLabels: datanodeSelector},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: datanodeLabel,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:            "hadoop-datanode",
						Image:           conf.HadoopImage,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports: []corev1.ContainerPort{
							{ContainerPort: 9000},
							{ContainerPort: 50070},
						},
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
	_, err = CreateDeploy(datanodeDeployName, ns, label, datanodeSpec)
	if err != nil {
		return nil, err
	}

	// 创建yarn-master的deploy
	//var yarnMasterReplicas int32 = 1
	yarnMasterSpec := appsv1.DeploymentSpec{
		Replicas: &yarnMasterReplicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: yarnMasterSelector,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: yarnMasterLabel,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:            "yarn-master",
						Image:           conf.HadoopImage,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports: []corev1.ContainerPort{
							{ContainerPort: 9000},
							{ContainerPort: 50070},
						},
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
	_, err = CreateDeploy(hadoopYarnMasterDeployName, ns, label, yarnMasterSpec)
	if err != nil {
		return nil, err
	}

	// yarn-master的service
	yarnMasterServiceSpec := corev1.ServiceSpec{
		Type:     corev1.ServiceTypeNodePort,
		Selector: yarnMasterSelector,
		Ports: []corev1.ServicePort{
			{Name: "8030", Port: 8030},
			{Name: "8031", Port: 8031},
			{Name: "8032", Port: 8032},
			{Name: "http", Port: 8088, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8088}},
		},
	}
	_, err = CreateService(hadoopYarnMasterServiceName, ns, label, yarnMasterServiceSpec)
	if err != nil {
		return nil, err
	}

	// yarn-node的deploy
	//var yarnNodeReplicas int32 = 3
	yarnNodeSpec := appsv1.DeploymentSpec{
		Replicas: &yarnNodeReplicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: yarnNodeSelector,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: yarnNodeLabel,
			},
			Spec: corev1.PodSpec{
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
	_, err = CreateDeploy(hadoopYarnNodeDeployName, ns, label, yarnNodeSpec)
	if err != nil {
		return nil, err
	}

	// yarn-node的service
	yarnNodeServiceSpec := corev1.ServiceSpec{
		Selector: yarnNodeSelector,
		Ports: []corev1.ServicePort{
			{Port: 8040},
		},
	}
	_, err = CreateService(hadoopYarnNodeServiceName, ns, label, yarnNodeServiceSpec)
	if err != nil {
		return nil, err
	}

	return &common.OK, nil
}

// GetHadoop 获取uid下的所有hadoop
func GetHadoop(u_id uint) (*common.HadoopListResponse, error) {
	label := map[string]string{
		"image": "hadoop",
	}
	if u_id > 0 {
		label["u_id"] = strconv.Itoa(int(u_id))
	}
	// 将map标签转换为string
	selector := labels.SelectorFromSet(label).String()
	hadoops, err := GetNs(selector)
	if err != nil {
		return nil, err
	}
	hadoopList := make([]common.Hadoop, hadoops.Length)
	for i, hadoop := range hadoops.NsList {
		// 获取pod
		podList, err := GetPod(hadoop.Name, "")
		if err != nil {
			return nil, err
		}
		// 获取deploy
		deployList, err := GetDeploy(hadoop.Name, "")
		if err != nil {
			return nil, err
		}
		var hdfsMaster, datanode, yarnMaster, yarnNode int32
		for j := 0; j < deployList.Length; j++ {
			deploy := deployList.DeployList[j]
			switch deploy.Name {
			case hadoopHdfsMasterDeployName:
				hdfsMaster = deploy.Replicas
			case datanodeDeployName:
				datanode = deploy.Replicas
			case hadoopYarnMasterDeployName:
				yarnMaster = deploy.Replicas
			case hadoopYarnNodeDeployName:
				yarnNode = deploy.Replicas
			}
		}
		// 获取service
		serviceList, err := GetService(hadoop.Name, "")
		if err != nil {
			return nil, err
		}
		hadoopList[i] = common.Hadoop{
			Name:               hadoop.Name,
			Uid:                u_id,
			Username:           hadoop.Username,
			Nickname:           hadoop.Nickname,
			CreatedAt:          hadoop.CreatedAt,
			Status:             hadoop.Status,
			PodList:            podList.PodList,
			DeployList:         deployList.DeployList,
			ServiceList:        serviceList.ServiceList,
			HdfsMasterReplicas: hdfsMaster,
			DatanodeReplicas:   datanode,
			YarnMasterReplicas: yarnMaster,
			YarnNodeReplicas:   yarnNode,
		}
	}
	return &common.HadoopListResponse{
		Response:   common.OK,
		Length:     hadoops.Length,
		HadoopList: hadoopList,
	}, nil
}

// DeleteHadoop 删除指定hadoop
func DeleteHadoop(ns string) (*common.Response, error) {
	var err1 error
	if _, err := DeleteService(hadoopYarnNodeServiceName, ns); err != nil {
		err1 = err
	}
	if _, err := DeleteDeploy(hadoopYarnNodeDeployName, ns); err != nil {
		err1 = err
	}
	if _, err := DeleteService(hadoopYarnMasterServiceName, ns); err != nil {
		err1 = err
	}
	if _, err := DeleteDeploy(hadoopYarnMasterDeployName, ns); err != nil {
		err1 = err
	}
	if _, err := DeleteDeploy(datanodeDeployName, ns); err != nil {
		err1 = err
	}
	if _, err := DeleteService(hadoopHdfsMasterServiceName, ns); err != nil {
		err1 = err
	}
	if _, err := DeleteDeploy(hadoopHdfsMasterDeployName, ns); err != nil {
		err1 = err
	}
	if _, err := DeleteConfigMap(hadoopConfigMapName, ns); err != nil {
		err1 = err
	}
	if _, err := DeleteNs(ns); err != nil {
		err1 = err
	}
	if err1 != nil {
		return nil, err1
	}
	return &common.OK, nil
}

// UpdateHadoop 更新hadoop的uid，以及replicas
func UpdateHadoop(name, uid string, hdfsMasterReplicas, datanodeReplicas, yarnMasterReplicas, yarnNodeReplicas int32) (*common.Response, error) {
	if _, err := UpdateNs(name, uid); err != nil {
		return nil, err
	}

	// 更新hdfsMaster的Replicas
	hdfsMaster, err := GetADeploy(hadoopHdfsMasterDeployName, name)
	if err != nil {
		return nil, err
	}
	hdfsMaster.Spec.Replicas = &hdfsMasterReplicas
	if _, err := UpdateDeploy(hdfsMaster); err != nil {
		return nil, err
	}

	// 更新datanode的Replicas
	datanode, err := GetADeploy(datanodeDeployName, name)
	if err != nil {
		return nil, err
	}
	datanode.Spec.Replicas = &datanodeReplicas
	if _, err := UpdateDeploy(datanode); err != nil {
		return nil, err
	}

	// 更新yarnMaster的Replicas
	yarnMaster, err := GetADeploy(hadoopYarnMasterDeployName, name)
	if err != nil {
		return nil, err
	}
	yarnMaster.Spec.Replicas = &yarnMasterReplicas
	if _, err := UpdateDeploy(yarnMaster); err != nil {
		return nil, err
	}

	// 更新yarnNode的Replicas
	yarnNode, err := GetADeploy(hadoopYarnNodeDeployName, name)
	if err != nil {
		return nil, err
	}
	yarnNode.Spec.Replicas = &yarnNodeReplicas
	if _, err := UpdateDeploy(yarnNode); err != nil {
		return nil, err
	}

	return &common.OK, nil
}
