package service

import (
	"Kube-CC/common/responses"
	"Kube-CC/conf"
	"Kube-CC/dao"
	"Kube-CC/service/ssh"
	"context"
	"errors"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
func GetNode(label string) (*responses.NodeListResponse, error) {
	nodes, err := dao.ClientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}

	num := len(nodes.Items)
	nodeList := make([]responses.Node, num)
	//遍历所有node实列
	for i, node := range nodes.Items {
		cpu := node.Status.Capacity.Cpu()
		usedCpu := cpu.DeepCopy()
		usedCpu.Sub(*node.Status.Allocatable.Cpu())
		memory := node.Status.Capacity.Memory()
		usedMemory := memory.DeepCopy()
		usedMemory.Sub(*node.Status.Allocatable.Memory())
		storage := node.Status.Capacity.StorageEphemeral()
		usedStorage := storage.DeepCopy()
		usedStorage.Sub(*node.Status.Allocatable.StorageEphemeral())
		pvc := node.Status.Capacity.Storage()
		usedPvc := pvc.DeepCopy()
		usedPvc.Sub(*node.Status.Allocatable.Storage())
		gpu := node.Status.Capacity.Name(GpuShare, resource.BinarySI)
		usedGpu := gpu.DeepCopy()
		usedGpu.Sub(*node.Status.Allocatable.Name(GpuShare, resource.BinarySI))
		tmp := responses.Node{
			Name:           node.Name,
			Ip:             node.Status.Addresses[0].Address,
			Ready:          node.Status.Conditions[len(node.Status.Conditions)-1].Status,
			CreatedAt:      node.CreationTimestamp.Format("2006-01-02 15:04:05"),
			OsImage:        node.Status.NodeInfo.OSImage,
			KubeletVersion: node.Status.NodeInfo.KubeletVersion,
			Resources: responses.Resources{
				Cpu:         cpu.String(),
				UsedCpu:     usedCpu.String(),
				Memory:      memory.String(),
				UsedMemory:  usedMemory.String(),
				Storage:     storage.String(),
				UsedStorage: usedStorage.String(),
				PVC:         pvc.String(),
				UsedPVC:     usedPvc.String(),
				GPU:         gpu.String(),
				UsedGPU:     usedGpu.String(),
			},
		}
		nodeList[i] = tmp
	}
	return &responses.NodeListResponse{Response: responses.OK, Length: num, NodeList: nodeList}, nil
}

// CreateNode 添加node
func CreateNode(configs []ssh.Config) (*responses.Response, error) {
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
	//return &responses.OK, nil
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
	return &responses.OK, nil
}

// DeleteNode 删除node节点
func DeleteNode(name string) (*responses.Response, error) {
	get, err := dao.ClientSet.CoreV1().Nodes().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if get.Status.Addresses[0].Address == conf.MasterInfo.Host {
		return nil, errors.New("不允许删除master")
	}
	err = dao.ClientSet.CoreV1().Nodes().Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}
