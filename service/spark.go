package service

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/conf"
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
func CreateSpark(u_id uint, masterReplicas int32, workerReplicas int32) (*common.Response, error) {
	// 随机生成ssh密码
	//pwd := CreatePWD(8)
	//fmt.Println(pwd)
	// 获取当前时间戳，纳秒
	s := strconv.FormatInt(time.Now().UnixNano(), 10)
	uid := strconv.Itoa(int(u_id))
	label := map[string]string{
		"image": "spark",
		"u_id":  uid,
	}
	masterLabel := map[string]string{
		"component": "spark-master",
		"u_id":      uid,
	}
	workerLabel := map[string]string{
		"component": "spark-worker",
		"u_id":      uid,
	}

	// 创建namespace
	_, err := CreateNs("spark-"+s, label)
	if err != nil {
		return nil, err
	}

	// spark的master控制器
	masterSpec := appsv1.DeploymentSpec{
		Replicas: &masterReplicas,
		Selector: &metav1.LabelSelector{MatchLabels: masterLabel},
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
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU: resource.MustParse("100m"),
							},
						},
					},
				},
			},
		},
	}
	_, err = CreateDeploy(sparkMasterDeployName, "spark-"+s, label, masterSpec)
	if err != nil {
		return nil, err
	}

	// spark的master的service
	masterServiceSpec := corev1.ServiceSpec{
		Type:     corev1.ServiceTypeNodePort,
		Selector: masterLabel,
		Ports: []corev1.ServicePort{ // 默认生成nodePort
			{Name: "spark", Port: 7077, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7077}},
			{Name: "http", Port: 8080, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8080}},
			{Name: "ssh", Port: 22, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 22}},
		},
	}
	_, err = CreateService(sparkMasterServiceName, "spark-"+s, label, masterServiceSpec)
	if err != nil {
		return nil, err
	}

	// spark的worker的控制器
	workerSpec := appsv1.DeploymentSpec{
		Replicas: &workerReplicas,
		Selector: &metav1.LabelSelector{MatchLabels: workerLabel},
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
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU: resource.MustParse("100m"),
							},
						},
					},
				},
			},
		},
	}
	_, err = CreateDeploy(sparkWorkerDeployName, "spark-"+s, label, workerSpec)
	if err != nil {
		return nil, err
	}

	// spark的worker的service
	workerServiceSpec := corev1.ServiceSpec{
		Type:     corev1.ServiceTypeNodePort,
		Selector: workerLabel,
		Ports: []corev1.ServicePort{
			{Name: "http", Port: 8081, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8081}},
			{Name: "ssh", Port: 22, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 22}},
		},
	}
	_, err = CreateService(sparkWorkerServiceName, "spark-"+s, label, workerServiceSpec)
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
	_, err = CreateIngress(sparkIngressName, "spark-"+s, label, ingressSpec)
	if err != nil {
		return nil, err
	}

	return &common.OK, nil
}

// GetSpark 获取uid用户下的所有spark
func GetSpark(u_id uint) (*common.SparkListResponse, error) {
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
	sparkList := make([]common.Spark, sparks.Length)
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
		sparkList[i] = common.Spark{
			Name:        spark.Name,
			CreatedAt:   spark.CreatedAt,
			Uid:         u_id,
			Username:    spark.Username,
			Nickname:    spark.Nickname,
			PodList:     podList.PodList,
			DeployList:  deployList.DeployList,
			ServiceList: serviceList.ServiceList,
			IngressList: ingressList.IngressList,
		}
	}

	return &common.SparkListResponse{
		Response:  common.OK,
		Length:    sparks.Length,
		SparkList: sparkList,
	}, nil
}

// DeleteSpark 删除spark
func DeleteSpark(ns string) (*common.Response, error) {
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
	return &common.OK, nil
}
