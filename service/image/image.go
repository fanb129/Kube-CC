package image

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"Kube-CC/models"
	"go.uber.org/zap"
	"strings"
	"time"
)

var (
	sealosHub       = "sealos.hub:5000"
	sealosHubAdmin  = "admin"
	sealosHubPasswd = "passw0rd"
)

// ListImages 镜像管理界面列出镜像
func ListImages(uid, gid uint) (*responses.ImageListResponse, error) {
	var dockers []models.Docker
	var err error
	if uid > 0 { // 普通用户查看自己
		dockers, err = dao.GetImagesByUid(uid)
		if err != nil {
			return nil, err
		}
	} else {
		if gid > 0 { // 组长查看全组
			users, err := dao.GetGroupUserById(gid)
			if err != nil {
				return nil, err
			}
			for _, user := range users {
				tmp, err := dao.GetPrivateImage(user.ID)
				if err != nil {
					return nil, err
				}
				dockers = append(dockers, tmp...)
			}
			tmp1, err := dao.GetPublicImage()
			dockers = append(dockers, tmp1...)
		} else { // 管理员查看所有人
			dockers, err = dao.GetAllImages()
			if err != nil {
				return nil, err
			}
		}
	}

	imageList, length := addToImageList(dockers)

	return &responses.ImageListResponse{
		Response:  responses.OK,
		Length:    length,
		ImageList: *imageList,
	}, nil
}

// ListOkImages 添加应用时选择镜像
func ListOkImages(uid uint) (*responses.ImageListResponse, error) {
	dockers, err := dao.GetOkIMages(uid)
	if err != nil {
		return nil, err
	}
	imageList, length := addToImageList(dockers)

	return &responses.ImageListResponse{
		Response:  responses.OK,
		Length:    length,
		ImageList: *imageList,
	}, nil
}

// PullImage 拉取镜像到平台
func PullImage(form forms.PullImageForm) (*responses.Response, error) {
	id, err := dao.CreateImage(form.TargetImage.Name, form.TargetImage.Tag, form.TargetImage.Description, form.TargetImage.Uid, form.TargetImage.Kind)
	if err != nil {
		return nil, err
	}
	// 后台pull and push
	go backendPullAndPush(id, form.PullImage)
	return &responses.Response{StatusCode: 1, StatusMsg: "操作成功,等待后台下载"}, nil
}

// SaveImage 将当前pod容器状态保存为新的镜像
func SaveImage(form forms.SaveImageForm) (*responses.Response, error) {
	id, err := dao.CreateImage(form.TargetImage.Name, form.TargetImage.Tag, form.TargetImage.Description, form.TargetImage.Uid, form.TargetImage.Kind)
	if err != nil {
		return nil, err
	}

	// 后台save and push
	go backendSaveAndPush(id, form.ContainerID, form.NodeIp)

	return &responses.Response{StatusCode: 1, StatusMsg: "操作成功,等待后台下载"}, nil
}

func DeleteImage(id uint) (*responses.Response, error) {
	err := dao.DeleteImage(id)
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}

func UpdateImage(form forms.UpdateImageForm) (*responses.Response, error) {
	err := dao.UpdateImage(form.Id, form.Kind, form.Description)
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}

func addToImageList(dockers []models.Docker) (*[]responses.ImageInfo, int) {
	length := len(dockers)
	imageList := make([]responses.ImageInfo, length)

	for i, docker := range dockers {
		username := ""
		nickname := ""
		user, err := dao.GetUserById(docker.UserId)
		if err != nil {
			zap.S().Errorln(err)
		} else {
			username = user.Username
			nickname = user.Nickname
		}

		imageList[i] = responses.ImageInfo{
			Id:        docker.ID,
			CreatedAt: docker.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: docker.UpdatedAt.Format("2006-01-02 15:04:05"),

			ImageId: docker.ImageId,
			Name:    docker.ImageName,
			Tag:     docker.Tag,
			Size:    docker.Size,

			Username: username,
			Nickname: nickname,
			Uid:      docker.UserId,

			Kind:   docker.Kind,
			Status: docker.Status,

			Description: docker.Description,
		}
	}

	return &imageList, length
}

// 后台pull and push
func backendPullAndPush(id uint, pullImage forms.PullImage) {
	docker, err := dao.GetImageById(id)
	if err != nil {
		zap.S().Errorln(err)
		docker.Status = 3
		dao.SaveImage(docker)
		return
	}

	//TODO ip暂时写死
	//cli, err := NewDockerCli(conf.MasterInfo.Host, 2375)
	cli, err := NewDockerCli("192.168.139.143", 2375)
	if cli != nil {
		defer cli.Close()
	}
	if err != nil {
		zap.S().Errorln(err)
		docker.Status = 3
		dao.SaveImage(docker)
		return
	}

	// 1. pull 到本地
	err = cli.Pull(pullImage.Name+":"+pullImage.Tag, pullImage.Username, pullImage.Password)
	if err != nil {
		zap.S().Errorln(err)
		docker.Status = 3
		dao.SaveImage(docker)
		return
	}

	// 2. 对镜像进行tag操作
	err = cli.Tag(pullImage.Name+":"+pullImage.Tag, docker.ImageName+":"+docker.Tag)
	if err != nil {
		zap.S().Errorln(err)
		//docker.Status = 3
		//dao.SaveImage(docker)
		//return
	}
	size, err := cli.GetSize(docker.ImageName + ":" + docker.Tag)
	if err != nil {
		zap.S().Errorln(err)
	}
	docker.Size = size

	// 3. push到仓库
	backendPush(docker, cli)
}

func backendSaveAndPush(id uint, containerId, nodeIp string) {
	docker, err := dao.GetImageById(id)
	if err != nil {
		zap.S().Errorln(err)
		docker.Status = 3
		dao.SaveImage(docker)
		return
	}

	cli, err := NewDockerCli(nodeIp, 2375)
	if cli != nil {
		defer cli.Close()
	}
	if err != nil {
		zap.S().Errorln(err)
		docker.Status = 3
		dao.SaveImage(docker)
		return
	}

	// 1. commit到本地
	err = cli.Commit(containerId, docker.ImageName+":"+docker.Tag)
	if err != nil {
		zap.S().Errorln(err)
		docker.Status = 3
		dao.SaveImage(docker)
		return
	}
	size, err := cli.GetSize(docker.ImageName + ":" + docker.Tag)
	if err != nil {
		zap.S().Errorln(err)
	}
	docker.Size = size

	// 2. push到仓库
	backendPush(docker, cli)
}

// 修改tag并push到仓库
func backendPush(docker *models.Docker, cli *DockerCli) {
	// push到仓库
	err := cli.Push(docker.ImageName+":"+docker.Tag, sealosHubAdmin, sealosHubPasswd)
	if err != nil {
		zap.S().Errorln(err)
		docker.Status = 3
		dao.SaveImage(docker)
		return
	}

	// 完善docker信息
	index := strings.Index(docker.ImageName, "/")
	name := docker.ImageName[index+1:]
	manifest, err := dao.Hub.ManifestV2(name, docker.Tag)
	if err != nil {
		zap.S().Errorln(err, "3s后第一次重试")
		// 3秒后重试
		time.Sleep(time.Second * 3)
		manifest, err = dao.Hub.ManifestV2(name, docker.Tag)
		if err != nil {
			zap.S().Errorln(err, "10s后第二次重试")
			// 10秒后重试
			time.Sleep(time.Second * 10)
			manifest, err = dao.Hub.ManifestV2(name, docker.Tag)
			if err != nil {
				zap.S().Errorln(err)
				docker.Status = 3
				dao.SaveImage(docker)
				return
			}
		}
	}
	docker.Status = 1
	//docker.Size = getSize(manifest.Config.Size)
	docker.ImageId = getImageId(manifest.Config.Digest.String())
	dao.SaveImage(docker)
}

func getImageId(digest string) string {
	index := strings.Index(digest, ":")
	return digest[index+1 : index+13]
}
