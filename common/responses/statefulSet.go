package responses

import (
	"Kube-CC/common/forms"
	corev1 "k8s.io/api/core/v1"
)

// AppStatefulSet 的信息,用于查看详细信息
type AppStatefulSet struct {
	Name      string               `json:"name"`
	Namespace string               `json:"namespace"`
	Replicas  int32                `json:"replicas"`
	Image     string               `json:"image"`
	Ports     []corev1.ServicePort `json:"ports"`
	Resources
	PvcPath           []string `json:"pvc_path"`
	CreatedAt         string   `json:"created_at"`
	UpdatedReplicas   int32    `json:"updated_replicas"`
	ReadyReplicas     int32    `json:"ready_replicas"`
	AvailableReplicas int32    `json:"available_replicas"`
	PodList           []StsPod `json:"pod_list"`
	PvcList           []Pvc    `json:"pvc_list"`
}

type AppStatefulSetList struct {
	Response
	Length          int              `json:"length"`
	StatefulSetList []AppStatefulSet `json:"stateful_set_list"`
}

type InfoStatefulSet struct {
	Response
	Form forms.StatefulSetAddForm
}
