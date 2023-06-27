package responses

import appsv1 "k8s.io/api/apps/v1"

type Deploy struct {
	Name              string `json:"name"`
	Namespace         string `json:"namespace"`
	CreatedAt         string `json:"created_at"`
	Replicas          int32  `json:"replicas"`
	UpdatedReplicas   int32  `json:"updated_replicas"`
	ReadyReplicas     int32  `json:"ready_replicas"`
	AvailableReplicas int32  `json:"available_replicas"`
	Uid               string `json:"u_id"`
}

// DeployListResponse pod控制器返回结果
type DeployListResponse struct {
	Response
	Length     int      `json:"length"`
	DeployList []Deploy `json:"deploy_list"`
}

type DeployInfo struct {
	Response
	Info appsv1.Deployment `json:"info"`
}
