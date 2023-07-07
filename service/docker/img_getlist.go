package docker

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	"github.com/docker/docker/api/types"
)

// GetImageList
// 获取当前用户镜像列表
// 备选参数
func GetImageListAll() (error, []types.ImageSummary) {
	// 设置同步信号
	ctx := context.Background()
	// images为获取到的[]types ImageSummary类型的存储镜像的特殊数据结构
	// imageList只能获取全部

	images, err := cli.ImageList(ctx, types.ImageListOptions{All: true}) // 等价于指令docker images
	if err != nil {
		return err, nil
	}
	return nil, images
}

// 获取指定镜像详细信息
func GetImage(imageID string) (response *responses.ImageInfoResponse, err error) {
	// 设置同步信号
	//ctx := context.Background()

	images, err := dao.GetImgById(imageID)

	if err != nil {
		return nil, err
	}

	return &responses.ImageInfoResponse{
		Response: responses.OK,
		ImageInfo: responses.ImageInfo{
			ID:        images.ID,
			CreatedAt: images.CreatedAt,
			UpdatedAt: images.UpdatedAt,
			ImageId:   images.ImageId,
			ImageName: images.ImageName,
			UserId:    images.UserId,
			Kind:      images.Kind,
			Tag:       images.Tag,
			Size:      images.Size,
		},
	}, nil
}
