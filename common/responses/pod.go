package responses

type DeployPod struct {
	Name   string `json:"name"`
	Phase  string `json:"phase"`
	HostIP string `json:"host_ip"`
	PodIP  string `json:"pod_ip"`
}

type StsPod struct {
	Name   string `json:"name"`
	Phase  string `json:"phase"`
	HostIP string `json:"host_ip"`
	PodIP  string `json:"pod_ip"`
}

type JobPod struct {
	Name     string `json:"name"`
	Phase    string `json:"phase"`
	Restarts int32  `json:"restarts"` // 重启次数
	HostIP   string `json:"host_ip"`
	PodIP    string `json:"pod_ip"`
}

type PodLogResponse struct {
	Response
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Log       string `json:"log"`
}
