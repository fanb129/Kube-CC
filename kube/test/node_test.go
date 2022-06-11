package test

import (
	"fmt"
	"k8s_deploy_gin/kube"
	"testing"
)

func TestGetNode(t *testing.T) {
	num, list := kube.GetNode()
	fmt.Println(num)
	for _, node := range list {
		fmt.Println(node)
	}
}
