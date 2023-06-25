package service

import (
	"Kube-CC/dao"
	"context"
	"fmt"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

// CreateLimitRange 为namespace创建LimitRange，进行namespace的默认资源限制
func CreateLimitRange(ns string, cpu, memory string, n int) error {
	cpu1, memory1, err := SplitSource(cpu, memory, n)
	if err != nil {
		return err
	}
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
						corev1.ResourceCPU:    resource.MustParse(cpu1),
						corev1.ResourceMemory: resource.MustParse(memory1),
					},
				},
			},
		},
	}
	_, err = dao.ClientSet.CoreV1().LimitRanges(ns).Create(context.Background(), &limitRange, metav1.CreateOptions{})
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

func UpdateLimitRange(ns string, cpu, memory string, n int) error {
	limit, err := dao.ClientSet.CoreV1().LimitRanges(ns).Get(context.Background(), ns+"-limitrange", metav1.GetOptions{})
	if err != nil {
		return err
	}
	// 默认每个container的limit为1/n
	cpu1, memory1, err := SplitSource(cpu, memory, n)
	if err != nil {
		return err
	}
	limit.Spec.Limits[0].Default[corev1.ResourceCPU] = resource.MustParse(cpu1)
	limit.Spec.Limits[0].Default[corev1.ResourceMemory] = resource.MustParse(memory1)

	_, err = dao.ClientSet.CoreV1().LimitRanges(ns).Update(context.Background(), limit, metav1.UpdateOptions{})
	return err
}

func SplitSource(cpu, memory string, n int) (string, string, error) {
	// 分割数字与单位
	index1 := 0
	index2 := 0
	for i, v := range cpu {
		if (v < '0' || v > '9') && v != '.' {
			index1 = i
			break
		}
	}
	for i, v := range memory {
		if (v < '0' || v > '9') && v != '.' {
			index2 = i
			break
		}
	}
	cpu1, err := strconv.ParseFloat(cpu[:index1], 64)
	if err != nil {
		return "", "", err
	}
	memory1, err := strconv.ParseFloat(memory[:index2], 64)
	if err != nil {
		return "", "", err
	}

	cpu11 := fmt.Sprintf("%.3f", cpu1/float64(n)) + cpu[index1:]
	memory11 := fmt.Sprintf("%.3f", memory1/float64(n)) + memory[index2:]
	zap.S().Infoln(cpu11, cpu11, memory1, memory11)
	return cpu11, memory11, nil

}
