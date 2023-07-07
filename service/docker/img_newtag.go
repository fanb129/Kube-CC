package docker

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	"errors"
)

func AlertTag(RepositoryName, oldTag, NewRepositoryName, newTag string) (*responses.Response, error) {
	ctx := context.Background()
	errs := cli.ImageTag(ctx, RepositoryName+":"+oldTag, NewRepositoryName+":"+newTag)

	if errs != nil {
		return nil, errs
	}

	image, err := dao.GetImgByName(RepositoryName)
	image.Tag = newTag
	row, err := dao.CreateImageByTag(image)

	if err != nil || row == 0 {
		return nil, errors.New("镜像tag更新失败")
	}
	return &responses.OK, nil
}
