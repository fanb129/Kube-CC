package service

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/dao"
)

// CreateDeploy 创建自定义控制器
func CreateDeploy(name, ns string, label map[string]string, spec appsv1.DeploymentSpec) (*appsv1.Deployment, error) {
	rs := appsv1.Deployment{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns, Labels: label},
		Spec:       spec,
	}
	create, err := dao.ClientSet.AppsV1().Deployments(ns).Create(&rs)
	return create, err
}

// GetDeploy 获得指定namespace下的控制器
func GetDeploy(ns, label string) (*common.DeployListResponse, error) {
	list, err := dao.ClientSet.AppsV1().Deployments(ns).List(metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	num := len(list.Items)
	deployList := make([]common.Deploy, num)
	for i, deploy := range list.Items {
		tmp := common.Deploy{
			Name: deploy.Name,
		}
		deployList[i] = tmp
	}
	return &common.DeployListResponse{Response: common.OK, Length: num, DeployList: deployList}, nil
}

// DeleteDeploy 删除指定namespace先的控制器
func DeleteDeploy(name, ns string) (*common.Response, error) {
	err := dao.ClientSet.AppsV1().Deployments(ns).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &common.OK, nil
}
