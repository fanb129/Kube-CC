package service

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/dao"
	"log"
)

// GetNs 获取所有namespace （default，kube-node-lease，kube-public，kube-system 系统自带除外）
func GetNs() (*common.NsListResponse, error) {
	namespace, err := dao.ClientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	num := len(namespace.Items)
	namespaceList := make([]interface{}, 0, num)
	for _, ns := range namespace.Items {
		//if ns.Name == "default" || ns.Name == "kube-node-lease" || ns.Name == "kube-public" || ns.Name == "kube-system" {
		//	continue
		//}
		tmpMap := map[string]interface{}{
			"name":     ns.Name,
			"status":   ns.Status.Phase,
			"createAt": ns.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		namespaceList = append(namespaceList, tmpMap)
	}
	return &common.NsListResponse{Response: common.OK, Length: num, NsList: namespaceList}, nil
}
