package service

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListSc 获取所有的storageClass
func ListSc() (*responses.ScListResponse, error) {
	list, err := dao.ClientSet.StorageV1().StorageClasses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	num := len(list.Items)
	scList := make([]responses.Sc, num)
	for i, sc := range list.Items {
		allow := sc.AllowVolumeExpansion
		if allow == nil {
			f := false
			allow = &f
		}
		tmp := responses.Sc{
			Name:                 sc.Name,
			CreatedAt:            sc.CreationTimestamp.Format("2006-01-02 15:04:05"),
			ReclaimPolicy:        string(*sc.ReclaimPolicy),
			Provisioner:          sc.Provisioner,
			VolumeBindingMode:    string(*sc.VolumeBindingMode),
			AllowVolumeExpansion: *allow,
		}
		scList[i] = tmp
	}
	return &responses.ScListResponse{Response: responses.OK, Length: num, PvcList: scList}, nil
}

// CreateJivaSc 创建jiva sc
func CreateJivaSc(name string) (*storagev1.StorageClass, error) {
	policy := corev1.PersistentVolumeReclaimDelete //回收策略
	sc := &storagev1.StorageClass{
		TypeMeta:      metav1.TypeMeta{Kind: "StorageClass", APIVersion: "storage.k8s.io/v1"},
		ObjectMeta:    metav1.ObjectMeta{Name: name},
		Provisioner:   "openebs.io/provisioner-iscsi",
		ReclaimPolicy: &policy,
	}
	storageClass, err := dao.ClientSet.StorageV1().StorageClasses().Create(context.Background(), sc, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return storageClass, nil
}
