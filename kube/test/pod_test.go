package test

import (
	"fmt"
	"k8s_deploy_gin/kube"
	"testing"
)

func TestGetPod(t *testing.T) {
	ns := "kube-system"
	num, pods := kube.GetPod(ns)
	fmt.Println(num)

	for _, pod := range pods {
		fmt.Println(pod)
	}
}
