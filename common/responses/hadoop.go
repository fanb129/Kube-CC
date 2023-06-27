package responses

import corev1 "k8s.io/api/core/v1"

type Hadoop struct {
	Name               string                `json:"name"`
	Uid                uint                  `json:"u_id"`
	Username           string                `json:"username"`
	Nickname           string                `json:"nickname"`
	Status             corev1.NamespacePhase `json:"status"`
	CreatedAt          string                `json:"created_at"`
	PodList            []Pod                 `json:"pod_list"`
	DeployList         []Deploy              `json:"deploy_list"`
	ServiceList        []Service             `json:"service_list"`
	HdfsMasterReplicas int32                 `json:"hdfs_master_replicas"`
	DatanodeReplicas   int32                 `json:"datanode_replicas"`
	YarnMasterReplicas int32                 `json:"yarn_master_replicas"`
	YarnNodeReplicas   int32                 `json:"yarn_node_replicas"`
	ExpiredTime        string                `json:"expired_time"`
	Cpu                string                `json:"cpu"`
	UsedCpu            string                `json:"used_cpu"`
	Memory             string                `json:"memory"`
	UsedMemory         string                `json:"used_memory"`
}

type HadoopListResponse struct {
	Response
	Length     int      `json:"length"`
	HadoopList []Hadoop `json:"hadoop_list"`
}
