package service

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetAPod 获得指定deploy
func GetAPod(name, ns string) (*corev1.Pod, error) {
	get, err := dao.ClientSet.CoreV1().Pods(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return get, nil
}

// GetPod 获得指定namespace下pod
func GetPod(ns string, label string) (*responses.PodListResponse, error) {
	pods, err := dao.ClientSet.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}

	num := len(pods.Items)
	podList := make([]responses.Pod, num)

	for i, pod := range pods.Items {
		tmp := responses.Pod{
			Name:              pod.Name,
			Namespace:         pod.Namespace,
			CreatedAt:         pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
			NodeIp:            pod.Status.HostIP,
			Phase:             pod.Status.Phase,
			ContainerStatuses: pod.Status.ContainerStatuses,
			Uid:               pod.Labels["u_id"],
		}
		podList[i] = tmp
	}
	return &responses.PodListResponse{Response: responses.OK, Length: num, PodList: podList}, nil
}

// DeletePod 删除指定pod
func DeletePod(name, ns string) (*responses.Response, error) {
	err := dao.ClientSet.CoreV1().Pods(ns).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}

func CreatePod(name, ns string, label map[string]string, spec corev1.PodSpec) (*corev1.Pod, error) {
	pod := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels:    label,
		},
		Spec: spec,
	}
	create, err := dao.ClientSet.CoreV1().Pods(ns).Create(context.Background(), &pod, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return create, nil
}

func UpdatePod(pod *corev1.Pod) (*responses.Response, error) {
	_, err := dao.ClientSet.CoreV1().Pods(pod.Name).Update(context.Background(), pod, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}
