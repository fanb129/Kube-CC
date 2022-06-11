package kube

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var clientset *kubernetes.Clientset

func init() {
	config, err := clientcmd.BuildConfigFromFlags("", "E:\\go\\k8s_deploy_gin\\conf\\kube.config")
	if err != nil {
		panic(err.Error())
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}
