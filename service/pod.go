package service

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/dao"
	"log"
)

// GetPod 获得指定namespace下pod
func GetPod(ns string) (*common.PodListResponse, error) {
	pods, err := dao.ClientSet.CoreV1().Pods(ns).List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	num := len(pods.Items)
	podList := make([]interface{}, 0, num)

	for _, pod := range pods.Items {
		tmpMap := map[string]interface{}{
			"name":      pod.Name,
			"namespace": pod.Namespace,
			"ready":     pod.Status.Conditions[0].Status,
			"status":    pod.Status.ContainerStatuses[0].Ready,
			"nodeIP":    pod.Status.HostIP,
		}
		podList = append(podList, tmpMap)
	}
	return &common.PodListResponse{Response: common.OK, Length: num, PodList: podList}, nil
}
