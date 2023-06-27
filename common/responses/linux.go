package responses

import corev1 "k8s.io/api/core/v1"

type Linux struct {
	Name        string                `json:"name"`
	Uid         uint                  `json:"u_id"`
	Username    string                `json:"username"`
	Nickname    string                `json:"nickname"`
	Status      corev1.NamespacePhase `json:"status"`
	CreatedAt   string                `json:"created_at"`
	PodList     []Pod                 `json:"pod_list"`
	DeployList  []Deploy              `json:"deploy_list"`
	ServiceList []Service             `json:"service_list"`
	ExpiredTime string                `json:"expired_time"`
	Cpu         string                `json:"cpu"`
	UsedCpu     string                `json:"used_cpu"`
	Memory      string                `json:"memory"`
	UsedMemory  string                `json:"used_memory"`
}

type LinuxListResponse struct {
	Response
	Image     string  `json:"image"`
	Length    int     `json:"length"`
	LinuxList []Linux `json:"linux_list"`
}
