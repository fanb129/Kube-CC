package service

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetDeploy 获得指定deploy
func GetDeploy(name, ns string) (*appsv1.Deployment, error) {
	get, err := dao.ClientSet.AppsV1().Deployments(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return get, nil
}

// CreateDeploy 创建自定义控制器
func CreateDeploy(name, ns, form string, label map[string]string, spec appsv1.DeploymentSpec) (*appsv1.Deployment, error) {
	annotation := map[string]string{}
	// 利用注释存储表单信息
	if form != "" {
		annotation["form"] = form
	}
	rs := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels:    label,
			// 保存创建时的表单信息
			Annotations: annotation,
		},
		Spec: spec,
	}
	create, err := dao.ClientSet.AppsV1().Deployments(ns).Create(context.Background(), &rs, metav1.CreateOptions{})
	return create, err
}

// ListDeploy 获得指定namespace下的控制器
func ListDeploy(ns, label string) ([]appsv1.Deployment, error) {
	list, err := dao.ClientSet.AppsV1().Deployments(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// DeleteDeploy 删除指定namespace的控制器
func DeleteDeploy(name, ns string) (*responses.Response, error) {
	err := dao.ClientSet.AppsV1().Deployments(ns).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}

// UpdateDeploy 更新deploy
func UpdateDeploy(name, ns, form string, spec appsv1.DeploymentSpec) (*appsv1.Deployment, error) {
	annotation := map[string]string{}
	// 利用注释存储表单信息
	if form != "" {
		annotation["form"] = form
	}
	deploy, err := GetDeploy(name, ns)
	if err != nil {
		return nil, err
	}
	deploy.Spec = spec
	deploy.Annotations = annotation
	update, err := dao.ClientSet.AppsV1().Deployments(deploy.Namespace).Update(context.Background(), deploy, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return update, nil
}
