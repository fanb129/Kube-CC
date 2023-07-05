package service

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPod 获得指定deploy
func GetPod(name, ns string) (*corev1.Pod, error) {
	get, err := dao.ClientSet.CoreV1().Pods(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return get, nil
}

// ListPod   获得指定namespace下pod
func ListPod(ns string, label string) ([]corev1.Pod, error) {
	list, err := dao.ClientSet.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// DeletePod 删除指定pod
func DeletePod(name, ns string) (*responses.Response, error) {
	err := dao.ClientSet.CoreV1().Pods(ns).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}

// ListDeployPod   获得指定namespace下pod
func ListDeployPod(ns string, label string) ([]responses.DeployPod, error) {
	list, err := dao.ClientSet.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	num := len(list.Items)
	podList := make([]responses.DeployPod, num)
	for i, pod := range list.Items {
		podList[i] = responses.DeployPod{
			Name:   pod.Name,
			Phase:  string(pod.Status.Phase),
			PodIP:  pod.Status.PodIP,
			HostIP: pod.Status.HostIP,
		}
	}
	return podList, nil
}

// ListStatefulSetPod   获得指定namespace下pod
func ListStatefulSetPod(ns string, label string) ([]responses.StsPod, error) {
	list, err := dao.ClientSet.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	num := len(list.Items)
	podList := make([]responses.StsPod, num)
	for i, pod := range list.Items {
		podList[i] = responses.StsPod{
			Name:   pod.Name,
			Phase:  string(pod.Status.Phase),
			PodIP:  pod.Status.PodIP,
			HostIP: pod.Status.HostIP,
		}
	}
	return podList, nil
}
