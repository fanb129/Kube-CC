package service

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/dao"
	"log"
)

// GetNode 获得所有node
func GetNode() (*common.NodeListResponse, error) {
	nodes, err := dao.ClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	num := len(nodes.Items)
	nodeList := make([]interface{}, 0, num)
	//遍历所有node实列
	for _, node := range nodes.Items {
		tmpMap := map[string]interface{}{
			"name":    node.Name,
			"ip":      node.Status.Addresses[0].Address,
			"status":  node.Status.Conditions[3].Status,
			"creatAt": node.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		nodeList = append(nodeList, tmpMap)
	}
	return &common.NodeListResponse{Response: common.OK, Length: num, NodeList: nodeList}, nil
}
