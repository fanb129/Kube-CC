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
			Name:           node.Name,
			Ip:             node.Status.Addresses[0].Address,
			Ready:          node.Status.Conditions[len(node.Status.Conditions)-1].Status,
			CreateAt:       node.CreationTimestamp.Format("2006-01-02 15:04:05"),
			OsImage:        node.Status.NodeInfo.OSImage,
			KubeletVersion: node.Status.NodeInfo.KubeletVersion,
			//CPU:            strconv.Itoa(int(node.Status.Capacity.Cpu().Value())),
			CPU:    node.Status.Allocatable.Cpu().String() + " / " + node.Status.Capacity.Cpu().String(),
			Memory: node.Status.Allocatable.Memory().String() + " / " + node.Status.Capacity.Memory().String(),
		}
		nodeList[i] = tmp
	}
	return &common.NodeListResponse{Response: common.OK, Length: num, NodeList: nodeList}, nil
}

//func SshNode(name,ns string){
//	dao.ClientSet.CoreV1().Pods(ns).Get(name,metav1.GetOptions{}).
//}
