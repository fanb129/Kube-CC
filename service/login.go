package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/dao"
	"k8s_deploy_gin/middleware"
	"time"
)

func Login(username, password string) (*common.LoginResponse, error) {
	user, err := dao.GetUserByName(username)
	if err != nil {
		return nil, err
	}
	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	// 设置token
	token, err := middleware.SetUpToken(user.ID, user.Status)
	if err != nil {
		return nil, errors.New("token生成失败")
	}
	return &common.LoginResponse{Response: common.OK, UserID: user.ID, Token: token}, nil
}

func Register(username, nickname, password string) (*common.RegisterResponse, error) {
	// 密码加密
	pwd, err := EncryptionPWD(password)
	if err != nil {
		return nil, err
	}
	user, _ := dao.GetUserByName(username)
	if user != nil {
		// 如果当前用户已注册
		if user.Status > 0 {
			return nil, errors.New("username exists")
		}
		// 注册过，但已删除
		user.Nickname = nickname
		user.Avatar = ""
		user.Status = 1
		user.Major = ""
		user.Class = ""
		user.CreateAt = time.Now()
		user.Password = pwd
		err = dao.UpdateUser(user)
		if err != nil {
			return nil, err
		}
		return &common.RegisterResponse{Response: common.OK, UserID: user.ID}, nil
	}
	// 当前用户名不存在时
	uid, err := dao.CreateUser(username, nickname, pwd)
	if err != nil {
		return nil, err
	}
	return &common.RegisterResponse{Response: common.OK, UserID: uid}, nil
}

// EncryptionPWD 对密码进行加密
func EncryptionPWD(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("密码加密失败")
	}
	return string(hash), nil
}
