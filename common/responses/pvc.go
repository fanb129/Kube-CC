package responses

import (
	corev1 "k8s.io/api/core/v1"
)

type Pvc struct {
	Name             string `json:"name"`
	Namespace        string `json:"namespace"`
	CreatedAt        string `json:"created_at"`
	Storage          string `json:"storage"`
	StorageClassName string `json:"storage_class_name"`
	AccessModes      string `json:"access_modes"`
	Status           string `json:"status"`
	Volume           string `json:"volume"`
}

type PvcListResponse struct {
	Response
	Length  int   `json:"length"`
	PvcList []Pvc `json:"pvc_list"`
}

type PvcInfo struct {
	Response
	Info corev1.PersistentVolumeClaim `json:"info"`
}
