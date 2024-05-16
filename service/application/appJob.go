package application

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service"
	"encoding/json"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/uuid"
)

// CreateAppJob 创建一个appJob一次性任务
func CreateAppJob(form forms.JobAddForm) (*responses.Response, error) {
	// 将form序列化为string，存入注释
	jsonBytes, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}
	strForm := string(jsonBytes)

	// 创建uuid，以便筛选出属于同一组的deploy、pod、service等
	newUuid := string(uuid.NewUUID())
	label := map[string]string{
		"uuid": newUuid,
	}

	manualSelector := true
	spec := batchv1.JobSpec{
		// 完成数
		Completions: &form.Completions,
		// 并发数
		Parallelism:    &form.Parallelism,
		ManualSelector: &manualSelector,
		Selector:       &metav1.LabelSelector{MatchLabels: label},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: label},
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
				Containers: []corev1.Container{
					{
						Name:            form.Name,
						Image:           form.Image,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Command:         form.Command,
						Args:            form.Args,
					},
				},
			},
		},
	}
	_, err = service.CreateJob(form.Name, form.Namespace, strForm, label, spec)
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}

// DeleteAppJob 删除指定job
func DeleteAppJob(name, ns string) (*responses.Response, error) {
	rsp, err := service.DeleteJob(name, ns)
	return rsp, err
}

// ListAppJob 列出joblist的详细信息
func ListAppJob(ns string, label string) (*responses.AppJobList, error) {
	jobs, err := service.ListJob(ns, label)
	if err != nil {
		return nil, err
	}
	num := len(jobs)
	jobList := make([]responses.AppJob, num)
	for i, job := range jobs {
		// 获取对应pod
		label1 := map[string]string{
			"uuid": job.Labels["uuid"],
		}
		selector := labels.SelectorFromSet(label1).String()
		podList, err := service.ListJobPod(ns, selector)
		if err != nil {
			return nil, err
		}
		startTime := job.Status.StartTime
		completionTime := job.Status.CompletionTime
		jobList[i] = responses.AppJob{
			Name:        job.Name,
			Namespace:   job.Namespace,
			CreatedAt:   job.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Succeeded:   job.Status.Succeeded,
			Completions: *job.Spec.Completions,
			Duration:    completionTime.Sub(startTime.Time).String(),
			Image:       job.Spec.Template.Spec.Containers[0].Image,
			PodList:     podList,
		}
	}
	return &responses.AppJobList{Response: responses.OK, Length: num, AppJobList: jobList}, nil
}

// GetAppJob 返回job填写时的表单信息，方便再次运行
func GetAppJob(name, ns string) (*responses.InfoJob, error) {
	form := forms.JobAddForm{}
	job, err := service.GetJob(name, ns)
	if err != nil {
		return nil, err
	}
	strForm := job.Annotations["form"]
	err = json.Unmarshal([]byte(strForm), &form)
	if err != nil {
		return nil, err
	}
	return &responses.InfoJob{Response: responses.OK, Form: form}, nil
}
