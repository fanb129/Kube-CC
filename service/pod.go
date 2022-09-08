package service

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/dao"
)

// GetPod 获得指定namespace下pod
func GetPod(ns string, label string) (*common.PodListResponse, error) {
	pods, err := dao.ClientSet.CoreV1().Pods(ns).List(metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}

	num := len(pods.Items)
	podList := make([]common.Pod, num)

	for i, pod := range pods.Items {
		tmp := common.Pod{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Ready:     pod.Status.ContainerStatuses[0].Ready,
			Status:    pod.Status.Conditions[0].Status,
			NodeIp:    pod.Status.HostIP,
		}
		podList[i] = tmp
	}
	return &common.PodListResponse{Response: common.OK, Length: num, PodList: podList}, nil
}

// DeletePod 删除指定pod
func DeletePod(name, ns string) (*common.Response, error) {
	err := dao.ClientSet.CoreV1().Pods(ns).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &common.OK, nil
}

//func AddNPod(ns string){
//	pod := corev1.Pod{
//		TypeMeta:metav1.TypeMeta{
//			Kind: "",
//			APIVersion: "v1",
//		},
//		ObjectMeta:metav1.ObjectMeta{
//
//		},
//		Spec: corev1.PodSpec{
//
//		},
//	}
//	create, err := dao.ClientSet.CoreV1().Pods(ns).Create(&pod)
//}
