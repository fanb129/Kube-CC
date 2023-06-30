package docker

import (
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"io"
	"os"
)

// PullImage
// 拉取指定镜像

// TODO 后续将方法规范重构
func PullImage(imageName string) (*responses.PullingResponse, error) {
	// 创建同步
	ctx := context.Background()

	reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})

	if err != nil {
		fmt.Println(err)
	}

	// TODO written方法需要后续经过测试进行完善

	// 读取所需要复制的内容
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return nil, err
	}

	image, err := dao.GetImgById(imageName)
	if err != nil {
		return nil, errors.New("镜像拉取失败")
	}

	return &responses.PullingResponse{
		Response: responses.OK,
		ImageInfo: responses.ImageInfo{
			ID:        image.ID,
			CreatedAt: image.CreatedAt,
			UpdatedAt: image.UpdatedAt,
			ImageId:   image.ImageId,
			UserId:    image.UserId,
			Kind:      image.Kind,
		},
	}, nil
}

// PullPrivateImage
// 拉取私有仓库的镜像
// TODO 测试后进行完善
/*func PullPrivateImage(imageName string, username string, passwd string) (*responses.Response, error) {
	ctx := context.Background()
	authConf := registry.AuthConfig{
		Username: username,
		Password: passwd,
	}
	encodeJson, _ := json.Marshal(authConf)
	authStr := base64.StdEncoding.EncodeToString(encodeJson)
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer out.Close()
	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return err
	}

	return nil
}
*/
