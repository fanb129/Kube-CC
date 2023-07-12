package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
)

// TODO 补充dao相关操作
func CreateImageBytDockerFile(DockerFile string) error {
	// 创建同步
	ctx := context.Background()

	_, err := cli.ImageBuild(ctx, nil, types.ImageBuildOptions{Dockerfile: DockerFile})

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
