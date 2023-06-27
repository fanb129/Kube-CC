package responses

import corev1 "k8s.io/api/core/v1"

type Ns struct {
	Name        string                `json:"name"`
	Status      corev1.NamespacePhase `json:"status"`
	CreatedAt   string                `json:"created_at"`
	Username    string                `json:"username"`
	Nickname    string                `json:"nickname"`
	Uid         uint                  `json:"u_id"`
	ExpiredTime string                `json:"expired_time"`
	Cpu         string                `json:"cpu"`
	UsedCpu     string                `json:"used_cpu"`
	Memory      string                `json:"memory"`
	UsedMemory  string                `json:"used_memory"`
}

type NsListResponse struct {
	Response
	Length int  `json:"length"`
	NsList []Ns `json:"ns_list"`
}
