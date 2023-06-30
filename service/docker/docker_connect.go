package docker

import (
	"fmt"
	"github.com/docker/docker/client"
)

// ConnectDocker
// 链接docker
var cli *client.Client

func ConnectDocker() (err error) {
	cli, err = client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.WithHost("tcp://192.168.239.16:2375"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
