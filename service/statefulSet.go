package service

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetStatefulSet 获得指定statefulSet
func GetStatefulSet(name, ns string) (*appsv1.StatefulSet, error) {
	get, err := dao.ClientSet.AppsV1().StatefulSets(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return get, nil
}

// CreateStatefulSet 创建自定义控制器
func CreateStatefulSet(name, ns string, label map[string]string, spec appsv1.StatefulSetSpec) (*appsv1.StatefulSet, error) {
	rs := appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{APIVersion: "apps/v1", Kind: "StatefulSet"},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			//ServiceName: serviceName,
			Namespace: ns,
			Labels:    label,
		},
		Spec: spec,
	}
	create, err := dao.ClientSet.AppsV1().StatefulSets(ns).Create(context.Background(), &rs, metav1.CreateOptions{})
	return create, err
}

// ListStatefulSet  获得指定namespace下的控制器
func ListStatefulSet(ns, label string) ([]appsv1.StatefulSet, error) {
	list, err := dao.ClientSet.AppsV1().StatefulSets(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// DeleteStatefulSet 删除指定namespace的控制器
func DeleteStatefulSet(name, ns string) (*responses.Response, error) {
	err := dao.ClientSet.AppsV1().StatefulSets(ns).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}

// UpdateStatefulSet 更新statefulSet
func UpdateStatefulSet(name, ns string, spec appsv1.StatefulSetSpec) (*appsv1.StatefulSet, error) {
	set, err := GetStatefulSet(name, ns)
	if err != nil {
		return nil, err
	}
	update, err := dao.ClientSet.AppsV1().StatefulSets(ns).Update(context.Background(), set, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return update, nil
}
