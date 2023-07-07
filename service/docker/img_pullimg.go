package docker

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"io"
	"os"
)

// PullImage
// 拉取指定公有仓库镜像

func PullImage(repositoryName, tag string, uid uint, kind int) (*responses.Response, error) {
	// 创建同步
	ctx := context.Background()

	reader, err := cli.ImagePull(ctx, repositoryName+":"+tag, types.ImagePullOptions{})

	if err != nil {
		fmt.Println(err)
	}

	// 读取所需要复制的内容
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return nil, err
	}

	id, size, err := GetIDandSize(repositoryName + ":" + tag)
	if err != nil {
		return nil, err
	}
	_, err = dao.CreateImage(repositoryName, id, uid, kind, tag, size)

	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}

// PullPrivateImage
// 拉取指定私有仓库的镜像
func PullPrivateImage(repositoryName, tag, username, passwd string, uid uint, kind int) (*responses.Response, error) {
	ctx := context.Background()
	authConf := registry.AuthConfig{
		Username: username,
		Password: passwd,
	}
	encodeJson, _ := json.Marshal(authConf)
	authStr := base64.StdEncoding.EncodeToString(encodeJson)
	out, err := cli.ImagePull(ctx, repositoryName, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer out.Close()
	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return nil, err
	}
	// 确保imagepull正确进行后进行数据库的写入操作

	id, size, err := GetIDandSize(repositoryName)
	if err != nil {
		return nil, err
	}
	_, err = dao.CreateImage(repositoryName, id, uid, kind, tag, size)

	if err != nil {
		return nil, err
	}

	return &responses.OK, nil
}
