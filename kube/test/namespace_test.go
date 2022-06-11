package test

import (
	"fmt"
	"k8s_deploy_gin/kube"
	"testing"
)

func TestGetNs(t *testing.T) {
	num, ns := kube.GetNs()
	fmt.Println(num)
	for _, n := range ns {
		fmt.Printf("%v\n", n)
	}
}
