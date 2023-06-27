package responses

import "k8s.io/api/extensions/v1beta1"

type Ingress struct {
	Name      string                `json:"name"`
	Namespace string                `json:"namespace"`
	CreatedAt string                `json:"created_at"`
	Rules     []v1beta1.IngressRule `json:"rules"`
	Uid       uint                  `json:"u_id"`
}

type IngressListResponse struct {
	Response
	Length      int       `json:"length"`
	IngressList []Ingress `json:"ingress_list"`
}
