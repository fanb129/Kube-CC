package service

import (
	"Kube-CC/dao"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateResourceQuota 为namespace创建ResourceQuota，进行namespace总的资源限制
func CreateResourceQuota(ns string, spec corev1.ResourceQuotaSpec) error {
	resourceQuota := corev1.ResourceQuota{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ResourceQuota",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ns + "-resourcequota",
			Namespace: ns,
		},
		Spec: spec,
	}
	_, err := dao.ClientSet.CoreV1().ResourceQuotas(ns).Create(&resourceQuota)
	return err
}

// GetResourceQuota 获取指定namespace下的ResourceQuota
func GetResourceQuota(ns string) (*corev1.ResourceQuota, error) {
	resourceQuota, err := dao.ClientSet.CoreV1().ResourceQuotas(ns).Get(ns+"-resourcequota", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return resourceQuota, nil
}

func UpdateResourceQuota(ns string, r *corev1.ResourceQuota) error {
	_, err := dao.ClientSet.CoreV1().ResourceQuotas(ns).Update(r)
	return err
}
