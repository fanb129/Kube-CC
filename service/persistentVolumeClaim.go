package service

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreatePVC 创建pvc
// storageClassName:选择的存储类 默认hostpath
// storageSize：申请的存储大小 5Gi 5G
// accessModes: 选择的读写模式 ReadWriteOnce ReadOnlyMany ReadWriteMany ReadWriteOncePod
func CreatePVC(namespace, name, storageClassName string, storageSize, accessModes string) (*responses.Response, error) {
	pvc := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			AccessModes:      []v1.PersistentVolumeAccessMode{v1.PersistentVolumeAccessMode(accessModes)},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(storageSize),
				},
			},
		},
	}

	_, err := dao.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Create(context.TODO(), pvc, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return &responses.OK, nil
}

// DeletePVC 删除pvc
func DeletePVC(namespace, name string) (*responses.Response, error) {
	err := dao.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}

func GetPVC(namespace, name string) (*v1.PersistentVolumeClaim, error) {
	pvc, err := dao.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return pvc, nil
}

// UpdatePVC 更新PVC,只能更改申请的存储大小
func UpdatePVC(namespace, name, newStorageSize string) error {
	pvc, err := dao.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	pvc.Spec.Resources.Requests[v1.ResourceStorage] = resource.MustParse(newStorageSize)

	_, err = dao.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Update(context.TODO(), pvc, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
