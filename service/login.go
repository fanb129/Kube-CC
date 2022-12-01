package service

import (
	"errors"
	"gorm.io/gorm"
	"time"

	"Kube-CC/common"
	"Kube-CC/dao"
	"Kube-CC/middleware"

	"golang.org/x/crypto/bcrypt"
)

func Login(username, password string) (*common.LoginResponse, error) {
	user, err := dao.GetUserByName(username)
	if err != nil {
		return nil, errors.New("获取用户失败")
	}
	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("密码错误")
	}
	// 设置token
	token, err := middleware.SetUpToken(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("token生成失败")
	}
	return &common.LoginResponse{Response: common.OK, UserID: user.ID, Token: token}, nil
}

func Register(username, password, nickname string) (*common.Response, error) {
	user, _ := dao.GetUserByName(username)
	if user != nil {
		return nil, errors.New("账号已注册")
	}

	// 密码加密
	pwd, err := EncryptionPWD(password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 是否被软删除
	user, err = dao.GetDeletedUserByName(username)
	if user != nil {
		user.CreatedAt = time.Now()
		user.DeletedAt = gorm.DeletedAt{}
		user.Password = pwd
		user.Nickname = nickname
		user.Role = 1
		user.Avatar = ""
		row, err := dao.UpdateUnscopedUser(user)
		if err != nil || row == 0 {
			return nil, errors.New("注册失败")
		}
	} else {
		// 当前用户名不存在时
		row, err := dao.CreateUser(username, nickname, pwd)
		if err != nil || row == 0 {
			return nil, errors.New("注册失败")
		}
	}
	return &common.OK, nil
}

// EncryptionPWD 对密码进行加密
func EncryptionPWD(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("密码加密失败")
	}
	return string(hash), nil
}
