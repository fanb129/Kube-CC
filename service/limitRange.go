package service

import (
	"Kube-CC/dao"
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateLimitRange 为namespace创建LimitRange，进行namespace的默认资源限制
func CreateLimitRange(ns string) error {
	limitRange := corev1.LimitRange{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "LimitRange",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ns + "-limitrange",
			Namespace: ns,
		},
		Spec: corev1.LimitRangeSpec{
			Limits: []corev1.LimitRangeItem{
				{
					Type: corev1.LimitTypeContainer,
					Default: corev1.ResourceList{
						corev1.ResourceCPU:              resource.MustParse("200m"),
						corev1.ResourceMemory:           resource.MustParse("256Mi"),
						corev1.ResourceEphemeralStorage: resource.MustParse("2Gi"),
					},
					DefaultRequest: corev1.ResourceList{
						corev1.ResourceCPU:              resource.MustParse("100m"),
						corev1.ResourceMemory:           resource.MustParse("64Mi"),
						corev1.ResourceEphemeralStorage: resource.MustParse("512Mi"),
					},
				},
			},
		},
	}
	_, err := dao.ClientSet.CoreV1().LimitRanges(ns).Create(context.Background(), &limitRange, metav1.CreateOptions{})
	return err
}

// DeleteLimitRange 删除limitrange
func DeleteLimitRange(ns string) error {
	err := dao.ClientSet.CoreV1().LimitRanges(ns).Delete(context.Background(), ns+"-limitrange", metav1.DeleteOptions{})
	return err
}

//// GetLimitRange 获取指定namespace下的LimitRange
//func GetLimitRange(ns string) (*corev1.LimitRange, error) {
//	limit, err := dao.ClientSet.CoreV1().LimitRanges(ns).Get(ns+"-limitrange", metav1.GetOptions{})
//	if err != nil {
//		return nil, err
//	}
//	return limit, nil
//}

//func UpdateLimitRange(ns string, cpu, memory string, n int) error {
//	limit, err := dao.ClientSet.CoreV1().LimitRanges(ns).Get(context.Background(), ns+"-limitrange", metav1.GetOptions{})
//	if err != nil {
//		return err
//	}
//	// 默认每个container的limit为1/n
//	cpu1, memory1, err := SplitSource(cpu, memory, n)
//	if err != nil {
//		return err
//	}
//	limit.Spec.Limits[0].Default[corev1.ResourceCPU] = resource.MustParse(cpu1)
//	limit.Spec.Limits[0].Default[corev1.ResourceMemory] = resource.MustParse(memory1)
//
//	_, err = dao.ClientSet.CoreV1().LimitRanges(ns).Update(context.Background(), limit, metav1.UpdateOptions{})
//	return err
//}
