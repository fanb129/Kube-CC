package docker

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	"github.com/docker/docker/api/types"
	"io"
	"os"
)

// 获取要创建的镜像，创建成功则返回相应的表单
func CreateImage(parent string) (*responses.ImageInfoResponse, error) {
	ctx := context.Background()
	reader, err := cli.ImageCreate(ctx, parent, types.ImageCreateOptions{})
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(os.Stdout, reader)

	if err != nil {
		return nil, err
	}
	image, err := dao.GetImgById(parent)
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
	}, nil
}
