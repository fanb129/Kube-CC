package kube

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func GetNs() (int, []interface{}) {
	namespace, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	num := len(namespace.Items)
	namespaceList := make([]interface{}, 0, num)
	for _, ns := range namespace.Items {
		tmpMap := map[string]interface{}{
			"name":     ns.Name,
			"status":   ns.Status.Phase,
			"createAt": ns.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		namespaceList = append(namespaceList, tmpMap)
	}
	return num, namespaceList
}

//func CreateNs(name string) bool{
//
//}