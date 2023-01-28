package dao

import (
	"Kube-CC/conf"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var ClientSet *kubernetes.Clientset

func InitKube() (err error) {
	config, err := clientcmd.BuildConfigFromFlags("", conf.KubeConfig)
	if err != nil {
		panic(err.Error())
		return err
	}
	ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
		return err
	}
	return
}
