package dao

import (
	"github.com/heroku/docker-registry-client/registry"
)

var Hub *registry.Registry

func InitRegistry() (err error) {
	url := "http://192.168.139.141:5000/"
	username := "admin"    // anonymous
	password := "passw0rd" // anonymous
	Hub, err = registry.New(url, username, password)
	// 取消默认log
	Hub.Logf = registry.Quiet
	return err
}
