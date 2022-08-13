package service

import (
	"errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/dao"
	"strconv"
	"time"
)

// CreateSpark 为uid创建spark，masterReplicas默认1， masterReplicas默认3
func CreateSpark(u_id uint, masterReplicas int32, workerReplicas int32) (*common.Response, error) {
	// 获取当前时间戳，纳秒
	s := strconv.FormatInt(time.Now().UnixNano(), 10)
	uid := strconv.Itoa(int(u_id))
	label := map[string]string{
		"time": s,
		"u_id": uid,
	}
	masterLabel := map[string]string{
		"component": "spark-master",
		"time":      s,
		"u_id":      uid,
	}
	workerLabel := map[string]string{
		"component": "spark-worker",
		"time":      s,
		"u_id":      uid,
	}
	// 创建namespace
	//_, err := CreateNs("spark", map[string]string{})
	//if err != nil {
	//	return nil, err
	//}

	// 记录当前用户新建的spark
	row, err := dao.CreateSpark(u_id, s)
	if err != nil || row == 0 {
		return nil, errors.New("创建spark失败")
	}

	// spark的master控制器
	masterSpec := appsv1.DeploymentSpec{
		Replicas: &masterReplicas,
		Selector: &metav1.LabelSelector{MatchLabels: masterLabel},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: masterLabel},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "spark-master",
						Image:           "registry.aliyuncs.com/google_containers/spark:1.5.2_v1",
						ImagePullPolicy: corev1.PullIfNotPresent, // 镜像拉取策略
						Command:         []string{"[/start-master]"},
						Ports: []corev1.ContainerPort{
							{ContainerPort: 7077},
							{ContainerPort: 8080},
						},
						//Resources: corev1.ResourceRequirements{Requests: }
					},
				},
			},
		},
	}
	_, err = CreateDeploy("spark-master-deploy"+s, "spark", label, masterSpec)
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
						Image:           "registry.aliyuncs.com/google_containers/spark:1.5.2_v1",
						ImagePullPolicy: corev1.PullIfNotPresent,
						Command:         []string{"[/start-worker]"},
						Ports:           []corev1.ContainerPort{{ContainerPort: 8081}},
					},
				},
			},
		},
	}
	_, err = CreateDeploy("spark-worker-deploy"+s, "spark", label, workerSpec)
	if err != nil {
		return nil, err
	}

	// spark的master的service
	masterServiceSpec := corev1.ServiceSpec{
		Type:     corev1.ServiceTypeNodePort,
		Selector: masterLabel,
		Ports: []corev1.ServicePort{
			{Name: "spark", Port: 7077, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7077}},
			{Name: "http", Port: 8080, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8080},
				Protocol: corev1.ProtocolTCP}, // 默认生成nodePort
		},
	}
	_, err = CreateService("spark-master-service"+s, "spark", label, masterServiceSpec)
	if err != nil {
		return nil, err
	}

	// spark的worker的service
	workerServiceSpec := corev1.ServiceSpec{
		Type:     corev1.ServiceTypeNodePort,
		Selector: workerLabel,
		Ports: []corev1.ServicePort{
			{Name: "http", Port: 8081, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8081},
				Protocol: corev1.ProtocolTCP},
		},
	}
	_, err = CreateService("spark-worker-service"+s, "spark", label, workerServiceSpec)
	if err != nil {
		return nil, err
	}

	return &common.OK, nil
}

// GetSpark 获取uid用户下的所有spark
func GetSpark(u_id uint) (*common.SparkListResponse, error) {
	sparks, err := dao.GetSparkListById(u_id)
	if err != nil {
		return nil, err
	}
	n := len(sparks)
	uid := strconv.Itoa(int(u_id))
	sparkList := make([]common.Spark, n)
	for i, spark := range sparks {
		t := spark.Time
		label := map[string]string{
			"u_id": uid,
			"time": t,
		}
		// 将map标签转换为string
		selector := labels.SelectorFromSet(label).String()
		// 获取pod
		podList, err := GetPod("spark", selector)
		if err != nil {
			return nil, err
		}
		// 获取deploy
		deployList, err := GetDeploy("spark", selector)
		if err != nil {
			return nil, err
		}
		// 获取service
		serviceList, err := GetService("spark", selector)
		if err != nil {
			return nil, err
		}
		sparkList[i] = common.Spark{
			Name:        "spark" + t,
			Uid:         u_id,
			Sid:         spark.ID,
			PodList:     podList.PodList,
			DeployList:  deployList.DeployList,
			ServiceList: serviceList.ServiceList,
		}
	}

	return &common.SparkListResponse{
		Response:  common.OK,
		Length:    n,
		SparkList: sparkList,
	}, nil
}

// DeleteSpark 删除spark
func DeleteSpark(s_id uint) (*common.Response, error) {
	spark, err := dao.GetSpark(s_id)
	if err != nil {
		return nil, err
	}
	t := spark.Time
	if _, err := DeleteService("spark-worker-service"+t, "spark"); err != nil {
		return nil, err
	}
	if _, err := DeleteService("spark-master-service"+t, "spark"); err != nil {
		return nil, err
	}
	if _, err := DeleteDeploy("spark-worker-deploy"+t, "spark"); err != nil {
		return nil, err
	}
	if _, err := DeleteDeploy("spark-master-deploy"+t, "spark"); err != nil {
		return nil, err
	}
	if row, err := dao.DeleteSpark(spark.ID); err != nil || row == 0 {
		return nil, errors.New("spark数据库删除失败")
	}
	return &common.OK, nil
}
