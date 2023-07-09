package service

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"time"
)

// ListNs 获取所有namespace
func ListNs(label string) (*responses.NsListResponse, error) {
	namespace, err := dao.ClientSet.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	num := len(namespace.Items)
	namespaceList := make([]responses.Ns, num)
	for i, ns := range namespace.Items {
		//if ns.Name == "default" || ns.Name == "kube-node-lease" || ns.Name == "kube-public" || ns.Name == "kube-system" {
		//	continue
		//}
		username := ""
		nickname := ""
		uid, err := strconv.Atoi(ns.Labels["u_id"])
		if err != nil {
			zap.S().Errorln(err)
		} else {
			user, err := dao.GetUserById(uint(uid))
			if err != nil {
				zap.S().Errorln(err)
			} else {
				username = user.Username
				nickname = user.Nickname
			}
		}

		// 增加ttl
		expiredTime := "null"
		ttl, err := dao.GetTtlByNs(ns.Name)
		if err != nil {
			//zap.S().Error(err)
		} else {
			expiredTime = ttl.ExpiredTime.Format("2006-01-02 15:04:05")
		}

		// 资源
		resources := responses.Resources{}
		quota, err := GetResourceQuota(ns.Name)
		if err != nil {
			zap.S().Errorln(err)
		} else {
			limitsCpu := quota.Status.Hard[corev1.ResourceLimitsCPU]
			limitsMemory := quota.Status.Hard[corev1.ResourceLimitsMemory]
			limitsStorage := quota.Status.Hard[corev1.ResourceLimitsEphemeralStorage] // [add] 限制临时存储
			requestPVC := quota.Status.Hard[corev1.ResourceRequestsStorage]           // [add] 限制PVC持久存储
			requestGPU := quota.Status.Hard[LimitsNvidiaGpu]

			usedLimitsCpu := quota.Status.Used[corev1.ResourceLimitsCPU]
			usedLimitsMemory := quota.Status.Used[corev1.ResourceLimitsMemory]
			usedSLimitsStorage := quota.Status.Used[corev1.ResourceLimitsEphemeralStorage] // [add] 已使用临时存储
			usedRequestPVC := quota.Status.Used[corev1.ResourceRequestsStorage]            // [add] 已使用PVC持久存储
			// TODO: GPU
			usedRequestGPU := quota.Status.Used[LimitsNvidiaGpu]

			resources.Cpu = limitsCpu.String()
			resources.Memory = limitsMemory.String()
			resources.Storage = limitsStorage.String()
			resources.PVC = requestPVC.String()
			resources.UsedCpu = usedLimitsCpu.String()
			resources.UsedMemory = usedLimitsMemory.String()
			resources.UsedStorage = usedSLimitsStorage.String()
			resources.UsedPVC = usedRequestPVC.String()
			resources.GPU = requestGPU.String()
			resources.UsedGPU = usedRequestGPU.String()
		}

		tmp := responses.Ns{
			Name:        ns.Name,
			Status:      ns.Status.Phase,
			CreatedAt:   ns.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Username:    username,
			Nickname:    nickname,
			Uid:         uint(uid),
			ExpiredTime: expiredTime,
			Resources:   resources,
		}
		namespaceList[i] = tmp
	}
	return &responses.NsListResponse{Response: responses.OK, Length: num, NsList: namespaceList}, nil
}

// CreateNs 新建属于指定用户的namespace
func CreateNs(name, form string, expiredTime *time.Time, label map[string]string, resources forms.Resources) (*responses.Response, error) {
	annotation := map[string]string{}
	// 利用注释存储表单信息
	if form != "" {
		annotation["form"] = form
	}
	ns := corev1.Namespace{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Namespace"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Labels: label, Annotations: annotation},
	}
	_, err := dao.ClientSet.CoreV1().Namespaces().Create(context.Background(), &ns, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	//创建resourceQuota
	if err = CreateResourceQuota(name, resources); err != nil {
		return nil, err
	}

	if err = CreateLimitRange(name); err != nil {
		return nil, err
	}

	// 创建ttl
	if expiredTime != nil {
		if err = CreateOrUpdateTtl(name, *expiredTime); err != nil {
			return nil, err
		}
	}

	return &responses.OK, nil
}

// DeleteNs 删除指定namespace
func DeleteNs(name string) (*responses.Response, error) {
	err := dao.ClientSet.CoreV1().Namespaces().Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	//TODO 会自动删除PVC吗？
	if err = DeleteTtl(name); err != nil {
		//return nil, err
	}
	return &responses.OK, nil
}

// UpdateNs 更新资源配额、过期时间
func UpdateNs(name, form string, expiredTime *time.Time, resources forms.Resources) (*responses.Response, error) {
	annotation := map[string]string{}
	// 利用注释存储表单信息
	if form != "" {
		annotation["form"] = form
	}
	get, err := dao.ClientSet.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	get.Annotations = annotation
	_, err = dao.ClientSet.CoreV1().Namespaces().Update(context.Background(), get, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	ns := get.Name
	//更新resourceQuota
	err = UpdateResourceQuota(ns, resources)
	if err != nil {
		return nil, err
	}

	// ttl
	if expiredTime != nil {
		if err = CreateOrUpdateTtl(ns, *expiredTime); err != nil {
			return nil, err
		}
	}

	return &responses.OK, nil
}

func GetNs(name string) (*corev1.Namespace, error) {
	get, err := dao.ClientSet.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return get, nil
}
