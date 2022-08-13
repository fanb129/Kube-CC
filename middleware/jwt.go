package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/conf"
	"net/http"
	"time"
)

type JWT struct {
	JwtKey []byte
}

// NewJWT 创建JWT实列
func NewJWT() *JWT {
	return &JWT{JwtKey: []byte(conf.JwtKey)}
}

// MyClaims 自定义Claim
type MyClaims struct {
	UserId uint `json:"user_id"` // 用户id
	Role   uint `json:"role"`    // 用户权限
	jwt.StandardClaims
}

var (
	TokenExpired     = errors.New("token已过期，请重新登录")
	TokenNotValidYet = errors.New("token无效，请重新登录")
	TokenMalFormed   = errors.New("token 不正确，请重新登录")
	TokenInvalid     = errors.New("这不是一个 token,请重新登录")
)

// SetUpToken 设置claims,为生成token准备
func SetUpToken(userID uint, role uint) (string, error) {
	j := NewJWT()
	claims := MyClaims{
		UserId: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 240,
			ExpiresAt: time.Now().Unix() + conf.TokenExpiredTime,
		},
	}
	token, err := j.CreatToken(claims)
	if err != nil {
		return "", errors.New("token生成失败")
	}
	return token, nil
}

// CreatToken 通过加密和claims生成token
func (j *JWT) CreatToken(claims MyClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JwtKey)
}

// ParseToken 解析token，返回claims
func (j *JWT) ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.JwtKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalFormed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	if token != nil {
		if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenExpired
	}

	return nil, TokenInvalid
}

// 提取token，若不存在返回空
func getToken(c *gin.Context) (string, bool) {
	var token string
	var ok bool

	token = c.Request.Header.Get("token")
	if token == "" {
		token, ok = c.GetQuery("token")
		// 如果Query 参数提取到 token 直接返回
		if ok {
			return token, true
		}
		// 否则继续从Form 参数里面提取
		token, ok = c.GetPostForm("token")
		if !ok {
			return "", false
		}
	}

	return token, true
}

// JWTToken 解析、验证token，并把解析出来的user_id 通过ctx.Set() 方法增加到 gin.Context 头部中
func JWTToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		var ok bool
		// 提取token，若不存在
		if token, ok = getToken(c); !ok {
			c.JSON(http.StatusUnauthorized, common.NoToken)
			c.Abort()
			return
		}

		j := NewJWT()
		// token存在，解析token
		claims, err := j.ParseToken(token)
		if err != nil {
			// token过期
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, common.TokenExpired)
				c.Abort()
				return
			}
			// 其他错误(不是一个token）
			c.JSON(http.StatusUnauthorized, common.NoToken)
			c.Abort()
			return
		}
		c.Set("u_id", claims.UserId) // 把解析出来的userID放进头部
		c.Set("role", claims.Role)

		c.Next()
	}
}
