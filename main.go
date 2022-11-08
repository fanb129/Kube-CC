package main

import (
	"go.uber.org/zap"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/conf"
	"k8s_deploy_gin/dao"
	"k8s_deploy_gin/log"
	"k8s_deploy_gin/routers"

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
	r := routers.InitRouter() //路由初始化
	// 初始化翻译
	if err := common.InitTrans("zh"); err != nil {
		zap.S().Panicln(err)
	}
	if err := r.Run(conf.Port); err != nil {
		zap.S().Panicln(err)
	}

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

}
