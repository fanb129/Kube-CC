package main

import (
	"Kube-CC/dao"
	"Kube-CC/service/docker"
	"fmt"
)

func main() {
	// 集成连接docker进行初始化
	dao.InitDB()
	err := docker.ConnectDocker()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 输出检测

	//dockerfile
	//ar, err := os.ReadFile("D:\\GitHubProject\\Kube-CC\\service\\docker\\Dockerfile")

	// r := strings.NewReader("-f")

	//dkf := string(ar)
	// docker.CreateImageBytDockerFile(dkf)

	// [ADD] 拉取公有仓库测试
	//err = docker.PullImage("registry.cn-shanghai.aliyuncs.com/fanb/mycentos", "7", 23)
	//
	//if err != nil {
	//	fmt.Println(err)
	//}

	//通过账号密码拉取私有仓库镜像的测试
	//err = docker.PullPrivateImage("registry.cn-shanghai.aliyuncs.com/fanb/myspark:1.5.2_v1", "发以稀为贵", "5115219452.62.fb")

	// [REMOVE] 删除镜像
	//_, err, _ = docker.DeleteImage("564cecc5c0ae")
	//if err != nil {
	//	fmt.Println(err)
	//}

	// [ADD] 通过修改TAG创建镜像
	//docker.AlertTag("registry.cn-shanghai.aliyuncs.com/fanb/mycentos", "7", "registry.cn-shanghai.aliyuncs.com/fanb/mycentos", "7.6")
}
