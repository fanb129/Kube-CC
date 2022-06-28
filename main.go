package main

import (
	"k8s_deploy_gin/conf"
	"k8s_deploy_gin/dao"
	"k8s_deploy_gin/routers"
)

func main() {
	dao.InitDB()   //数据库初始化
	dao.InitKube() // client-go k8s初始化
	//dao.InitRedisPool() //Redis 初始化(暂时不用)
	r := routers.InitRouter() //路由初始化
	err := r.Run(conf.Port)
	if err != nil {
		panic(err)
	}
}
