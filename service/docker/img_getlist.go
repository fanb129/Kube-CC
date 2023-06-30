package docker

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	"github.com/docker/docker/api/types"
)

// GetImageList
// 获取镜像列表
func GetImageList(uid string) (*responses.ImageInfoResponse, error, []types.ImageSummary) {
	// 设置同步信号
	ctx := context.Background()
	//
	images, err := cli.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		return nil, err, nil
	}

	image, err := dao.GetImgById(uid)
	return &responses.ImageInfoResponse{
		Response: responses.OK,
		ImageInfo: responses.ImageInfo{
			ID:        image.ID,
			CreatedAt: image.CreatedAt,
			UpdatedAt: image.UpdatedAt,
			ImageId:   image.ImageId,
			UserId:    image.UserId,
			Kind:      image.Kind,
		},
	}, nil, images
}

// 获取指定镜像相关信息
func GetImage(imageID string) (imageInfo types.ImageInspect, err error) {
	// 设置同步信号
	ctx := context.Background()
	imageInfo, _, err = cli.ImageInspectWithRaw(ctx, imageID)
	if err != nil {
		return imageInfo, err
	}
	return imageInfo, nil
}
