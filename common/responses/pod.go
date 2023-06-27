package responses

import corev1 "k8s.io/api/core/v1"

type Pod struct {
	Name              string                   `json:"name"`
	Namespace         string                   `json:"namespace"`
	CreatedAt         string                   `json:"created_at"`
	Phase             corev1.PodPhase          `json:"phase"`
	NodeIp            string                   `json:"node_ip"`
	ContainerStatuses []corev1.ContainerStatus `json:"container_statuses"`
	Uid               string                   `json:"u_id"`
}

type PodListResponse struct {
	Response
	Length  int   `json:"length"`
	PodList []Pod `json:"pod_list"`
}

type PodInfo struct {
	Response
	Info corev1.Pod `json:"info"`
}
