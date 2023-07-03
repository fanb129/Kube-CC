package yamlApply

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"Kube-CC/service"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceApply namespace的更新或新建
func NamespaceApply(ns *corev1.Namespace) (*responses.Response, error) {
	name := ns.Name
	if _, err := dao.ClientSet.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			if _, err := dao.ClientSet.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{}); err != nil {
				return nil, err
			}
			return &responses.OK, nil
		} else {
			return nil, err
		}
	} else {
		if _, err := dao.ClientSet.CoreV1().Namespaces().Update(context.Background(), ns, metav1.UpdateOptions{}); err != nil {
			return nil, err
		}
		return &responses.OK, nil
	}

}

// DeployApply  deploy的更新或者创建
func DeployApply(deploy *appsv1.Deployment) (*responses.Response, error) {
	name := deploy.Name
	ns := deploy.Namespace
	labels := deploy.Labels
	if labels == nil {
		labels = make(map[string]string)
	}
	if _, err := dao.ClientSet.AppsV1().Deployments(ns).Get(context.Background(), name, metav1.GetOptions{}); err != nil {
		// 不存在则创建
		if errors.IsNotFound(err) {
			// 获取namespace，(不再)提取出uid的label
			//_, err := dao.ClientSet.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			if _, err := service.CreateDeploy(name, ns, labels, deploy.Spec); err != nil {
				return nil, err
			}
			return &responses.OK, nil
		} else { // 其他错误直接返回
			return nil, err
		}
	} else { // 存在则更新
		response, err := service.UpdateDeploy(deploy)
		return response, err
	}

}

// StatefulSetApply  statefulSet的更新或者创建
func StatefulSetApply(statefulSet *appsv1.StatefulSet) (*responses.Response, error) {
	name := statefulSet.Name
	ns := statefulSet.Namespace
	labels := statefulSet.Labels
	if labels == nil {
		labels = make(map[string]string)
	}
	if _, err := dao.ClientSet.AppsV1().StatefulSets(ns).Get(context.Background(), name, metav1.GetOptions{}); err != nil {
		// 不存在则创建
		if errors.IsNotFound(err) {
			// 获取namespace，提取出uid的label
			//_, err := dao.ClientSet.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			if _, err := service.CreateStatefulSet(name, ns, labels, statefulSet.Spec); err != nil {
				return nil, err
			}
			return &responses.OK, nil
		} else { // 其他错误直接返回
			return nil, err
		}
	} else { // 存在则更新
		response, err := service.UpdateStatefulSet(statefulSet)
		return response, err
	}

}

// ServiceApply service的更新或者创建
func ServiceApply(svc *corev1.Service) (*responses.Response, error) {
	name := svc.Name
	ns := svc.Namespace
	labels := svc.Labels
	if labels == nil {
		labels = make(map[string]string)
	}
	if _, err := dao.ClientSet.CoreV1().Services(ns).Get(context.Background(), name, metav1.GetOptions{}); err != nil {
		// 不存在则创建
		if errors.IsNotFound(err) {
			// 获取namespace，提取出uid的label
			//get, err := dao.ClientSet.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
			//if err != nil {
			//	return nil, err
			//}
			//uid, ok := get.Labels["u_id"]
			//if ok {
			//	labels["u_id"] = uid //service的label
			//}
			if _, err := service.CreateService(name, ns, labels, svc.Spec); err != nil {
				return nil, err
			}
			return &responses.OK, nil
		} else { // 其他错误直接返回
			return nil, err
		}
	} else {
		res, err := service.UpdateService(svc)
		return res, err
	}
}

func PodApply(pod *corev1.Pod) (*responses.Response, error) {
	name := pod.Name
	ns := pod.Namespace
	labels := pod.Labels
	if labels == nil {
		labels = make(map[string]string)
	}
	if _, err := dao.ClientSet.CoreV1().Pods(ns).Get(context.Background(), name, metav1.GetOptions{}); err != nil {
		// 不存在则创建
		if errors.IsNotFound(err) {
			// 获取namespace，提取出uid的label
			//get, err := dao.ClientSet.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
			//if err != nil {
			//	return nil, err
			//}
			//uid, ok := get.Labels["u_id"]
			//if ok {
			//	labels["u_id"] = uid //service的label
			//}
			if _, err := service.CreatePod(name, ns, labels, pod.Spec); err != nil {
				return nil, err
			}
			return &responses.OK, nil
		} else { // 其他错误直接返回
			return nil, err
		}
	} else {
		res, err := service.UpdatePod(pod)
		return res, err
	}
}

func JobApply(job *corev1.Pod) (*responses.Response, error) {
	name := job.Name
	ns := job.Namespace
	labels := job.Labels
	if labels == nil {
		labels = make(map[string]string)
	}
	if _, err := dao.ClientSet.CoreV1().Pods(ns).Get(context.Background(), name, metav1.GetOptions{}); err != nil {
		// 不存在则创建
		if errors.IsNotFound(err) {
			// 获取namespace，提取出uid的label
			//_, err := dao.ClientSet.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
			//if err != nil {
			//	return nil, err
			//}
			if _, err := service.CreateJob(name, ns, labels, job.Spec); err != nil {
				return nil, err
			}
			return &responses.OK, nil
		} else { // 其他错误直接返回
			return nil, err
		}
	} else {
		res, err := service.UpdateJob(job)
		return res, err
	}
}
