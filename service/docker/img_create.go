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

// 获取要创建的镜像，创建成功则返回相应的表单
func CreateImage(parent, username, passwd, tag string, uid uint, kind int) (*responses.Response, error) {
	// 创建同步
	ctx := context.Background()
	authConf := registry.AuthConfig{
		Username: username,
		Password: passwd,
	}
	encodeJson, _ := json.Marshal(authConf)
	authStr := base64.StdEncoding.EncodeToString(encodeJson)
	out, err := cli.ImageCreate(ctx, parent, types.ImageCreateOptions{RegistryAuth: authStr})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer out.Close()

	// 读取所需要复制的内容
	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return nil, err
	}

	id, size, err := GetIDandSize(parent)
	if err != nil {
		return nil, err
	}
	_, err = dao.CreateImage(parent, id, uid, kind, tag, size)

	if err != nil {
		return nil, err
	}

	return &responses.OK, nil
}
