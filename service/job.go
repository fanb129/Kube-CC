package service

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetJob 获得指定namespace下job
func GetJob(ns string, label string) (*responses.JobListResponse, error) {
	pods, err := dao.ClientSet.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}

	num := len(pods.Items)
	jobList := make([]responses.Job, num)

	for i, job := range pods.Items {
		tmp := responses.Job{
			Name:              job.Name,
			Namespace:         job.Namespace,
			CreatedAt:         job.CreationTimestamp.Format("2006-01-02 15:04:05"),
			NodeIp:            job.Status.HostIP,
			Phase:             job.Status.Phase,
			ContainerStatuses: job.Status.ContainerStatuses,
		}
		jobList[i] = tmp
	}
	return &responses.JobListResponse{Response: responses.OK, Length: num, JobList: jobList}, nil
}

// DeleteJob 删除指定job
func DeleteJob(name, ns string) (*responses.Response, error) {
	err := dao.ClientSet.CoreV1().Pods(ns).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}

func CreateJob(name, ns string, label map[string]string, spec corev1.PodSpec) (*corev1.Pod, error) {
	job := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "v2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels:    label,
		},
		Spec: spec,
	}
	create, err := dao.ClientSet.CoreV1().Pods(ns).Create(context.Background(), &job, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return create, nil
}

func UpdateJob(job *corev1.Pod) (*responses.Response, error) {
	_, err := dao.ClientSet.CoreV1().Pods(job.Name).Update(context.Background(), job, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}
