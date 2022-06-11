package kube

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func GetPod(ns string) (int, []interface{}) {
	pods, err := clientset.CoreV1().Pods(ns).List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
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
	return num, podList
}
