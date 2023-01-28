package service

import (
	"Kube-CC/common"
	"Kube-CC/conf"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/uuid"
	"strconv"
	"strings"
	"time"
)

var linuxImage = [2]string{"centos", "ubuntu"}
var cmd = [2][]string{{"/usr/sbin/init"}, {"/init.sh"}}
var privileged = [2]bool{true, false}

// CreateLinux 为uid创建linux 1-centos，2-ubuntu
func CreateLinux(u_id, kind uint, expiredTime *time.Time, cpu, memory string) (*common.Response, error) {

	// 随机生成ssh密码
	//pwd := CreatePWD(8)
	//fmt.Println(pwd)

	// uuid
	s := string(uuid.NewUUID())
	selector := map[string]string{
		"image": linuxImage[kind-1],
	}
	label := map[string]string{
		"image": linuxImage[kind-1],
	}
	if u_id != 0 {
		uid := strconv.Itoa(int(u_id))
		label["u_id"] = uid
	}
	// 创建namespace
	_, err := CreateNs(linuxImage[kind-1]+"-"+s, expiredTime, label, cpu, memory, 1)
	if err != nil {
		return nil, err
	}

	// 创建centos的控制器
	var replicas int32

	replicas = 1
	deploySpec := appsv1.DeploymentSpec{
		Replicas: &replicas,
		Selector: &metav1.LabelSelector{MatchLabels: selector},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: label},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:            linuxImage[kind-1],
						Image:           conf.LinuxImage[kind-1],
						ImagePullPolicy: corev1.PullIfNotPresent,
						Command:         cmd[kind-1],
						//Env:             []corev1.EnvVar{{Name: "mypwd", Value: pwd}},
						SecurityContext: &corev1.SecurityContext{Privileged: &privileged[kind-1]}, // 以特权模式进入容器
						Args:            []string{conf.SshPwd},
						Ports: []corev1.ContainerPort{
							{ContainerPort: 22},
						},
						//Resources: corev1.ResourceRequirements{
						//	Requests: corev1.ResourceList{
						//		corev1.ResourceCPU:    resource.MustParse("1m"),
						//		corev1.ResourceMemory: resource.MustParse("1m"),
						//	},
						//	Limits: corev1.ResourceList{
						//		corev1.ResourceCPU:    resource.MustParse(cpu),
						//		corev1.ResourceMemory: resource.MustParse(memory),
						//	},
						//},
						//Resources: corev1.ResourceRequirements{
						//	Requests: corev1.ResourceList{
						//		corev1.ResourceCPU: resource.MustParse("100m"),
						//	},
						//},
					},
				},
			},
		},
	}
	_, err = CreateDeploy(linuxImage[kind-1]+"-deploy", linuxImage[kind-1]+"-"+s, label, deploySpec)
	if err != nil {
		return nil, err
	}

	// centos的service
	serviceSpec := corev1.ServiceSpec{
		Type:     corev1.ServiceTypeNodePort,
		Selector: selector,
		Ports: []corev1.ServicePort{
			{Name: "ssh", Port: 22, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 22}},
		},
	}
	_, err = CreateService(linuxImage[kind-1]+"-service", linuxImage[kind-1]+"-"+s, label, serviceSpec)
	if err != nil {
		return nil, err
	}
	return &common.OK, nil
}

// GetLinux 获取uid用户下的所有kind类型的linux
func GetLinux(u_id, kind uint) (*common.LinuxListResponse, error) {
	label := map[string]string{
		"image": linuxImage[kind-1],
	}
	if u_id > 0 {
		label["u_id"] = strconv.Itoa(int(u_id))
	}
	// 将map标签转换为string
	selector := labels.SelectorFromSet(label).String()
	ns, err := GetNs(selector)
	if err != nil {
		return nil, err
	}
	LinuxList := make([]common.Linux, ns.Length)

	for i, linux := range ns.NsList {
		//获取pod
		podList, err := GetPod(linux.Name, "")
		if err != nil {
			return nil, err
		}
		// 获取deploy
		deployList, err := GetDeploy(linux.Name, "")
		if err != nil {
			return nil, err
		}
		// 获取service
		serviceList, err := GetService(linux.Name, "")
		if err != nil {
			return nil, err
		}
		LinuxList[i] = common.Linux{
			Name:        linux.Name,
			Uid:         u_id,
			Username:    linux.Username,
			Nickname:    linux.Nickname,
			Status:      linux.Status,
			CreatedAt:   linux.CreatedAt,
			PodList:     podList.PodList,
			DeployList:  deployList.DeployList,
			ServiceList: serviceList.ServiceList,
			ExpiredTime: linux.ExpiredTime,
			Cpu:         linux.Cpu,
			Memory:      linux.Memory,
			UsedCpu:     linux.UsedCpu,
			UsedMemory:  linux.UsedMemory,
		}
	}
	return &common.LinuxListResponse{
		Response:  common.OK,
		Length:    ns.Length,
		Image:     linuxImage[kind-1],
		LinuxList: LinuxList,
	}, nil
}

// DeleteLinux 删除linux
func DeleteLinux(ns string) (*common.Response, error) {
	image := strings.Split(ns, "-")[0]
	var err1 error
	if _, err := DeleteService(image+"-service", ns); err != nil {
		err1 = err
	}
	if _, err := DeleteDeploy(image+"-deploy", ns); err != nil {
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

func UpdateLinux(name, uid string, expiredTime *time.Time, cpu, memory string) (*common.Response, error) {
	res, err := UpdateNs(name, uid, expiredTime, cpu, memory, 1)
	return res, err
}
