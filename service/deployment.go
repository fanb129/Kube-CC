package service

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetADeploy 获得指定deploy
func GetADeploy(name, ns string) (*appsv1.Deployment, error) {
	get, err := dao.ClientSet.AppsV1().Deployments(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return get, nil
}

// CreateDeploy 创建自定义控制器
func CreateDeploy(name, ns string, spec appsv1.DeploymentSpec) (*appsv1.Deployment, error) { //此处去掉形参”label map[string]string“
	rs := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			//Labels:    label,
		},
		Spec: spec,
	}
	create, err := dao.ClientSet.AppsV1().Deployments(ns).Create(context.Background(), &rs, metav1.CreateOptions{})
	return create, err
}

// GetDeploy 获得指定namespace下的控制器
func GetDeploy(ns, label string) (*responses.DeployListResponse, error) {
	list, err := dao.ClientSet.AppsV1().Deployments(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	num := len(list.Items)
	deployList := make([]responses.Deploy, num)
	for i, deploy := range list.Items {
		tmp := responses.Deploy{
			Name:              deploy.Name,
			Namespace:         deploy.Namespace,
			CreatedAt:         deploy.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Replicas:          deploy.Status.Replicas,
			UpdatedReplicas:   deploy.Status.UpdatedReplicas,
			ReadyReplicas:     deploy.Status.ReadyReplicas,
			AvailableReplicas: deploy.Status.AvailableReplicas,
			Uid:               deploy.Labels["u_id"],
			//SshPwd:        deploy.Spec.Template.Spec.Containers[0].Args[0],
			//SshPwd: deploy.Spec.Template.Spec.Containers[0].Env[0].Value,
		}
		deployList[i] = tmp
	}
	return &responses.DeployListResponse{Response: responses.OK, Length: num, DeployList: deployList}, nil
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
func UpdateDeploy(deploy *appsv1.Deployment) (*responses.Response, error) {
	_, err := dao.ClientSet.AppsV1().Deployments(deploy.Namespace).Update(context.Background(), deploy, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}
