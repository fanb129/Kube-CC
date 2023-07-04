package responses

type DeployPod struct {
	Name   string `json:"name"`
	Phase  string `json:"phase"`
	HostIP string `json:"host_ip"`
	PodIP  string `json:"pod_ip"`
}
