package docker

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
)

func DeleteImage(imageId string) ([]types.ImageDeleteResponseItem, error, *responses.Response) {
	ctx := context.Background()
	imgdeletersp, err := cli.ImageRemove(ctx, imageId, types.ImageRemoveOptions{Force: true})
	if err != nil {

		fmt.Println(err)
		return nil, err, nil
	}
	_, err = dao.DeletedImgByImageId(imageId)
	if err != nil {
		return nil, err, nil
	}
	return imgdeletersp, nil, &responses.OK
}
