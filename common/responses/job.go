package responses

import corev1 "k8s.io/api/core/v1"

type Job struct {
	Name              string                   `json:"name"`
	Namespace         string                   `json:"namespace"`
	CreatedAt         string                   `json:"created_at"`
	Phase             corev1.PodPhase          `json:"phase"`
	NodeIp            string                   `json:"node_ip"`
	ContainerStatuses []corev1.ContainerStatus `json:"container_statuses"`
}

type JobListResponse struct {
	Response
	Length  int   `json:"length"`
	JobList []Job `json:"job_list"`
}

type JobInfo struct {
	Response
	Info corev1.Pod `json:"info"`
}
