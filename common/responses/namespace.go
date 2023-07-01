package responses

import corev1 "k8s.io/api/core/v1"

type Resources struct {
	Cpu         string `json:"cpu"`
	UsedCpu     string `json:"used_cpu"`
	Memory      string `json:"memory"`
	UsedMemory  string `json:"used_memory"`
	Storage     string `json:"storage"`
	UsedStorage string `json:"used_storage"`
	PVC         string `json:"pvc"`
	UsedPVC     string `json:"used_pvc"`
	GPU         string `json:"gpu"`
	UsedGPU     string `json:"used_gpu"`
}

type Ns struct {
	Name        string                `json:"name"`
	Status      corev1.NamespacePhase `json:"status"`
	CreatedAt   string                `json:"created_at"`
	Username    string                `json:"username"`
	Nickname    string                `json:"nickname"`
	Uid         uint                  `json:"u_id"`
	ExpiredTime string                `json:"expired_time"`
	Resources
}

type NsListResponse struct {
	Response
	Length int  `json:"length"`
	NsList []Ns `json:"ns_list"`
}
