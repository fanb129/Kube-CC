package common

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
)

type Deploy struct {
	Name          string `json:"name"`
	Namespace     string `json:"namespace"`
	Replicas      int32  `json:"replicas"`
	ReadyReplicas int32  `json:"ready_replicas"`
}

type Ns struct {
	Name     string                `json:"name"`
	Status   corev1.NamespacePhase `json:"status"`
	CreateAt string                `json:"create_at"`
}

type Node struct {
	Name     string                 `json:"name"`
	Ip       string                 `json:"ip"`
	Status   corev1.ConditionStatus `json:"status"`
	CreateAt string                 `json:"create_at"`
}

type Pod struct {
	Name      string                 `json:"name"`
	Namespase string                 `json:"namespase"`
	Ready     bool                   `json:"ready"`
	Status    corev1.ConditionStatus `json:"status"`
	NodeIp    string                 `json:"node_ip"`
}

type Service struct {
	Name      string               `json:"name"`
	Namespase string               `json:"namespase"`
	Ports     []corev1.ServicePort `json:"ports"`
	SshPwd    string               `json:"ssh_pwd,omitempty"`
}

type Ingress struct {
	Name      string                `json:"name"`
	Namespace string                `json:"namespace"`
	Rules     []v1beta1.IngressRule `json:"rules"`
}

type Spark struct {
	Name        string    `json:"name"`
	Uid         uint      `json:"u_id"`
	PodList     []Pod     `json:"pod_list"`
	DeployList  []Deploy  `json:"deploy_list"`
	ServiceList []Service `json:"service_list"`
	IngressList []Ingress `json:"ingress_list"`
}
type Linux struct {
	Name        string    `json:"name"`
	Uid         uint      `json:"u_id"`
	PodList     []Pod     `json:"pod_list"`
	DeployList  []Deploy  `json:"deploy_list"`
	ServiceList []Service `json:"service_list"`
}
type Hadoop struct {
	Name        string    `json:"name"`
	Uid         uint      `json:"u_id"`
	PodList     []Pod     `json:"pod_list"`
	DeployList  []Deploy  `json:"deploy_list"`
	ServiceList []Service `json:"service_list"`
}
type UserInfo struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Role      uint   `json:"role"`
	Avatar    string `json:"avatar"`
}

type Response struct {
	StatusCode int    `json:"code"`
	StatusMsg  string `json:"msg,omitempty"`
}

type ResponseOfValidator struct {
	StatusCode int         `json:"code"`
	StatusMsg  interface{} `json:"msg,omitempty"`
}

func ValidatorResponse(err error) ResponseOfValidator {
	return ResponseOfValidator{
		-1,
		translate(err),
	}
}

var OK = Response{StatusCode: 1, StatusMsg: "success"}
var NoRole = Response{StatusCode: -1, StatusMsg: "权限不够"}
var NoToken = Response{StatusCode: 50008, StatusMsg: "No Token"}
var TokenExpired = Response{StatusCode: 50014, StatusMsg: "token过期"}
var NoUid = Response{StatusCode: -1, StatusMsg: "Uid获取失败"}

type LoginResponse struct {
	Response
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

type UserInfoResponse struct {
	Response
	UserInfo UserInfo `json:"user_info"`
}
type UserListResponse struct {
	Response
	Page     int        `json:"page"`
	Total    int        `json:"total"`
	UserList []UserInfo `json:"user_list"`
}

type NodeListResponse struct {
	Response
	Length   int    `json:"length"`
	NodeList []Node `json:"node_list"`
}
type NsListResponse struct {
	Response
	Length int  `json:"length"`
	NsList []Ns `json:"ns_list"`
}
type PodListResponse struct {
	Response
	Length  int   `json:"length"`
	PodList []Pod `json:"pod_list"`
}

// DeployListResponse pod控制器返回结果
type DeployListResponse struct {
	Response
	Length     int      `json:"length"`
	DeployList []Deploy `json:"deploy_list"`
}

// ServiceListResponse 服务返回结果
type ServiceListResponse struct {
	Response
	Length      int       `json:"length"`
	ServiceList []Service `json:"service_list"`
}

type SparkListResponse struct {
	Response
	Length    int     `json:"length"`
	SparkList []Spark `json:"spark_list"`
}
type LinuxListResponse struct {
	Response
	Image     string  `json:"image"`
	Length    int     `json:"length"`
	LinuxList []Linux `json:"linux_list"`
}
type IngressListResponse struct {
	Response
	Length      int       `json:"length"`
	IngressList []Ingress `json:"ingress_list"`
}
type HadoopListResponse struct {
	Response
	Length     int      `json:"length"`
	HadoopList []Hadoop `json:"hadoop_list"`
}
