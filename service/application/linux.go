package application

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/conf"
	"Kube-CC/service"
	"encoding/json"
	"errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/uuid"
)

var linuxImage = [2]string{"centos", "ubuntu"}
var cmd = [2][]string{{"/usr/sbin/init"}, {"/init.sh"}}
var privileged = [2]bool{true, false}

// CreateLinux 为uid创建linux 1-centos，2-ubuntu
func CreateLinux(name, ns string, kind uint, resources forms.ApplyResources) (*responses.Response, error) {
	// 将form序列化为string，存入deploy的注释
	form := forms.LinuxUpdateForm{
		Name:           name,
		Namespace:      ns,
		ApplyResources: resources,
	}
	jsonBytes, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}
	strForm := string(jsonBytes)
	// 准备工作
	// 分割申请资源
	requestCpu, err := service.SplitRSC(resources.Cpu, n)
	if err != nil {
		return nil, err
	}
	requestMemory, err := service.SplitRSC(resources.Memory, n)
	if err != nil {
		return nil, err
	}
	requestStorage, err := service.SplitRSC(resources.Storage, n)
	if err != nil {
		return nil, err
	}

	// 创建uuid，以便筛选出属于同一组的deploy、pod、service等
	newUuid := string(uuid.NewUUID())
	label := map[string]string{
		"image": linuxImage[kind-1],
		"uuid":  newUuid,
	}
	// 创建PVC，持久存储
	var volumes []corev1.Volume
	var volumeMounts []corev1.VolumeMount
	if resources.PvcStorage != "" {
		volumes = make([]corev1.Volume, 1)
		volumeMounts = make([]corev1.VolumeMount, 1)
		if resources.StorageClassName == "" {
			return nil, errors.New("已填写PvcStorage,StorageClassName不能为空")
		}
		pvcName := name + "-pvc"
		_, err = service.CreatePVC(ns, pvcName, resources.StorageClassName, resources.PvcStorage, readWriteMany)
		if err != nil {
			return nil, err
		}
		volumes[0] = corev1.Volume{
			Name: pvcName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvcName,
				},
			},
		}
		// 写死为/data目录
		volumeMounts[0] = corev1.VolumeMount{
			Name:      pvcName,
			MountPath: "/data",
		}
	}
	// 创建centos的控制器
	var replicas int32
	replicas = 1
	deploySpec := appsv1.DeploymentSpec{
		Replicas: &replicas,
		Selector: &metav1.LabelSelector{MatchLabels: label},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: label},
			Spec: corev1.PodSpec{
				Volumes:       volumes,
				RestartPolicy: corev1.RestartPolicyAlways,
				Containers: []corev1.Container{
					{
						Name:            linuxImage[kind-1],
						Image:           conf.LinuxImage[kind-1],
						ImagePullPolicy: corev1.PullIfNotPresent,
						Command:         cmd[kind-1],
						//Env:             []corev1.EnvVar{{Name: "mypwd", Value: pwd}},
						SecurityContext: &corev1.SecurityContext{Privileged: &privileged[0]}, // 以特权模式进入容器
						Args:            []string{conf.SshPwd},
						Ports: []corev1.ContainerPort{
							{ContainerPort: 22},
						},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:              resource.MustParse(requestCpu),
								corev1.ResourceMemory:           resource.MustParse(requestMemory),
								corev1.ResourceEphemeralStorage: resource.MustParse(requestStorage),
								//TODO GPU
							},
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:              resource.MustParse(resources.Cpu),
								corev1.ResourceMemory:           resource.MustParse(resources.Memory),
								corev1.ResourceEphemeralStorage: resource.MustParse(resources.Storage),
							},
						},
						VolumeMounts: volumeMounts,
					},
				},
			},
		},
	}
	_, err = service.CreateDeploy(name, ns, strForm, label, deploySpec)
	if err != nil {
		DeleteLinux(name, ns)
		return nil, err
	}

	// centos的service
	serviceSpec := corev1.ServiceSpec{
		Type:     corev1.ServiceTypeNodePort,
		Selector: label,
		Ports: []corev1.ServicePort{
			{Name: "ssh", Port: 22, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 22}},
		},
	}
	_, err = service.CreateService(name+"-service", ns, label, serviceSpec)
	if err != nil {
		DeleteLinux(name, ns)
		return nil, err
	}
	return &responses.OK, nil
}

// ListLinux 获取uid用户下的所有kind类型的linux
func ListLinux(ns string, kind uint) (*responses.AppDeployList, error) {
	label := map[string]string{
		"image": linuxImage[kind-1],
	}
	// 将map标签转换为string
	selector := labels.SelectorFromSet(label).String()
	list, err := ListAppDeploy(ns, selector)
	if err != nil {
		return nil, err
	}
	return list, nil

}

// DeleteLinux 删除linux
func DeleteLinux(name, ns string) (*responses.Response, error) {
	response, err := DeleteAppDeploy(name, ns)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// GetLinux 更新之前先get
func GetLinux(name, ns string) (*forms.LinuxUpdateForm, error) {
	form := forms.LinuxUpdateForm{}
	deploy, err := service.GetDeploy(name, ns)
	if err != nil {
		return nil, err
	}
	strForm := deploy.Annotations["form"]
	err = json.Unmarshal([]byte(strForm), &form)
	if err != nil {
		return nil, err
	}
	return &form, nil
}

func UpdateLinux(name, ns string, resources forms.ApplyResources) (*responses.Response, error) {
	// 将form序列化为string，存入deploy的注释
	form := forms.LinuxUpdateForm{
		Name:           name,
		Namespace:      ns,
		ApplyResources: resources,
	}
	jsonBytes, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}
	strForm := string(jsonBytes)
	// 准备工作
	// 分割申请资源
	requestCpu, err := service.SplitRSC(resources.Cpu, n)
	if err != nil {
		return nil, err
	}
	requestMemory, err := service.SplitRSC(resources.Memory, n)
	if err != nil {
		return nil, err
	}
	requestStorage, err := service.SplitRSC(resources.Storage, n)
	if err != nil {
		return nil, err
	}
	pvcName := name + "-pvc"

	// 更新deploy
	deploy, err := service.GetDeploy(name, ns)
	if err != nil {
		return nil, err
	}
	// 创建PVC，持久存储
	var volumes []corev1.Volume
	var volumeMounts []corev1.VolumeMount
	if resources.PvcStorage != "" {
		volumes = make([]corev1.Volume, 1)
		volumeMounts = make([]corev1.VolumeMount, 1)
		_, err = service.UpdateOrCreatePvc(ns, pvcName, resources.StorageClassName, resources.PvcStorage, readWriteMany)
		if err != nil {
			return nil, err
		}
		volumes[0] = corev1.Volume{
			Name: pvcName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvcName,
				},
			},
		}
		// 写死为/data目录
		volumeMounts[0] = corev1.VolumeMount{
			Name:      pvcName,
			MountPath: "/data",
		}
	}

	deploySpec := deploy.Spec
	deploySpec.Template.Spec.Containers[0].Resources = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:              resource.MustParse(requestCpu),
			corev1.ResourceMemory:           resource.MustParse(requestMemory),
			corev1.ResourceEphemeralStorage: resource.MustParse(requestStorage),
			//TODO GPU
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:              resource.MustParse(resources.Cpu),
			corev1.ResourceMemory:           resource.MustParse(resources.Memory),
			corev1.ResourceEphemeralStorage: resource.MustParse(resources.Storage),
		},
	}
	deploySpec.Template.Spec.Containers[0].VolumeMounts = volumeMounts
	deploySpec.Template.Spec.Volumes = volumes

	if _, err := service.UpdateDeploy(name, ns, strForm, deploySpec); err != nil {
		return nil, err
	}

	return &responses.OK, nil
}
