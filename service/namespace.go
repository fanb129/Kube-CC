package service

import (
	"Kube-CC/common"
	"Kube-CC/dao"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

// GetNs 获取所有namespace
func GetNs(label string) (*common.NsListResponse, error) {
	namespace, err := dao.ClientSet.CoreV1().Namespaces().List(metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	num := len(namespace.Items)
	namespaceList := make([]common.Ns, num)
	for i, ns := range namespace.Items {
		//if ns.Name == "default" || ns.Name == "kube-node-lease" || ns.Name == "kube-public" || ns.Name == "kube-system" {
		//	continue
		//}
		uid, err := strconv.Atoi(ns.Labels["u_id"])
		username := ""
		nickname := ""
		if err == nil {
			user, err := dao.GetUserById(uint(uid))
			if err == nil {
				username = user.Username
				nickname = user.Nickname
			}
		}

		tmp := common.Ns{
			Name:      ns.Name,
			Status:    ns.Status.Phase,
			CreatedAt: ns.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Username:  username,
			Nickname:  nickname,
			Uid:       uint(uid),
		}
		namespaceList[i] = tmp
	}
	return &common.NsListResponse{Response: common.OK, Length: num, NsList: namespaceList}, nil
}

// CreateNs 新建属于指定用户的namespace，u_id == 0 则不添加标签
func CreateNs(name string, label map[string]string) (*common.Response, error) {
	ns := corev1.Namespace{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Namespace"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Labels: label},
	}
	_, err := dao.ClientSet.CoreV1().Namespaces().Create(&ns)
	if err != nil {
		return nil, err
	}
	// 创建resourceQuota
	//spec := corev1.ResourceQuotaSpec{
	//	Hard: corev1.ResourceList{
	//		// CPU, in cores. (500m = .5 cores)
	//		corev1.ResourceRequestsCPU: resource.MustParse("100m"),
	//		corev1.ResourceLimitsCPU:   resource.MustParse("100m"),
	//
	//		// Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	//		corev1.ResourceRequestsMemory: resource.MustParse("100m"),
	//		corev1.ResourceLimitsMemory:   resource.MustParse("100m"),
	//
	//		// Volume size, in bytes (e,g. 5Gi = 5GiB = 5 * 1024 * 1024 * 1024)
	//		//corev1.ResourceRequestsStorage: resource.MustParse("100m"),
	//		// Local ephemeral storage, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	//		// The resource name for ResourceEphemeralStorage is alpha and it can change across releases.
	//		//corev1.ResourceEphemeralStorage: resource.Quantity{},
	//	},
	//}
	//if err = createResourceQuota(name, spec); err != nil {
	//	return nil, err
	//}
	return &common.OK, nil
}

// DeleteNs 删除指定namespace
func DeleteNs(name string) (*common.Response, error) {
	err := dao.ClientSet.CoreV1().Namespaces().Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &common.OK, nil
}

// UpdateNs 分配namespace
func UpdateNs(name, uid string) (*common.Response, error) {
	get, err := dao.ClientSet.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// 未改变
	if get.Labels["u_id"] == uid {
		return &common.OK, nil
	}
	// 更新namespace的uid
	if uid == "" {
		delete(get.Labels, "u_id")
	} else {
		get.Labels["u_id"] = uid
	}
	if _, err := dao.ClientSet.CoreV1().Namespaces().Update(get); err != nil {
		return nil, err
	}

	ns := get.Name

	// 更新namespace下所有deploy的uid
	deployList, err := GetDeploy(ns, "")
	if err == nil {
		for i := 0; i < deployList.Length; i++ {
			name := deployList.DeployList[i].Name
			deployment, err := dao.ClientSet.AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			if uid == "" {
				delete(deployment.Labels, "u_id")
				delete(deployment.Spec.Template.Labels, "u_id")
			} else {
				deployment.Labels["u_id"] = uid
				deployment.Spec.Template.Labels["u_id"] = uid
			}
			if _, err := UpdateDeploy(deployment); err != nil {
				return nil, err
			}
		}
	}

	// 更新namespace下所有service的uid
	serviceList, err := GetService(ns, "")
	if err == nil {
		for i := 0; i < serviceList.Length; i++ {
			name := serviceList.ServiceList[i].Name
			service, err := dao.ClientSet.CoreV1().Services(ns).Get(name, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			if uid == "" {
				delete(service.Labels, "u_id")
			} else {
				service.Labels["u_id"] = uid
			}
			if _, err := UpdateService(service); err != nil {
				return nil, err
			}
		}
	}

	// 更新namespace下所有pod的uid
	podList, err := GetPod(ns, "")
	if err == nil {
		for i := 0; i < podList.Length; i++ {
			name := podList.PodList[i].Name
			pod, err := dao.ClientSet.CoreV1().Pods(ns).Get(name, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			if uid == "" {
				delete(pod.Labels, "u_id")
			} else {
				pod.Labels["u_id"] = uid
			}
			if _, err := UpdatePod(pod); err != nil {
				return nil, err
			}
		}
	}

	return &common.OK, nil
}

// 为namespace创建ResourceQuota，进行namespace总的资源限制
func createResourceQuota(ns string, spec corev1.ResourceQuotaSpec) error {
	resourceQuota := corev1.ResourceQuota{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ResourceQuota",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ns + "-ResourceQuota",
			Namespace: ns,
		},
		Spec: spec,
	}
	_, err := dao.ClientSet.CoreV1().ResourceQuotas(ns).Create(&resourceQuota)
	return err
}

// 获取指定namespace下的ResourceQuota
func getResourceQuota(ns string) (*corev1.ResourceQuota, error) {
	resourceQuota, err := dao.ClientSet.CoreV1().ResourceQuotas(ns).Get(ns+"-ResourceQuota", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return resourceQuota, nil
}

func updateResourceQuota(ns string, r *corev1.ResourceQuota) error {
	_, err := dao.ClientSet.CoreV1().ResourceQuotas(ns).Update(r)
	return err
}
