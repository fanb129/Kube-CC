package responses

import corev1 "k8s.io/api/core/v1"

type Spark struct {
	Name           string                `json:"name"`
	Uid            uint                  `json:"u_id"`
	Username       string                `json:"username"`
	Nickname       string                `json:"nickname"`
	CreatedAt      string                `json:"created_at"`
	Status         corev1.NamespacePhase `json:"status"`
	PodList        []Pod                 `json:"pod_list"`
	DeployList     []Deploy              `json:"deploy_list"`
	ServiceList    []Service             `json:"service_list"`
	IngressList    []Ingress             `json:"ingress_list"`
	MasterReplicas int32                 `json:"master_replicas"`
	WorkerReplicas int32                 `json:"worker_replicas"`
	ExpiredTime    string                `json:"expired_time"`
	Cpu            string                `json:"cpu"`
	UsedCpu        string                `json:"used_cpu"`
	Memory         string                `json:"memory"`
	UsedMemory     string                `json:"used_memory"`
}

type SparkListResponse struct {
	Response
	Length    int     `json:"length"`
	SparkList []Spark `json:"spark_list"`
}
