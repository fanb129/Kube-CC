package yamlApply

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/dao"
	"k8s_deploy_gin/service"
)

func NamespaceCreate(ns *corev1.Namespace) (*common.Response, error) {
	if _, err := dao.ClientSet.CoreV1().Namespaces().Create(ns); err != nil {
		return nil, err
	}
	return &common.OK, nil
}

// DeployCreate  deploy的创建
func DeployCreate(deploy *appsv1.Deployment) (*common.Response, error) {
	name := deploy.Name
	ns := deploy.Namespace
	labels := deploy.Labels
	// 获取namespace，提取出uid的label
	get, err := dao.ClientSet.CoreV1().Namespaces().Get(ns, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	uid, ok := get.Labels["u_id"]
	if ok {
		labels["u_id"] = uid                      //deploy的label
		deploy.Spec.Template.Labels["u_id"] = uid // pod的label
	}
	if _, err := service.CreateDeploy(name, ns, labels, deploy.Spec); err != nil {
		return nil, err
	}
	return &common.OK, nil
}

// ServiceCreate service的更新或者创建
func ServiceCreate(svc *corev1.Service) (*common.Response, error) {
	name := svc.Name
	ns := svc.Namespace
	labels := svc.Labels
	// 获取namespace，提取出uid的label
	get, err := dao.ClientSet.CoreV1().Namespaces().Get(ns, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	uid, ok := get.Labels["u_id"]
	if ok {
		labels["u_id"] = uid //service的label
	}
	if _, err := service.CreateService(name, ns, labels, svc.Spec); err != nil {
		return nil, err
	}
	return &common.OK, nil
}

// PodCreate 创建pod
func PodCreate(pod *corev1.Pod) (*common.Response, error) {
	name := pod.Name
	ns := pod.Namespace
	labels := pod.Labels
	// 获取namespace，提取出uid的label
	get, err := dao.ClientSet.CoreV1().Namespaces().Get(ns, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	uid, ok := get.Labels["u_id"]
	if ok {
		labels["u_id"] = uid //service的label
	}
	if _, err := service.CreatePod(name, ns, labels, pod.Spec); err != nil {
		return nil, err
	}
	return &common.OK, nil
}
