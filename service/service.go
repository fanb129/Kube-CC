package service

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/conf"
	"k8s_deploy_gin/dao"
)

// CreateService 创建自定义服务
func CreateService(name, ns string, label map[string]string, spec corev1.ServiceSpec) (*corev1.Service, error) {
	service := corev1.Service{
		TypeMeta:   metav1.TypeMeta{Kind: "service", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns, Labels: label},
		Spec:       spec,
	}
	create, err := dao.ClientSet.CoreV1().Services(ns).Create(&service)
	if err != nil {
		return nil, err
	}
	return create, err
}

// GetService 获得指定ns下的service
func GetService(ns string, label string) (*common.ServiceListResponse, error) {
	list, err := dao.ClientSet.CoreV1().Services(ns).List(metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	num := len(list.Items)
	serviceList := make([]common.Service, num)
	for i, sc := range list.Items {
		tmp := common.Service{
			Name:      sc.Name,
			Namespase: sc.Namespace,
			Ports:     sc.Spec.Ports,
			SshPwd:    conf.SshPwd,
		}
		serviceList[i] = tmp
	}
	return &common.ServiceListResponse{Response: common.OK, Length: num, ServiceList: serviceList}, nil
}

// DeleteService 删除指定service
func DeleteService(name, ns string) (*common.Response, error) {
	err := dao.ClientSet.CoreV1().Services(ns).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &common.OK, nil
}
