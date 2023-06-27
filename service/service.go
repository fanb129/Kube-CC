package service

import (
	"Kube-CC/common/responses"
	"Kube-CC/conf"
	"Kube-CC/dao"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetAService 获得指定deploy
func GetAService(name, ns string) (*corev1.Service, error) {
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

// GetService 获得指定ns下的service
func GetService(ns string, label string) (*responses.ServiceListResponse, error) {
	list, err := dao.ClientSet.CoreV1().Services(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	num := len(list.Items)
	serviceList := make([]responses.Service, num)
	for i, sc := range list.Items {
		tmp := responses.Service{
			Name:      sc.Name,
			Namespace: sc.Namespace,
			CreatedAt: sc.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Ports:     sc.Spec.Ports,
			ClusterIP: sc.Spec.ClusterIP,
			Type:      sc.Spec.Type,
			SshPwd:    conf.SshPwd,
			Uid:       sc.Labels["u_id"],
		}
		serviceList[i] = tmp
	}
	return &responses.ServiceListResponse{Response: responses.OK, Length: num, ServiceList: serviceList}, nil
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
func UpdateService(service *corev1.Service) (*responses.Response, error) {
	_, err := dao.ClientSet.CoreV1().Services(service.Namespace).Update(context.Background(), service, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}
