package service

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/dao"
)

// CreateConfigMap 创建存储配置
func CreateConfigMap(name, ns string, label, data map[string]string) (*corev1.ConfigMap, error) {
	cm := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels:    label,
		},
		Data: data,
	}
	configMap, err := dao.ClientSet.CoreV1().ConfigMaps(ns).Create(&cm)
	return configMap, err
}

// DeleteConfigMap 删除指定namespace的configMap
func DeleteConfigMap(name, ns string) (*common.Response, error) {
	err := dao.ClientSet.CoreV1().ConfigMaps(ns).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &common.OK, nil
}
