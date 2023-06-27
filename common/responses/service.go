package responses

import corev1 "k8s.io/api/core/v1"

type Service struct {
	Name      string               `json:"name"`
	Namespace string               `json:"namespace"`
	CreatedAt string               `json:"created_at"`
	ClusterIP string               `json:"cluster_ip"`
	Type      corev1.ServiceType   `json:"type"`
	Ports     []corev1.ServicePort `json:"ports"`
	SshPwd    string               `json:"ssh_pwd,omitempty"`
	Uid       string               `json:"u_id"`
}

// ServiceListResponse 服务返回结果
type ServiceListResponse struct {
	Response
	Length      int       `json:"length"`
	ServiceList []Service `json:"service_list"`
}

type ServiceInfo struct {
	Response
	Info corev1.Service `json:"info"`
}
