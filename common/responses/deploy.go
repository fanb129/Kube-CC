package responses

import (
	corev1 "k8s.io/api/core/v1"
)

// AppDeploy appDeploy的信息,用于查看详细信息
type AppDeploy struct {
	Name      string               `json:"name" form:"name" binding:"required,min=3,max=16"`
	Namespace string               `json:"namespace" form:"namespace" binding:"required,min=3,max=16"`
	Replicas  int32                `json:"replicas" form:"replicas" binding:"required,gte=1"`
	Image     string               `json:"image" form:"image" binding:"required"`
	Ports     []corev1.ServicePort `json:"ports"`
	Resources
	PvcPath           []string `json:"pvc_path"`
	Volume            string   `json:"volume"`
	CreatedAt         string   `json:"created_at"`
	UpdatedReplicas   int32    `json:"updated_replicas"`
	ReadyReplicas     int32    `json:"ready_replicas"`
	AvailableReplicas int32    `json:"available_replicas"`
	PodList           []Pod    `json:"pod_list"`
}

type AppDeployList struct {
	Response
	Length     int         `json:"length"`
	DeployList []AppDeploy `json:"deploy_list"`
}
