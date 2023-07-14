package responses

import (
	corev1 "k8s.io/api/core/v1"
)

// AppDeploy appDeploy的信息,用于查看详细信息
type AppDeploy struct {
	Name      string               `json:"name"`
	Namespace string               `json:"namespace"`
	Replicas  int32                `json:"replicas"`
	Image     string               `json:"image"`
	Ports     []corev1.ServicePort `json:"ports"`
	Resources
	PvcPath           []string    `json:"pvc_path"`
	Volume            string      `json:"volume"`
	CreatedAt         string      `json:"created_at"`
	UpdatedReplicas   int32       `json:"updated_replicas"`
	ReadyReplicas     int32       `json:"ready_replicas"`
	AvailableReplicas int32       `json:"available_replicas"`
	PodList           []DeployPod `json:"pod_list"`
}

type AppDeployList struct {
	Response
	Length     int         `json:"length"`
	DeployList []AppDeploy `json:"deploy_list"`
}
