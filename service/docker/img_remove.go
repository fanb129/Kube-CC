package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
)

func DeleteImage(imageId string) (imageDeleteResponseItems []types.ImageDeleteResponseItem, err error) {
	ctx := context.Background()
	imageDeleteResponseItems, err = cli.ImageRemove(ctx, imageId, types.ImageRemoveOptions{Force: true})
	if err != nil {
		fmt.Println(err)
		return imageDeleteResponseItems, err
	}
	return imageDeleteResponseItems, nil
}
