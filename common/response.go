package common

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
)

type Deploy struct {
	Name              string `json:"name"`
	Namespace         string `json:"namespace"`
	CreatedAt         string `json:"created_at"`
	Replicas          int32  `json:"replicas"`
	UpdatedReplicas   int32  `json:"updated_replicas"`
	ReadyReplicas     int32  `json:"ready_replicas"`
	AvailableReplicas int32  `json:"available_replicas"`
	Uid               string `json:"u_id"`
}

type Ns struct {
	Name      string                `json:"name"`
	Status    corev1.NamespacePhase `json:"status"`
	CreatedAt string                `json:"created_at"`
	Username  string                `json:"username"`
	Nickname  string                `json:"nickname"`
	Uid       uint                  `json:"u_id"`
}

type Node struct {
	Name           string                 `json:"name"`
	Ip             string                 `json:"ip"`
	Ready          corev1.ConditionStatus `json:"ready"`
	CreatedAt      string                 `json:"created_at"`
	OsImage        string                 `json:"os_image"`
	KubeletVersion string                 `json:"kubelet_version"`
	CPU            string                 `json:"cpu"`
	Memory         string                 `json:"memory"`
}

type Pod struct {
	Name              string                   `json:"name"`
	Namespace         string                   `json:"namespace"`
	CreatedAt         string                   `json:"created_at"`
	Phase             corev1.PodPhase          `json:"phase"`
	NodeIp            string                   `json:"node_ip"`
	ContainerStatuses []corev1.ContainerStatus `json:"container_statuses"`
	Uid               string                   `json:"u_id"`
}

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

type Ingress struct {
	Name      string                `json:"name"`
	Namespace string                `json:"namespace"`
	CreatedAt string                `json:"created_at"`
	Rules     []v1beta1.IngressRule `json:"rules"`
	Uid       uint                  `json:"u_id"`
}

type Spark struct {
	Name           string                `json:"name"`
	Uid            uint                  `json:"u_id"`
	Username       string                `json:"username"`
	Nickname       string                `json:"nickname"`
	CreatedAt      string                `json:"created_at"`
	Status         corev1.NamespacePhase `json:"status"`
	PodList        []Pod                 `json:"pod_list"`
	DeployList     []Deploy              `json:"deploy_list"`
	ServiceList    []Service             `json:"service_list"`
	IngressList    []Ingress             `json:"ingress_list"`
	MasterReplicas int32                 `json:"master_replicas"`
	WorkerReplicas int32                 `json:"worker_replicas"`
}
type Linux struct {
	Name        string                `json:"name"`
	Uid         uint                  `json:"u_id"`
	Username    string                `json:"username"`
	Nickname    string                `json:"nickname"`
	Status      corev1.NamespacePhase `json:"status"`
	CreatedAt   string                `json:"created_at"`
	PodList     []Pod                 `json:"pod_list"`
	DeployList  []Deploy              `json:"deploy_list"`
	ServiceList []Service             `json:"service_list"`
}
type Hadoop struct {
	Name               string                `json:"name"`
	Uid                uint                  `json:"u_id"`
	Username           string                `json:"username"`
	Nickname           string                `json:"nickname"`
	Status             corev1.NamespacePhase `json:"status"`
	CreatedAt          string                `json:"created_at"`
	PodList            []Pod                 `json:"pod_list"`
	DeployList         []Deploy              `json:"deploy_list"`
	ServiceList        []Service             `json:"service_list"`
	HdfsMasterReplicas int32                 `json:"hdfs_master_replicas"`
	DatanodeReplicas   int32                 `json:"datanode_replicas"`
	YarnMasterReplicas int32                 `json:"yarn_master_replicas"`
	YarnNodeReplicas   int32                 `json:"yarn_node_replicas"`
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

type DeployInfo struct {
	Response
	Info appsv1.Deployment `json:"info"`
}

type ServiceInfo struct {
	Response
	Info corev1.Service `json:"info"`
}

type PodInfo struct {
	Response
	Info corev1.Pod `json:"info"`
}
