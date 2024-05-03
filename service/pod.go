package service

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"bytes"
	"context"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

// GetPod 获得指定deploy
func GetPod(name, ns string) (*corev1.Pod, error) {
	get, err := dao.ClientSet.CoreV1().Pods(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return get, nil
}

// ListPod   获得指定namespace下pod
func ListPod(ns string, label string) ([]corev1.Pod, error) {
	list, err := dao.ClientSet.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// DeletePod 删除指定pod
func DeletePod(name, ns string) (*responses.Response, error) {
	err := dao.ClientSet.CoreV1().Pods(ns).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}

// ListDeployPod   获得指定namespace下pod
func ListDeployPod(ns string, label string) ([]responses.DeployPod, error) {
	list, err := dao.ClientSet.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	num := len(list.Items)
	podList := make([]responses.DeployPod, num)
	for i, pod := range list.Items {
		//pod.Status.ContainerStatuses[0].State
		containerName := ""
		if len(pod.Status.ContainerStatuses) == 1 {
			containerName = pod.Status.ContainerStatuses[0].Name
		}
		podList[i] = responses.DeployPod{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Phase:     string(pod.Status.Phase),
			PodIP:     pod.Status.PodIP,
			HostIP:    pod.Status.HostIP,
			Container: containerName,
		}
	}
	return podList, nil
}

// ListStatefulSetPod   获得指定namespace下pod
func ListStatefulSetPod(ns string, label string) ([]responses.StsPod, error) {
	list, err := dao.ClientSet.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	num := len(list.Items)
	podList := make([]responses.StsPod, num)
	for i, pod := range list.Items {
		podList[i] = responses.StsPod{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Phase:     string(pod.Status.Phase),
			PodIP:     pod.Status.PodIP,
			HostIP:    pod.Status.HostIP,
			Container: pod.Status.ContainerStatuses[0].Name,
		}
	}
	return podList, nil
}

func ListJobPod(ns string, label string) ([]responses.JobPod, error) {
	list, err := dao.ClientSet.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	num := len(list.Items)
	podList := make([]responses.JobPod, num)
	for i, pod := range list.Items {
		podList[i] = responses.JobPod{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Phase:     string(pod.Status.Phase),
			Restarts:  pod.Status.ContainerStatuses[0].RestartCount,
			PodIP:     pod.Status.PodIP,
			HostIP:    pod.Status.HostIP,
		}
	}
	return podList, nil
}

func GetPodLog(ns, name string) (*responses.PodLogResponse, error) {
	event, err := GetPodEvent(ns, name)
	if err != nil {
		return nil, err
	}

	req := dao.ClientSet.CoreV1().Pods(ns).GetLogs(name, &corev1.PodLogOptions{})
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return nil, err
	}
	defer podLogs.Close()

	// Copy the logs
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return nil, err
	}
	return &responses.PodLogResponse{
		Response:  responses.OK,
		Namespace: ns,
		Name:      name,
		Log:       event + buf.String() + "1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n",
	}, nil
}

func GetPodEvent(ns, name string) (string, error) {
	events, err := dao.ClientSet.CoreV1().Events(ns).List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s", name),
	})
	if err != nil {
		return "", err
	}
	// 拼接事件信息到 res
	var res strings.Builder
	for _, event := range events.Items {
		res.WriteString(fmt.Sprintf("%s\t%s\t[%s]\t%s\n", event.CreationTimestamp.Format("2006-01-02 15:04:05"), event.Type, event.Reason, event.Message))
	}
	return res.String(), nil
}
