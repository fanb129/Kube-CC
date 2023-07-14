package main

import (
	"Kube-CC/common"
	"Kube-CC/conf"
	"Kube-CC/dao"
	"Kube-CC/log"
	"Kube-CC/routers"
	"Kube-CC/service"
	"time"

	"go.uber.org/zap"

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
	// if err := docker.ConnectDocker(); err != nil {
	// 	zap.S().Panicln(err)
	// }

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

	go func() {
		// 每隔一小时检测用户是否过期
		ticker := time.NewTicker(time.Hour)
		for {
			select {
			case <-ticker.C:
				zap.S().Infoln("开始扫描过期时间")
				startTime := time.Now()
				users, err := dao.ListAllUser()
				if err != nil {
					zap.S().Errorln("获取all user失败:", err)
				}
				for _, user := range users {
					// 如果过期时间在现在之前，则删除
					if user.ExpiredTime.Before(time.Now()) {
						// 删除user
						_, err := service.DeleteUSer(user.ID)
						if err != nil {
							zap.S().Errorln("删除user失败:", err)
						} else {
							zap.S().Infoln("delete user:", user.Username)
						}
					}
				}
				endTime := time.Now()
				zap.S().Infoln("扫描结束,用时", endTime.Sub(startTime).String())
			}
		}
	}()

	if err := r.Run(conf.Port); err != nil {
		zap.S().Panicln(err)
	}
	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
