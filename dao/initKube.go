package dao

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s_deploy_gin/conf"
)

var ClientSet *kubernetes.Clientset

func InitKube() (err error) {
	config, err := clientcmd.BuildConfigFromFlags("", conf.KubeConfig)
	if err != nil {
		panic(err.Error())
		return
	}
	ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
		return
	}
	return
}
