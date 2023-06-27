package responses

import corev1 "k8s.io/api/core/v1"

type Node struct {
	Name           string                 `json:"name"`
	Ip             string                 `json:"ip"`
	Ready          corev1.ConditionStatus `json:"ready"`
	CreatedAt      string                 `json:"created_at"`
	OsImage        string                 `json:"os_image"`
	KubeletVersion string                 `json:"kubelet_version"`
	CPU            string                 `json:"cpu"`
	Memory         string                 `json:"memory"`
}

type NodeListResponse struct {
	Response
	Length   int    `json:"length"`
	NodeList []Node `json:"node_list"`
}
