package service

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/dao"
	"strconv"
)

// GetNs 获取所有namespace
func GetNs(label string) (*common.NsListResponse, error) {
	namespace, err := dao.ClientSet.CoreV1().Namespaces().List(metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}

	num := len(namespace.Items)
	namespaceList := make([]common.Ns, num)
	for i, ns := range namespace.Items {
		//if ns.Name == "default" || ns.Name == "kube-node-lease" || ns.Name == "kube-public" || ns.Name == "kube-system" {
		//	continue
		//}
		uid, err := strconv.Atoi(ns.Labels["u_id"])
		username := ""
		nickname := ""
		if err == nil {
			user, err := dao.GetUserById(uint(uid))
			if err == nil {
				username = user.Username
				nickname = user.Nickname
			}
		}

		tmp := common.Ns{
			Name:     ns.Name,
			Status:   ns.Status.Phase,
			CreateAt: ns.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Username: username,
			Nickname: nickname,
		}
		namespaceList[i] = tmp
	}
	return &common.NsListResponse{Response: common.OK, Length: num, NsList: namespaceList}, nil
}

// CreateNs 新建属于指定用户的namespace，u_id == 0 则不添加标签
func CreateNs(name string, label map[string]string) (*common.Response, error) {
	ns := v1.Namespace{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Namespace"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Labels: label},
	}
	_, err := dao.ClientSet.CoreV1().Namespaces().Create(&ns)
	if err != nil {
		return nil, err
	}
	return &common.OK, nil
}

// DeleteNs 删除指定namespace
func DeleteNs(name string) (*common.Response, error) {
	err := dao.ClientSet.CoreV1().Namespaces().Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &common.OK, nil
}
