package service

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/dao"
)

// GetNode 获得所有node
func GetNode(label string) (*common.NodeListResponse, error) {
	nodes, err := dao.ClientSet.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}

	num := len(nodes.Items)
	nodeList := make([]common.Node, num)
	//遍历所有node实列
	for i, node := range nodes.Items {
		tmp := common.Node{
			Name:     node.Name,
			Ip:       node.Status.Addresses[0].Address,
			Status:   node.Status.Conditions[3].Status,
			CreateAt: node.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		nodeList[i] = tmp
	}
	return &common.NodeListResponse{Response: common.OK, Length: num, NodeList: nodeList}, nil
}
