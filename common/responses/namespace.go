package responses

import corev1 "k8s.io/api/core/v1"

type Resources struct {
	Cpu          string `json:"cpu"`
	UsedCpu      string `json:"used_cpu"`
	UsedCpuValue int64  `json:"used_cpu_value"`

	Memory          string `json:"memory"`
	UsedMemory      string `json:"used_memory"`
	UsedMemoryValue int64  `json:"used_memory_value"`

	Storage          string `json:"storage"`
	UsedStorage      string `json:"used_storage"`
	UsedStorageValue int64  `json:"used_storage_value"`

	PVC          string `json:"pvc"`
	UsedPVC      string `json:"used_pvc"`
	UsedPVCValue int64  `json:"used_pvc_value"`

	GPU          string `json:"gpu"`
	UsedGPU      string `json:"used_gpu"`
	UsedGPUValue int64  `json:"used_gpu_value"`
}

type Ns struct {
	Name      string                `json:"name"`
	Status    corev1.NamespacePhase `json:"status"`
	CreatedAt string                `json:"created_at"`
	Username  string                `json:"username"`
	Nickname  string                `json:"nickname"`
	Uid       uint                  `json:"u_id"`
	//ExpiredTime string                `json:"expired_time"`
	Resources
}

type NsListResponse struct {
	Response
	Length int  `json:"length"`
	NsList []Ns `json:"ns_list"`
}

type UserTotalNs struct {
	Response
	Cpu      string  `json:"cpu"`
	UsedCpu  string  `json:"used_cpu"`
	CpuRatio float64 `json:"cpu_ratio"`

	Memory      string  `json:"memory"`
	UsedMemory  string  `json:"used_memory"`
	MemoryRatio float64 `json:"memory_ratio"`

	Storage      string  `json:"storage"`
	UsedStorage  string  `json:"used_storage"`
	StorageRatio float64 `json:"storage_ratio"`

	PVC      string  `json:"pvc"`
	UsedPVC  string  `json:"used_pvc"`
	PvcRatio float64 `json:"pvc_ratio"`

	GPU      string  `json:"gpu"`
	UsedGPU  string  `json:"used_gpu"`
	GpuRatio float64 `json:"gpu_ratio"`
}
