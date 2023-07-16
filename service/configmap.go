package service

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	configMap, err := dao.ClientSet.CoreV1().ConfigMaps(ns).Create(context.Background(), &cm, metav1.CreateOptions{})
	return configMap, err
}

// DeleteConfigMap 删除指定namespace的configMap
func DeleteConfigMap(name, ns string) (*responses.Response, error) {
	err := dao.ClientSet.CoreV1().ConfigMaps(ns).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}

// GetConfigMap 获取configMap
func GetConfigMap(name, ns string) (*corev1.ConfigMap, error) {
	get, err := dao.ClientSet.CoreV1().ConfigMaps(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return get, err
}

// UpdateConfigMap 更新configMap
func UpdateConfigMap(name, ns string, data map[string]string) (*corev1.ConfigMap, error) {
	configMap, err := GetConfigMap(name, ns)
	if err != nil {
		return nil, err
	}
	configMap.Data = data
	update, err := dao.ClientSet.CoreV1().ConfigMaps(ns).Update(context.Background(), configMap, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return update, nil
}

func CreateOrUpdateConfigMap(name, ns string, data map[string]string) (*corev1.ConfigMap, error) {
	configMap, _ := GetConfigMap(name, ns)
	if configMap == nil {
		create, err := CreateConfigMap(name, ns, map[string]string{}, data)
		if err != nil {
			return nil, err
		}
		return create, nil
	} else {
		update, err := UpdateConfigMap(name, ns, data)
		if err != nil {
			return nil, err
		}
		return update, nil
	}
}
