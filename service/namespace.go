package service

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	"errors"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"strconv"
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
		//expiredTime := "null"
		//ttl, err := dao.GetTtlByNs(ns.Name)
		//if err != nil {
		//	//zap.S().Error(err)
		//} else {
		//	expiredTime = ttl.ExpiredTime.Format("2006-01-02 15:04:05")
		//}

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

			resources.UsedCpuValue = usedLimitsCpu.MilliValue()
			resources.UsedMemoryValue = usedLimitsMemory.Value()
			resources.UsedStorageValue = usedSLimitsStorage.Value()
			resources.UsedPVCValue = usedRequestPVC.Value()
			resources.UsedGPUValue = usedRequestGPU.Value()
		}

		tmp := responses.Ns{
			Name:      ns.Name,
			Status:    ns.Status.Phase,
			CreatedAt: ns.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Username:  username,
			Nickname:  nickname,
			Uid:       uint(uid),
			//ExpiredTime: expiredTime,
			Resources: resources,
		}
		namespaceList[i] = tmp
	}
	return &responses.NsListResponse{Response: responses.OK, Length: num, NsList: namespaceList}, nil
}

// CreateNs 新建属于指定用户的namespace
func CreateNs(name, form string, label map[string]string, resources forms.Resources) (*responses.Response, error) {
	uid := label["u_id"]
	err := VerifyNsResource(uid, "", resources)
	if err != nil {
		return nil, err
	}
	// 利用注释存储表单信息
	annotation := map[string]string{}
	if form != "" {
		annotation["form"] = form
	}
	ns := corev1.Namespace{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Namespace"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Labels: label, Annotations: annotation},
	}
	_, err = dao.ClientSet.CoreV1().Namespaces().Create(context.Background(), &ns, metav1.CreateOptions{})
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

	//// 创建ttl
	//if expiredTime != nil {
	//	if err = CreateOrUpdateTtl(name, *expiredTime); err != nil {
	//		return nil, err
	//	}
	//}

	return &responses.OK, nil
}

// DeleteNs 删除指定namespace
func DeleteNs(name string) (*responses.Response, error) {
	err := dao.ClientSet.CoreV1().Namespaces().Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	//TODO 会自动删除PVC吗？
	//if err = DeleteTtl(name); err != nil {
	//	//return nil, err
	//}
	return &responses.OK, nil
}

// UpdateNs 更新资源配额、过期时间
func UpdateNs(name, form string, resources forms.Resources) (*responses.Response, error) {
	annotation := map[string]string{}
	// 利用注释存储表单信息
	if form != "" {
		annotation["form"] = form
	}
	get, err := GetNs(name)
	if err != nil {
		return nil, err
	}
	uid := get.Labels["u_id"]
	err = VerifyNsResource(uid, name, resources)
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
	//if expiredTime != nil {
	//	if err = CreateOrUpdateTtl(ns, *expiredTime); err != nil {
	//		return nil, err
	//	}
	//}

	return &responses.OK, nil
}

func GetNs(name string) (*corev1.Namespace, error) {
	get, err := dao.ClientSet.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return get, nil
}

// VerifyNsResource
//如果是更新前的资源验证的话，计算剩余量时排除掉更新所选的ns
//如果是创建ns前的资源验证，ns参数为空，则应该计算所有的资源剩余量
func VerifyNsResource(uid, name string, resources forms.Resources) error {
	// 请求的资源量
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
	requestCpu := resource.MustParse(resources.Cpu)
	requestMemory := resource.MustParse(resources.Memory)
	requestGpu := resource.MustParse(resources.Gpu)
	requestStorage := resource.MustParse(resources.Storage)
	requestPvc := resource.MustParse(resources.PvcStorage)
	// 该用户的资源限额
	intuid, err := strconv.Atoi(uid)
	if err != nil {
		return err
	}
	user, err := dao.GetUserById(uint(intuid))
	if err != nil {
		return err
	}
	cpu := resource.MustParse(user.Cpu)
	memory := resource.MustParse(user.Memory)
	gpu := resource.MustParse(user.Gpu)
	storage := resource.MustParse(user.Storage)
	pvc := resource.MustParse(user.Pvcstorage)

	// 创建ns前先计算资源是否超额
	label := map[string]string{
		"u_id": uid,
	}
	// 将map标签转换为string
	selector := labels.SelectorFromSet(label).String()
	nsList, err := ListNs(selector)
	if err != nil {
		return err
	}
	for _, ns := range nsList.NsList {
		if ns.Name != name {
			cpu.Sub(resource.MustParse(ns.Cpu))
			memory.Sub(resource.MustParse(ns.Memory))
			gpu.Sub(resource.MustParse(ns.GPU))
			storage.Sub(resource.MustParse(ns.Storage))
			pvc.Sub(resource.MustParse(ns.PVC))
		}
	}
	if cpu.MilliValue() < requestCpu.MilliValue() {
		return errors.New("left cpu:" + cpu.String() + " less than request cpu:" + requestCpu.String())
	}
	if gpu.Value() < requestGpu.Value() {
		return errors.New("left gpu:" + gpu.String() + " is less than request gpu:" + requestGpu.String())
	}
	if memory.Value() < requestMemory.Value() {
		return errors.New("left memory:" + memory.String() + " is less than request memory:" + requestMemory.String())
	}
	if storage.Value() < requestStorage.Value() {
		return errors.New("left storage:" + storage.String() + " is less than request storage:" + requestStorage.String())
	}
	if pvc.Value() < requestPvc.Value() {
		return errors.New("left pvc:" + pvc.String() + " is less than request pvc:" + requestPvc.String())
	}
	return nil
}

// GetUserNsTotal 返回当前用户总的ns的使用情况
func GetUserNsTotal(uid string) (*responses.UserTotalNs, error) {
	label := map[string]string{
		"u_id": uid,
	}
	// 将map标签转换为string
	selector := labels.SelectorFromSet(label).String()
	nsList, err := ListNs(selector)
	if err != nil {
		return nil, err
	}
	cpu := resource.MustParse("0")
	usedCpu := resource.MustParse("0")
	memory := resource.MustParse("0")
	usedMemory := resource.MustParse("0")
	gpu := resource.MustParse("0")
	usedGpu := resource.MustParse("0")
	storage := resource.MustParse("0")
	usedStorage := resource.MustParse("0")
	pvc := resource.MustParse("0")
	usedPvc := resource.MustParse("0")
	for _, ns := range nsList.NsList {
		cpu.Add(resource.MustParse(ns.Cpu))
		usedCpu.Add(resource.MustParse(ns.UsedCpu))

		memory.Add(resource.MustParse(ns.Memory))
		usedMemory.Add(resource.MustParse(ns.UsedMemory))

		gpu.Add(resource.MustParse(ns.GPU))
		usedGpu.Add(resource.MustParse(ns.UsedGPU))

		storage.Add(resource.MustParse(ns.Storage))
		usedStorage.Add(resource.MustParse(ns.UsedStorage))

		pvc.Add(resource.MustParse(ns.PVC))
		usedPvc.Add(resource.MustParse(ns.UsedPVC))
	}
	rsp := responses.UserTotalNs{
		Response: responses.OK,
		Cpu:      cpu.String(),
		UsedCpu:  usedCpu.String(),
		CpuRatio: float64(usedCpu.MilliValue()) / float64(cpu.MilliValue()),

		Memory:      memory.String(),
		UsedMemory:  usedMemory.String(),
		MemoryRatio: float64(usedMemory.Value()) / float64(memory.Value()),

		Storage:      storage.String(),
		UsedStorage:  usedStorage.String(),
		StorageRatio: float64(usedStorage.Value()) / float64(storage.Value()),

		PVC:      pvc.String(),
		UsedPVC:  usedPvc.String(),
		PvcRatio: float64(usedPvc.Value()) / float64(pvc.Value()),

		GPU:      gpu.String(),
		UsedGPU:  usedGpu.String(),
		GpuRatio: float64(usedGpu.Value()) / float64(gpu.Value()),
	}

	return &rsp, nil
}
