package docker

import (
	"Kube-CC/conf"
	"fmt"
	"github.com/docker/docker/client"
)

// ConnectDocker
// 链接docker
var cli *client.Client

func ConnectDocker() (err error) {
	// 连接
	//TODO 改成自由登录 \ 或者改成对应物理IP使用其他方法进行进一步处理
	cli, err = client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.WithHost("tcp://"+conf.DockerHost))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
