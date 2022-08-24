package service

import (
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/dao"
)

// CreateIngress 创建ingress
func CreateIngress(name, ns string, label map[string]string, spec v1beta1.IngressSpec) (*v1beta1.Ingress, error) {
	ingress := v1beta1.Ingress{
		TypeMeta:   metav1.TypeMeta{Kind: "Ingress", APIVersion: "extensions/v1beta1"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns, Labels: label},
		Spec:       spec,
	}
	create, err := dao.ClientSet.ExtensionsV1beta1().Ingresses(ns).Create(&ingress)
	if err != nil {
		return nil, err
	}
	return create, nil
}

// GetIngress 获得指定namespace下的ingress
func GetIngress(ns string, label string) (*common.IngressListResponse, error) {
	list, err := dao.ClientSet.ExtensionsV1beta1().Ingresses(ns).List(metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	num := len(list.Items)
	ingressList := make([]common.Ingress, num)
	for i, ing := range list.Items {
		tmp := common.Ingress{
			Name:      ing.Name,
			Namespace: ing.Namespace,
			Rules:     ing.Spec.Rules,
		}
		ingressList[i] = tmp
	}
	return &common.IngressListResponse{
		Response:    common.OK,
		Length:      num,
		IngressList: ingressList,
	}, nil
}

// DeleteIngress 删除指定ingress
func DeleteIngress(name, ns string) (*common.Response, error) {
	err := dao.ClientSet.ExtensionsV1beta1().Ingresses(ns).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &common.OK, nil
}
