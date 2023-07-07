package main

import (
	"Kube-CC/common"
	"Kube-CC/conf"
	"Kube-CC/dao"
	"Kube-CC/log"
	"Kube-CC/routers"
	"Kube-CC/service"
	"Kube-CC/service/docker"
	"go.uber.org/zap"
	"time"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 初始化日志库
	log.InitLogger()

	//数据库初始化
	if err := dao.InitDB(); err != nil {
		zap.S().Panicln(err)
	}
	dao.InitRedis() //Redis 初始化(暂时不用)
	// client-go k8s初始化
	if err := dao.InitKube(); err != nil {
		zap.S().Panicln(err)
	}

	// 初始化连接docker
	if err := docker.ConnectDocker(); err != nil {
		zap.S().Panicln(err)
	}

	// TODO login提交完毕后对功能进行本地测试
	//docker.CreateImage()
	//docker.GetImageList()
	//docker.AddNewTag()

	//docker.DeleteImage()

	//TODO 待测

	//docker.PullImage()
	//docker.SaveImage()

	r := routers.InitRouter() //路由初始化
	// 初始化翻译
	if err := common.InitTrans("zh"); err != nil {
		zap.S().Panicln(err)
	}
	if err := r.Run(conf.Port); err != nil {
		zap.S().Panicln(err)
	}

	go func() {
		// 每隔一小时检测ttl
		ticker := time.NewTicker(time.Hour)
		for {
			select {
			case <-ticker.C:
				ttls, err := dao.ListTtl()
				if err != nil {
					zap.S().Errorln("获取ttl失败:", err)
				}
				for _, ttl := range ttls {
					// 如果过期时间在现在之前，则删除
					if ttl.ExpiredTime.Before(time.Now()) {
						// 删除ns
						_, err := service.DeleteNs(ttl.Namespace)
						if err != nil {
							zap.S().Errorln("删除ns失败:", err)
						}
						err = service.DeleteTtl(ttl.Namespace)
						if err != nil {
							zap.S().Errorln("删除ttl失败:", err)
						}
					}
				}
			}
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

}
