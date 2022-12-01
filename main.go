package main

import (
	"Kube-CC/common"
	"Kube-CC/conf"
	"Kube-CC/dao"
	"Kube-CC/log"
	"Kube-CC/routers"
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
