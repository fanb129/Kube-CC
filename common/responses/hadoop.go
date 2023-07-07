package responses

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type Hadoop struct {
	Ns
	PodList     []corev1.Pod        `json:"pod_list"`
	DeployList  []appsv1.Deployment `json:"deploy_list"`
	ServiceList []corev1.Service    `json:"service_list"`
}

type HadoopListResponse struct {
	Response
	Length     int      `json:"length"`
	HadoopList []Hadoop `json:"hadoop_list"`
}
