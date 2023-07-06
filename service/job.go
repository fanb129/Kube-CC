package service

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListJob 获得指定namespace下job
func ListJob(ns string, label string) ([]batchv1.Job, error) {
	list, err := dao.ClientSet.BatchV1().Jobs(ns).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// DeleteJob 删除指定job
func DeleteJob(name, ns string) (*responses.Response, error) {
	err := dao.ClientSet.BatchV1().Jobs(ns).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}

func CreateJob(name, ns string, label map[string]string, spec batchv1.JobSpec) (*batchv1.Job, error) {
	job := batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels:    label,
		},
		Spec: spec,
	}
	create, err := dao.ClientSet.BatchV1().Jobs(ns).Create(context.Background(), &job, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return create, nil
}

func GetJob(name, ns string) (*batchv1.Job, error) {
	get, err := dao.ClientSet.BatchV1().Jobs(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return get, err
}

func UpdateJob(name, ns string, spec batchv1.JobSpec) (*batchv1.Job, error) {
	job, err := GetJob(name, ns)
	if err != nil {
		return nil, err
	}
	job.Spec = spec
	update, err := dao.ClientSet.BatchV1().Jobs(ns).Update(context.Background(), job, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return update, nil
}
