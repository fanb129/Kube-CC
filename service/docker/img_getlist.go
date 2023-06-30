package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
)

// GetImageList
// 获取镜像列表
func GetImageList() error {
	ctx := context.Background()
	images, err := cli.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		return err
	}
	//打印结果
	for _, image := range images {
		fmt.Printf("================\n%q %q\n", image.RepoTags, image.ID)
	}
	return nil
}

// 获取指定镜像相关信息
func GetImage(imageID string) (imageInfo types.ImageInspect, err error) {
	ctx := context.Background()
	imageInfo, _, err = cli.ImageInspectWithRaw(ctx, imageID)
	if err != nil {
		return imageInfo, err
	}
	return imageInfo, nil
}
