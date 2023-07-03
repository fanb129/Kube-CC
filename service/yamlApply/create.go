package yamlApply

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"Kube-CC/service"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NamespaceCreate(ns *corev1.Namespace) (*responses.Response, error) {
	if _, err := dao.ClientSet.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{}); err != nil {
		return nil, err
	}
	return &responses.OK, nil
}

// DeployCreate  deploy的创建
func DeployCreate(deploy *appsv1.Deployment) (*responses.Response, error) {
	name := deploy.Name
	ns := deploy.Namespace
	labels := deploy.Labels
	if labels == nil {
		labels = make(map[string]string)
	}
	// 获取namespace，提取出uid的label
	get, err := dao.ClientSet.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
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
	return &responses.OK, nil
}

// StatefulSetCreate  deploy的创建
func StatefulSetCreate(sts *appsv1.StatefulSet) (*responses.Response, error) {
	name := sts.Name
	ns := sts.Namespace
	labels := sts.Labels
	if labels == nil {
		labels = make(map[string]string)
	}
	// 获取namespace，提取出uid的label
	get, err := dao.ClientSet.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	uid, ok := get.Labels["u_id"]
	if ok {
		labels["u_id"] = uid //sts的label
		sts.Spec.Template.Labels["u_id"] = uid
	}
	if _, err := service.CreateStatefulSet(name, ns, labels, sts.Spec); err != nil {
		return nil, err
	}
	return &responses.OK, nil
}

// ServiceCreate service的更新或者创建
func ServiceCreate(svc *corev1.Service) (*responses.Response, error) {
	name := svc.Name
	ns := svc.Namespace
	labels := svc.Labels
	if labels == nil {
		labels = make(map[string]string)
	}
	// 获取namespace，提取出uid的label
	get, err := dao.ClientSet.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
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
	return &responses.OK, nil
}

// PodCreate 创建pod
func PodCreate(pod *corev1.Pod) (*responses.Response, error) {
	name := pod.Name
	ns := pod.Namespace
	labels := pod.Labels
	if labels == nil {
		labels = make(map[string]string)
	}
	// 获取namespace，提取出uid的label
	get, err := dao.ClientSet.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
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
	return &responses.OK, nil
}

// PodCreate 创建job
func JobCreate(job *corev1.Pod) (*responses.Response, error) {
	name := job.Name
	ns := job.Namespace
	labels := job.Labels
	if labels == nil {
		labels = make(map[string]string)
	}
	// 获取namespace，提取出uid的label
	get, err := dao.ClientSet.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	uid, ok := get.Labels["u_id"]
	if ok {
		labels["u_id"] = uid //service的label
	}
	if _, err := service.CreateJob(name, ns, labels, job.Spec); err != nil {
		return nil, err
	}
	return &responses.OK, nil
}
