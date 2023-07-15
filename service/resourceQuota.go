package service

import (
	"Kube-CC/common/forms"
	"Kube-CC/dao"
	"context"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//var n = 1 // request == limit
var (
	NvidiaGPU       corev1.ResourceName = "nvidia.com/gpu"
	LimitsNvidiaGpu corev1.ResourceName = "limits.nvidia.com/gpu"
	AmdGpu          corev1.ResourceName = "amd.com/gpu"
	LimitsAmdGpu    corev1.ResourceName = "limits.amd.com/gpu"
	GpuShare        corev1.ResourceName = "aliyun.com/gpu-mem"
	LimitsGpuShare  corev1.ResourceName = "limits.aliyun.com/gpu-mem"
)

//GPU 只能在 limits 部分指定，这意味着：
//
//你可以指定 GPU 的 limits 而不指定其 requests，因为 Kubernetes 将默认使用限制值作为请求值。
//你可以同时指定 limits 和 requests，不过这两个值必须相等。
//你不可以仅指定 requests 而不指定 limits。

// CreateResourceQuota 为namespace创建ResourceQuota，进行namespace总的资源限制
func CreateResourceQuota(ns string, resouces forms.Resources) error {
	if resouces.PvcStorage == "" {
		resouces.PvcStorage = "0"
	}
	if resouces.Gpu == "" {
		resouces.Gpu = "0"
	}
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

				//LimitsNvidiaGpu: resource.MustParse(resouces.Gpu),
				LimitsGpuShare: resource.MustParse(resouces.Gpu),
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
	if resouces.PvcStorage == "" {
		resouces.PvcStorage = "0"
	}
	if resouces.Gpu == "" {
		resouces.Gpu = "0"
	}
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
	//quota.Spec.Hard[LimitsNvidiaGpu] = resource.MustParse(resouces.Gpu)
	quota.Spec.Hard[LimitsGpuShare] = resource.MustParse(resouces.Gpu)
	_, err = dao.ClientSet.CoreV1().ResourceQuotas(ns).Update(context.Background(), quota, metav1.UpdateOptions{})
	return err
}

// SplitRSC 将资源除以n，用作request
func SplitRSC(rsc string, n int) (string, error) {
	quantity, err := resource.ParseQuantity(rsc)
	if err != nil {
		return "", errors.New("failed to parse resource quantity: " + err.Error())
	}
	// 获取该资源的数量单位
	unit := quantity.Format
	// 如果是十进制，即cpu的单位
	if unit == resource.DecimalSI {
		quantity.SetMilli(quantity.MilliValue() / int64(n))
		return quantity.String(), nil
	} else {
		// 将该资源的数量转化为字节数
		bytes := quantity.Value()
		// 将字节数除以n
		bytesPerPart := bytes / int64(n)
		// 将结果转换为适当的单位
		var result string
		switch {
		case bytesPerPart >= 1<<60:
			// 保留到整数位，故将单位转换到下一级单位，
			result = fmt.Sprintf("%.0fPi", float64(bytesPerPart)/(1<<50))
		case bytesPerPart >= 1<<50:
			result = fmt.Sprintf("%.0fTi", float64(bytesPerPart)/(1<<40))
		case bytesPerPart >= 1<<40:
			result = fmt.Sprintf("%.0fGi", float64(bytesPerPart)/(1<<30))
		case bytesPerPart >= 1<<30:
			result = fmt.Sprintf("%.0fMi", float64(bytesPerPart)/(1<<20))
		case bytesPerPart >= 1<<20:
			result = fmt.Sprintf("%.0fKi", float64(bytesPerPart)/(1<<10))
		default:
			quantity.Set(bytesPerPart)
			result = quantity.String()
		}
		return result, nil
	}
}

func VerifyCpu(rsc string) error {
	quantity, err := resource.ParseQuantity(rsc)
	if err != nil {
		return errors.New("failed to parse resource quantity: " + err.Error())
	}
	// 获取该资源的数量单位
	unit := quantity.Format
	// 如果不是十进制，即cpu的单位
	if unit != resource.DecimalSI {
		return errors.New("please use DecimalSI")
	}
	return nil
}
func VerifyResource(rsc string) error {
	quantity, err := resource.ParseQuantity(rsc)
	if err != nil {
		return errors.New("failed to parse resource quantity: " + err.Error())
	}
	// 获取该资源的数量单位
	unit := quantity.Format
	// 如果是十进制，即cpu的单位,
	if unit != resource.BinarySI {
		return errors.New("please use BinarySI")
	}
	return nil
}

func VerifyResourceForm(resources forms.ApplyResources) error {
	err := VerifyCpu(resources.Cpu)
	if err != nil {
		return err
	}
	err = VerifyResource(resources.Memory)
	if err != nil {
		return err
	}
	err = VerifyResource(resources.Gpu)
	if err != nil {
		return err
	}
	err = VerifyResource(resources.Storage)
	if err != nil {
		return err
	}
	err = VerifyResource(resources.PvcStorage)
	if err != nil {
		return err
	}
	return nil
}
