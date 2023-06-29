package responses

import appsv1 "k8s.io/api/apps/v1"

type StatefulSet struct {
	Name string `json:"name"`
	//ServiceName       string `json:"serviceName"`
	Namespace       string `json:"namespace"`
	CreatedAt       string `json:"created_at"`
	Replicas        int32  `json:"replicas"`
	UpdatedReplicas int32  `json:"updated_replicas"`
	ReadyReplicas   int32  `json:"ready_replicas"`
	CurrentReplicas int32  `json:"current_replicas"`
	CurrentRevision string `json:"current_revision"`
	Uid             string `json:"u_id"`
}

// StatefulSetListResponse pod控制器返回结果
type StatefulSetListResponse struct {
	Response
	Length          int           `json:"length"`
	StatefulSetList []StatefulSet `json:"stateful_set_list"`
}

type StatefulSetInfo struct {
	Response
	Info appsv1.StatefulSet `json:"info"`
}
