package service

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/conf"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/uuid"
	"strconv"
	"time"
)

const (
	sparkMasterDeployName  = "spark-master-deploy"
	sparkWorkerDeployName  = "spark-worker-deploy"
	sparkMasterServiceName = "spark-master"
	sparkWorkerServiceName = "spark-worker-service"
	sparkIngressName       = "spark-ingress"
)

// CreateSpark 为uid创建spark，masterReplicas默认1， masterReplicas默认2
func CreateSpark(u_id uint, masterReplicas int32, workerReplicas int32, expiredTime *time.Time, resources forms.Resources) (*responses.Response, error) {
	// 随机生成ssh密码
	//pwd := CreatePWD(8)
	//fmt.Println(pwd)

	// 获取当前时间戳，纳秒
	//s := strconv.FormatInt(time.Now().UnixNano(), 10)

	// uuid
	s := string(uuid.NewUUID())
	label := map[string]string{
		"image": "spark",
	}
	masterSelector := map[string]string{
		"component": "spark-master",
	}
	masterLabel := map[string]string{
		"component": "spark-master",
	}
	workerSelector := map[string]string{
		"component": "spark-worker",
	}
	workerLabel := map[string]string{
		"component": "spark-worker",
	}
	if u_id != 0 {
		uid := strconv.Itoa(int(u_id))
		label["u_id"] = uid
		//masterLabel["u_id"] = uid
		//workerLabel["u_id"] = uid
	}
	// 创建namespace
	_, err := CreateNs("spark-"+s, expiredTime, label, resources)
	if err != nil {
		return nil, err
	}

	// spark的master控制器
	masterSpec := appsv1.DeploymentSpec{
		Replicas: &masterReplicas,
		Selector: &metav1.LabelSelector{MatchLabels: masterSelector},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: masterLabel},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:            "spark-master",
						Image:           conf.SparkImage,
						ImagePullPolicy: corev1.PullIfNotPresent, // 镜像拉取策略
						Command:         []string{"/start-master"},
						Args:            []string{conf.SshPwd},
						Ports: []corev1.ContainerPort{
							{ContainerPort: 7077},
							{ContainerPort: 8080},
							{ContainerPort: 22},
						},
						//Resources: corev1.ResourceRequirements{
						//	Requests: corev1.ResourceList{
						//		corev1.ResourceCPU: resource.MustParse("100m"),
						//	},
						//	Limits: corev1.ResourceList{
						//		corev1.ResourceCPU:    resource.MustParse(cpu),
						//		corev1.ResourceMemory: resource.MustParse(memory),
						//	},
						//},
					},
				},
				RestartPolicy: corev1.RestartPolicyAlways,
			},
		},
	}
	_, err = CreateDeploy(sparkMasterDeployName, "spark-"+s, map[string]string{}, masterSpec)
	if err != nil {
		return nil, err
	}

	// spark的master的service
	masterServiceSpec := corev1.ServiceSpec{
		Type:     corev1.ServiceTypeNodePort,
		Selector: masterSelector,
		Ports: []corev1.ServicePort{ // 默认生成nodePort
			{Name: "spark", Port: 7077, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7077}},
			{Name: "http", Port: 8080, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8080}},
			{Name: "ssh", Port: 22, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 22}},
		},
	}
	_, err = CreateService(sparkMasterServiceName, "spark-"+s, map[string]string{}, masterServiceSpec)
	if err != nil {
		return nil, err
	}

	// spark的worker的控制器
	workerSpec := appsv1.DeploymentSpec{
		Replicas: &workerReplicas,
		Selector: &metav1.LabelSelector{MatchLabels: workerSelector},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: workerLabel},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "spark-worker",
						Image:           conf.SparkImage,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Command:         []string{"/start-worker"},
						Args:            []string{conf.SshPwd},
						Ports: []corev1.ContainerPort{
							{ContainerPort: 8081},
							{ContainerPort: 22},
						},
						//Resources: corev1.ResourceRequirements{
						//	Requests: corev1.ResourceList{
						//		corev1.ResourceCPU: resource.MustParse("100m"),
						//	},
						//	Limits: corev1.ResourceList{
						//		corev1.ResourceCPU:    resource.MustParse(cpu),
						//		corev1.ResourceMemory: resource.MustParse(memory),
						//	},
						//},
					},
				},
				RestartPolicy: corev1.RestartPolicyAlways,
			},
		},
	}
	_, err = CreateDeploy(sparkWorkerDeployName, "spark-"+s, map[string]string{}, workerSpec)
	if err != nil {
		return nil, err
	}

	// spark的worker的service
	workerServiceSpec := corev1.ServiceSpec{
		Type:     corev1.ServiceTypeNodePort,
		Selector: workerSelector,
		Ports: []corev1.ServicePort{
			{Name: "http", Port: 8081, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8081}},
			{Name: "ssh", Port: 22, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 22}},
		},
	}
	_, err = CreateService(sparkWorkerServiceName, "spark-"+s, map[string]string{}, workerServiceSpec)
	if err != nil {
		return nil, err
	}

	// spark的ingress
	ingressSpec := v1beta1.IngressSpec{
		Rules: []v1beta1.IngressRule{
			// 代理master服务
			{
				Host: fmt.Sprintf("spark.%s", conf.ProjectName),
				IngressRuleValue: v1beta1.IngressRuleValue{
					HTTP: &v1beta1.HTTPIngressRuleValue{
						Paths: []v1beta1.HTTPIngressPath{
							{
								Path:    "/master-" + s,
								Backend: v1beta1.IngressBackend{ServiceName: sparkMasterServiceName, ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 8080}},
							},
						},
					},
				},
			},
			// 代理worker服务
			{
				Host: fmt.Sprintf("spark.%s", conf.ProjectName),
				IngressRuleValue: v1beta1.IngressRuleValue{
					HTTP: &v1beta1.HTTPIngressRuleValue{
						Paths: []v1beta1.HTTPIngressPath{
							{
								Path:    "/worker-" + s,
								Backend: v1beta1.IngressBackend{ServiceName: sparkWorkerServiceName, ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 8081}},
							},
						},
					},
				},
			},
		},
	}
	_, err = CreateIngress(sparkIngressName, "spark-"+s, map[string]string{}, ingressSpec)
	if err != nil {
		return nil, err
	}

	return &responses.OK, nil
}

// GetSpark 获取uid用户下的所有spark
func GetSpark(u_id uint) (*responses.SparkListResponse, error) {
	label := map[string]string{
		"image": "spark",
	}
	if u_id > 0 {
		label["u_id"] = strconv.Itoa(int(u_id))
	}
	// 将map标签转换为string
	selector := labels.SelectorFromSet(label).String()
	sparks, err := GetNs(selector)
	if err != nil {
		return nil, err
	}
	sparkList := make([]responses.Spark, sparks.Length)
	for i, spark := range sparks.NsList {
		// 获取pod
		podList, err := GetPod(spark.Name, "")
		if err != nil {
			return nil, err
		}
		// 获取deploy
		deployList, err := GetDeploy(spark.Name, "")
		if err != nil {
			return nil, err
		}
		var master, worker int32
		for j := 0; j < deployList.Length; j++ {
			deploy := deployList.DeployList[j]
			if deploy.Name == sparkMasterDeployName {
				master = deploy.Replicas
			} else if deploy.Name == sparkWorkerDeployName {
				worker = deploy.Replicas
			}
		}
		// 获取service
		serviceList, err := GetService(spark.Name, "")
		if err != nil {
			return nil, err
		}
		// 获取ingress
		ingressList, err := GetIngress(spark.Name, "")
		if err != nil {
			return nil, err
		}
		sparkList[i] = responses.Spark{
			Name:           spark.Name,
			CreatedAt:      spark.CreatedAt,
			Uid:            u_id,
			Status:         spark.Status,
			Username:       spark.Username,
			Nickname:       spark.Nickname,
			PodList:        podList.PodList,
			DeployList:     deployList.DeployList,
			ServiceList:    serviceList.ServiceList,
			IngressList:    ingressList.IngressList,
			MasterReplicas: master,
			WorkerReplicas: worker,
			ExpiredTime:    spark.ExpiredTime,
			Cpu:            spark.Cpu,
			Memory:         spark.Memory,
			UsedCpu:        spark.UsedCpu,
			UsedMemory:     spark.UsedMemory,
		}
	}

	return &responses.SparkListResponse{
		Response:  responses.OK,
		Length:    sparks.Length,
		SparkList: sparkList,
	}, nil
}

// DeleteSpark 删除spark
func DeleteSpark(ns string) (*responses.Response, error) {
	var err1 error
	if _, err := DeleteIngress(sparkIngressName, ns); err != nil {
		err1 = err
	}
	if _, err := DeleteService(sparkWorkerServiceName, ns); err != nil {
		err1 = err
	}
	if _, err := DeleteService(sparkMasterServiceName, ns); err != nil {
		err1 = err
	}
	if _, err := DeleteDeploy(sparkWorkerDeployName, ns); err != nil {
		err1 = err
	}
	if _, err := DeleteDeploy(sparkMasterDeployName, ns); err != nil {
		err1 = err
	}
	if _, err := DeleteNs(ns); err != nil {
		err1 = err
	}
	if err1 != nil {
		return nil, err1
	}
	return &responses.OK, nil
}

// UpdateSpark 更新spark的uid以及replicas
func UpdateSpark(name, uid string, masterReplicas int32, workerReplicas int32, expiredTime *time.Time, resources forms.Resources) (*responses.Response, error) {
	if _, err := UpdateNs(name, uid, expiredTime, resources); err != nil {
		return nil, err
	}

	// 更新master的Replicas
	master, err := GetADeploy(sparkMasterDeployName, name)
	if err != nil {
		return nil, err
	}
	master.Spec.Replicas = &masterReplicas
	if _, err := UpdateDeploy(master); err != nil {
		return nil, err
	}

	// 更新worker的Replicas
	worker, err := GetADeploy(sparkWorkerDeployName, name)
	if err != nil {
		return nil, err
	}
	worker.Spec.Replicas = &workerReplicas
	if _, err := UpdateDeploy(worker); err != nil {
		return nil, err
	}

	return &responses.OK, nil
}
