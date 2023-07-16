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
	"go.uber.org/zap"
	"io"
	"os"
)

// PullImage
// 拉取指定公有仓库镜像

func PullImage(repositoryName, tag string, uid uint, kind int) (*responses.Response, error) {
	// 协程后台下载
	go func() {
		// 创建同步
		ctx := context.Background()
		// TODO 通过api实现pull镜像到指定仓库地址
		reader, err := cli.ImagePull(ctx, repositoryName+":"+tag, types.ImagePullOptions{})

		if err != nil {
			zap.S().Errorln(err)
			return
		}

		// 读取所需要复制的内容
		_, err = io.Copy(os.Stdout, reader)
		if err != nil {
			zap.S().Errorln(err)
			return
		}

		id, size, err := GetIDandSize(repositoryName + ":" + tag)
		if err != nil {
			zap.S().Errorln(err)
			return
		}
		_, err = dao.CreateImage(repositoryName, id, uid, kind, tag, size)

		if err != nil {
			zap.S().Errorln(err)
			return
		}
	}()

	return &responses.OK, nil
}

// PullPrivateImage
// 拉取指定私有仓库的镜像
func PullPrivateImage(repositoryName, tag, username, passwd string, uid uint, kind int) (*responses.Response, error) {
	go func() {
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
			zap.S().Errorln(err)
			return
		}
		defer out.Close()
		_, err = io.Copy(os.Stdout, out)
		if err != nil {
			zap.S().Errorln(err)
			return
		}
		// 确保imagepull正确进行后进行数据库的写入操作

		id, size, err := GetIDandSize(repositoryName)
		if err != nil {
			zap.S().Errorln(err)
			return
		}
		_, err = dao.CreateImage(repositoryName, id, uid, kind, tag, size)

		if err != nil {
			zap.S().Errorln(err)
			return
		}
	}()
	return &responses.OK, nil
}
