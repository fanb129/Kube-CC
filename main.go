package main

import (
	"fmt"
	"k8s_deploy_gin/pkg/setting"
	"k8s_deploy_gin/routers"
	"net/http"
)

func main() {
	r := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        r,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}

	//r.Run(fmt.Sprintf(":%d", setting.HTTPPort))

}
