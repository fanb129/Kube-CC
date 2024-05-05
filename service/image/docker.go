package image

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

type DockerCli struct {
	ip   string
	port int
	cli  *client.Client
}

func NewDockerCli(ip string, port int) (*DockerCli, error) {
	dockerCli := DockerCli{
		ip:   ip,
		port: port,
	}
	host := fmt.Sprintf("tcp://%s:%d", dockerCli.ip, dockerCli.port)
	var err error
	dockerCli.cli, err = client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.WithHost(host))
	if err != nil {
		zap.S().Errorln(err)
		return nil, err
	}
	return &dockerCli, nil
}

func (d *DockerCli) Close() error {
	var err error
	if d.cli != nil {
		err = d.cli.Close()
	}
	return err
}

// Pull 拉取共有或者私有镜像
func (d *DockerCli) Pull(repositoryName, username, password string) error {
	options := types.ImagePullOptions{}
	// 从私有仓库拉取时
	if username != "" && password != "" {
		authConf := registry.AuthConfig{
			Username: username,
			Password: password,
		}
		encodeJson, _ := json.Marshal(authConf)
		authStr := base64.StdEncoding.EncodeToString(encodeJson)
		options.RegistryAuth = authStr
	}
	_, err := d.cli.ImagePull(context.Background(), repositoryName, options)
	return err
}

//func (d *DockerCli) Remove(imageId string) error{
//	remove, err := d.cli.ImageRemove(context.Background(), imageId, types.ImageRemoveOptions{Force: true})
//
//}

// Tag 重命名
func (d *DockerCli) Tag(source, target string) error {
	err := d.cli.ImageTag(context.Background(), source, target)
	return err
}

func (d *DockerCli) Push(repositoryName, username, password string) error {
	// 设置账号密码
	pushOptions := types.ImagePushOptions{}
	authConf := registry.AuthConfig{
		Username: username,
		Password: password,
	}
	encodeJson, _ := json.Marshal(authConf)
	authStr := base64.StdEncoding.EncodeToString(encodeJson)
	pushOptions.RegistryAuth = authStr

	_, err := d.cli.ImagePush(context.Background(), repositoryName, pushOptions)
	//io.Copy(os.Stdout, reader)
	//defer reader.Close()
	return err
}

func (d *DockerCli) GetSize(repositoryName string) (string, error) {
	filter := filters.NewArgs()
	filter.Add("reference", repositoryName)
	images, err := d.cli.ImageList(context.Background(), types.ImageListOptions{
		Filters: filter,
	})
	if err != nil {
		return "", err
	}
	for _, image := range images {
		size := float64(image.Size) / (1000 * 1000) // 将大小转换为兆字节（MB）
		sizeFormatted := fmt.Sprintf("%.1fMB", size)
		return sizeFormatted, nil
	}

	return "", fmt.Errorf("image not found: %s", repositoryName)
}

// Commit 将容器保存为镜像到本地
func (d *DockerCli) Commit(containerID, target string) error {
	response, err := d.cli.ContainerCommit(context.Background(), containerID, types.ContainerCommitOptions{
		Reference: target,
	})
	if err != nil {
		return err
	}

	// 提取保存的镜像 ID
	imageID := response.ID

	// 保存镜像到本地文件系统
	_, err = d.cli.ImageSave(context.Background(), []string{imageID})
	return err
}
