package dao

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s_deploy_gin/conf"
)

var ClientSet *kubernetes.Clientset

func InitKube() {
	config, err := clientcmd.BuildConfigFromFlags("", conf.KubeConfig)
	if err != nil {
		panic(err.Error())
	}
	ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}
