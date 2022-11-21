package service

import (
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/conf"
	"k8s_deploy_gin/dao"
	"k8s_deploy_gin/service/ssh"
	"sync"
)

// master创建token
func createToken() (string, error) {
	config := ssh.Config{
		Host:     conf.MasterInfo.Host,
		Port:     conf.MasterInfo.Port,
		User:     conf.MasterInfo.User,
		Type:     ssh.TypePassword,
		Password: conf.MasterInfo.Password,
	}
	newSsh, err := ssh.NewSsh(config)
	if err != nil {
		zap.S().Errorln(err)
		return "", err
	}
	defer newSsh.CloseClient()
	r, err := newSsh.SendCmd("kubeadm token create --print-join-command 2> /dev/null")
	if err != nil {
		zap.S().Errorln(err)
		return "", err
	}
	zap.S().Debug(r)
	return r, nil
}

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
			CreatedAt:      node.CreationTimestamp.Format("2006-01-02 15:04:05"),
			OsImage:        node.Status.NodeInfo.OSImage,
			KubeletVersion: node.Status.NodeInfo.KubeletVersion,
			CPU:            node.Status.Allocatable.Cpu().String() + " / " + node.Status.Capacity.Cpu().String(),
			Memory:         node.Status.Allocatable.Memory().String() + " / " + node.Status.Capacity.Memory().String(),
		}
		nodeList[i] = tmp
	}
	return &common.NodeListResponse{Response: common.OK, Length: num, NodeList: nodeList}, nil
}

// CreateNode 添加node
func CreateNode(configs []ssh.Config) (*common.Response, error) {
	//node := corev1.Node{
	//	ObjectMeta: metav1.ObjectMeta{
	//		Name: name,
	//	},
	//}
	//create, err := dao.ClientSet.CoreV1().Nodes().Create(&node)
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println(create)
	//return &common.OK, nil
	token, err := createToken()
	if err != nil {
		zap.S().Errorln(err)
		return nil, err
	}
	// 使用协程，并行批量添加
	group := sync.WaitGroup{}
	group.Add(len(configs))
	for _, config := range configs {
		go func(config ssh.Config) {
			zap.S().Info(config)
			newSsh, err := ssh.NewSsh(config)
			defer newSsh.CloseClient()
			if err != nil {
				zap.S().Errorln(err)
				group.Done()
			}
			// 在join之前，先reset
			reset := "echo y|kubeadm reset"
			if _, err = newSsh.SendCmd(reset + "&&" + token); err != nil {
				zap.S().Errorln(err)
			}
			group.Done()
		}(config)
	}
	group.Wait()
	return &common.OK, nil
}

// DeleteNode 删除node节点
func DeleteNode(name string) (*common.Response, error) {
	err := dao.ClientSet.CoreV1().Nodes().Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &common.OK, nil
}
