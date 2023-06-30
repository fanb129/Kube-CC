package service

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetAStatefulSet 获得指定statefulSet
func GetAStatefulSet(name, ns string) (*appsv1.StatefulSet, error) {
	get, err := dao.ClientSet.AppsV1().StatefulSets(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return get, nil
}

// CreateStatefulSet 创建自定义控制器
func CreateStatefulSet(name, ns string, spec appsv1.StatefulSetSpec) (*appsv1.StatefulSet, error) {
	rs := appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{APIVersion: "apps/v1", Kind: "StatefulSet"},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			//ServiceName: serviceName,
			Namespace: ns,
			//Labels:    label,
		},
		Spec: spec,
	}
	create, err := dao.ClientSet.AppsV1().StatefulSets(ns).Create(context.Background(), &rs, metav1.CreateOptions{})
	return create, err
}

// GetStatefulSet 获得指定namespace下的控制器
func GetStatefulSet(ns, label string) (*responses.StatefulSetListResponse, error) {
	list, err := dao.ClientSet.AppsV1().StatefulSets(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	num := len(list.Items)
	stslist := make([]responses.StatefulSet, num)
	for i, statefulSet := range list.Items {
		tmp := responses.StatefulSet{
			Name:            statefulSet.Name,
			Namespace:       statefulSet.Namespace,
			CreatedAt:       statefulSet.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Replicas:        statefulSet.Status.Replicas,
			UpdatedReplicas: statefulSet.Status.UpdatedReplicas,
			ReadyReplicas:   statefulSet.Status.ReadyReplicas,
			CurrentReplicas: statefulSet.Status.CurrentReplicas,
			CurrentRevision: statefulSet.Status.CurrentRevision,
			Uid:             statefulSet.Labels["u_id"],
			//SshPwd:        deploy.Spec.Template.Spec.Containers[0].Args[0],
			//SshPwd: deploy.Spec.Template.Spec.Containers[0].Env[0].Value,
		}
		stslist[i] = tmp
	}
	return &responses.StatefulSetListResponse{Response: responses.OK, Length: num, StatefulSetList: stslist}, nil
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
func UpdateStatefulSet(deploy *appsv1.StatefulSet) (*responses.Response, error) {
	_, err := dao.ClientSet.AppsV1().StatefulSets(deploy.Namespace).Update(context.Background(), deploy, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}
