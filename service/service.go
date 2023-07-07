package service

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetService 获得指定
func GetService(name, ns string) (*corev1.Service, error) {
	get, err := dao.ClientSet.CoreV1().Services(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return get, nil
}

// CreateService 创建自定义服务
func CreateService(name, ns string, label map[string]string, spec corev1.ServiceSpec) (*corev1.Service, error) {
	service := corev1.Service{
		TypeMeta:   metav1.TypeMeta{Kind: "service", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns, Labels: label},
		Spec:       spec,
	}
	create, err := dao.ClientSet.CoreV1().Services(ns).Create(context.Background(), &service, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return create, err
}

// ListService 获得指定ns下的service
func ListService(ns string, label string) ([]corev1.Service, error) {
	list, err := dao.ClientSet.CoreV1().Services(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// DeleteService 删除指定service
func DeleteService(name, ns string) (*responses.Response, error) {
	err := dao.ClientSet.CoreV1().Services(ns).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}

// UpdateService 更新service
func UpdateService(name, ns string, spec corev1.ServiceSpec) (*corev1.Service, error) {
	service, err := GetService(name, ns)
	if err != nil {
		return nil, err
	}
	service.Spec = spec
	update, err := dao.ClientSet.CoreV1().Services(service.Namespace).Update(context.Background(), service, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return update, nil
}
