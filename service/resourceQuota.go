package service

import (
	"Kube-CC/common/forms"
	"Kube-CC/dao"
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

//var n = 1 // request == limit
var ResourceGPU corev1.ResourceName = "requests.nvidia.com/gpu"

// CreateResourceQuota 为namespace创建ResourceQuota，进行namespace总的资源限制
func CreateResourceQuota(ns string, resouces forms.Resources) error {
	resourceQuota := corev1.ResourceQuota{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ResourceQuota",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ns + "-resourcequota",
			Namespace: ns,
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: corev1.ResourceList{
				// 存储 pvc限制
				corev1.ResourceRequestsStorage: resource.MustParse(resouces.PvcStorage),
				// 临时存储
				corev1.ResourceRequestsEphemeralStorage: resource.MustParse(resouces.Storage),
				corev1.ResourceLimitsEphemeralStorage:   resource.MustParse(resouces.Storage),
				// cpu限制
				corev1.ResourceRequestsCPU: resource.MustParse(resouces.Cpu),
				corev1.ResourceLimitsCPU:   resource.MustParse(resouces.Cpu),
				// 内存限制
				corev1.ResourceRequestsMemory: resource.MustParse(resouces.Memory),
				corev1.ResourceLimitsMemory:   resource.MustParse(resouces.Memory),

				// TODO:GPU
				//ResourceGPU: resource.MustParse(resouces.Gpu),
			},
		},
	}
	_, err := dao.ClientSet.CoreV1().ResourceQuotas(ns).Create(context.Background(), &resourceQuota, metav1.CreateOptions{})
	return err
}

// GetResourceQuota 获取指定namespace下的ResourceQuota
func GetResourceQuota(ns string) (*corev1.ResourceQuota, error) {
	resourceQuota, err := dao.ClientSet.CoreV1().ResourceQuotas(ns).Get(context.Background(), ns+"-resourcequota", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return resourceQuota, nil
}

// UpdateResourceQuota 更新
func UpdateResourceQuota(ns string, resouces forms.Resources) error {
	quota, err := GetResourceQuota(ns)
	if err != nil {
		return err
	}
	quota.Spec.Hard[corev1.ResourceRequestsStorage] = resource.MustParse(resouces.PvcStorage)

	quota.Spec.Hard[corev1.ResourceRequestsEphemeralStorage] = resource.MustParse(resouces.Storage)
	quota.Spec.Hard[corev1.ResourceLimitsEphemeralStorage] = resource.MustParse(resouces.Storage)

	quota.Spec.Hard[corev1.ResourceRequestsCPU] = resource.MustParse(resouces.Cpu)
	quota.Spec.Hard[corev1.ResourceLimitsCPU] = resource.MustParse(resouces.Cpu)

	quota.Spec.Hard[corev1.ResourceRequestsMemory] = resource.MustParse(resouces.Memory)
	quota.Spec.Hard[corev1.ResourceLimitsMemory] = resource.MustParse(resouces.Memory)
	// TODO:GPU
	//quota.Spec.Hard[ResourceGPU] = resource.MustParse(resouces.Gpu)
	_, err = dao.ClientSet.CoreV1().ResourceQuotas(ns).Update(context.Background(), quota, metav1.UpdateOptions{})
	return err
}

// SplitRSC 将资源除以n，用作request
func SplitRSC(rsc string, n int) (string, error) {
	if n <= 1 {
		return rsc, nil
	}
	// 分割出数字与单位
	index := 0
	for _, v := range rsc {
		if (v < '0' || v > '9') && v != '.' {
			break
		}
		index++
	}
	//转换为float
	float, err := strconv.ParseFloat(rsc[:index], 64)
	if err != nil {
		return "", err
	}
	// 如果没有单位，x1000，加上单位m
	if index == len(rsc) {
		m := int(float * 1000 / float64(n))
		str := strconv.Itoa(m)
		return str + "m", nil
	} else {
		m := int(float / float64(n))
		str := strconv.Itoa(m)
		return str + rsc[index:], nil
	}

}
