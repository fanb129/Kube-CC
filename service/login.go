package service

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/models"
	"errors"
	"time"

	"gorm.io/gorm"

	"Kube-CC/dao"
	"Kube-CC/middleware"

	"golang.org/x/crypto/bcrypt"
)

// <<修改>>
func Login(usernameoremail, password string) (*responses.LoginResponse, error) {
	//分别通过用户名和邮箱查找用户
	var (
		user *models.User
		err  error
	)
	if dao.VerifyEmailFormat(usernameoremail) {
		user, err = dao.GetUserByEmail(usernameoremail)
	} else {
		user, err = dao.GetUserByName(usernameoremail)
	}
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
	return &responses.LoginResponse{Response: responses.OK, UserID: user.ID, Token: token}, nil
}

// <<修改>>
//func Register(username, password, nickname, email string) (*responses.Response, error) {
//	user, _ := dao.GetUserByName(username)
//	if user != nil {
//		return nil, errors.New("账号已注册")
//	}
//	user, _ = dao.GetUserByEmail(email)
//	if user != nil {
//		return nil, errors.New("邮箱已使用")
//	}
//
//	// 密码加密
//	pwd, err := EncryptionPWD(password)
//	if err != nil {
//		return nil, errors.New("密码加密失败")
//	}
//
//	// 是否被软删除
//	user, err = dao.GetDeletedUserByName(username)
//	if user != nil {
//		user.CreatedAt = time.Now()
//		user.DeletedAt = gorm.DeletedAt{}
//		user.Password = pwd
//		user.Nickname = nickname
//		user.Role = 1
//		user.Avatar = ""
//		row, err := dao.UpdateUnscopedUser(user)
//		if err != nil || row == 0 {
//			return nil, errors.New("注册失败")
//		}
//	} else {
//		// 当前用户名不存在时
//		row, err := dao.CreateUser(username, nickname, pwd, email)
//		if err != nil || row == 0 {
//			return nil, errors.New("注册失败")
//		}
//	}
//	return &responses.OK, nil
//}

// EncryptionPWD 对密码进行加密
func EncryptionPWD(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("密码加密失败")
	}
	return string(hash), nil
}

// RegisterUser 添加用户
func RegisterUser(username, password, nickname string, gid uint) (*responses.Response, error) {
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
		user.Email = ""
		user.Groupid = gid
		row, err := dao.UpdateUnscopedUser(user)
		if err != nil || row == 0 {
			return nil, errors.New("注册失败")
		}
	} else {
		// 当前用户名不存在时
		row, err := dao.CreateUser(username, nickname, pwd, gid)
		if err != nil || row == 0 {
			return nil, errors.New("注册失败")
		}
	}
	return &responses.OK, nil
}

// FindPass 找回密码
func FindPass(form forms.FindPass) (*responses.Response, error) {
	err := validateEmailCode(form.Email, form.VCode)
	if err != nil {
		return nil, err
	}
	user, err := dao.GetUserByEmail(form.Email)
	if err != nil {
		return nil, err
	}
	// 密码加密
	pwd, err := EncryptionPWD(form.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}
	user.Password = pwd
	_, err = dao.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return &responses.OK, nil
}
