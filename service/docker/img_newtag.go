package docker

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	"errors"
)

func AddNewTag(ImageName string, data forms.ImageUpdateForm) (*responses.Response, error) {
	ctx := context.Background()
	image, err := dao.GetImgById(ImageName)
	if err != nil {
		return nil, errors.New("修改镜像名称失败")
	}
	image.ImageId = data.ImageId

	errs := cli.ImageTag(ctx, ImageName, data.ImageId)

	if errs != nil {
		return nil, errs
	}

	row, err := dao.UpdateImage(image)

	if err != nil || row == 0 {
		return nil, errors.New("镜像tag更新失败")
	}
	return &responses.OK, nil
}
