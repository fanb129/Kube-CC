package myjwt

import (
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/dao"
	"k8s_deploy_gin/models"
)

// IAuthorizator 授权规则接口
type IAuthorizator interface {
	HandleAuthorizator(data interface{}, c *gin.Context) bool
}

// FourAuthorizator 超级管理员授权规则 status==4
type FourAuthorizator struct{}

// HandleAuthorizator 处理超级管理员授权规则
func (*FourAuthorizator) HandleAuthorizator(data interface{}, c *gin.Context) bool {
	if u, ok := data.(*models.User); ok {
		user := dao.GetUserByName(u.Username)
		if user.Status == 4 {
			return true
		}
	}
	return false
}

// ThreeAuthorizator 集群管理员授权规则 status==3
type ThreeAuthorizator struct{}

// HandleAuthorizator 处理集群管理员授权规则
func (*ThreeAuthorizator) HandleAuthorizator(data interface{}, c *gin.Context) bool {
	if u, ok := data.(*models.User); ok {
		user := dao.GetUserByName(u.Username)
		if user.Status == 3 {
			return true
		}
	}
	return false
}

// TwoAuthorizator 集群管理员授权规则 status==2
type TwoAuthorizator struct{}

// HandleAuthorizator 处理集群管理员授权规则
func (*TwoAuthorizator) HandleAuthorizator(data interface{}, c *gin.Context) bool {
	if u, ok := data.(*models.User); ok {
		user := dao.GetUserByName(u.Username)
		if user.Status == 2 {
			return true
		}
	}
	return false
}

//AllUserAuthorizator 普通用户授权规则 status == 1
type AllUserAuthorizator struct {
}

//HandleAuthorizator 处理普通用户授权规则
func (*AllUserAuthorizator) HandleAuthorizator(data interface{}, c *gin.Context) bool {
	return true
}
