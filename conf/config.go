package conf

import (
	"gopkg.in/ini.v1"
	"log"
	"strconv"
)

// 解析配置文件
var (
	AppMode          string //服务器启动模式，默认debug模式
	Port             string //服务启动端口
	JwtKey           string //JWT签名
	DbType           string //数据库类型
	DbHost           string //数据库服务器主机
	DbPort           string //数据库服务器端口
	DbUser           string //数据库用户名
	DbPassword       string //数据库密码
	BcryptCost       int    //bcrypt 生成密码时的cost
	DbName           string //数据库名
	TokenExpiredTime int64  //JWT过期时间
	PageSize         int    // 分页大小
	KubeConfig       string // kube.config文件位置
)

func init() {
	f, err := ini.Load("./conf/config.ini")
	if err != nil {
		log.Fatal("配置文件初始化失败")
	}

	loadServer(f)
	loadKubernetes(f)
	loadDb(f)
	loadJWt(f)

	BcryptCost, err = strconv.Atoi(f.Section("password").Key("bcryptCost").MustString("10"))
	if err != nil {
		log.Fatal("BcryptCost加载失败")
	}
}

// 加载服务器配置
func loadServer(file *ini.File) {
	s := file.Section("server")
	AppMode = s.Key("AppMode").MustString("debug")
	Port = s.Key("Port").MustString("8888")
	PageSize = s.Key("PageSize").MustInt(10)
}
func loadKubernetes(file *ini.File) {
	s := file.Section("kubernetes")
	KubeConfig = s.Key("KubeConfig").MustString("")
}

// 加载数据库相关配置
func loadDb(file *ini.File) {
	s := file.Section("database")
	DbType = s.Key("DbType").MustString("mysql")
	DbHost = s.Key("DbHost").MustString("39.103.195.185")
	DbPort = s.Key("DbPort").MustString("3306")
	DbUser = s.Key("DbUser").MustString("root")
	DbPassword = s.Key("DbPassWord").MustString("Fb123456.")
	DbName = s.Key("DbName").MustString("k8s_deploy_gin")
}

// 加载JWT相关配置
func loadJWt(file *ini.File) {
	s := file.Section("jwt")
	JwtKey = s.Key("JwtKey").MustString("")
	TokenExpiredTime, _ = strconv.ParseInt(s.Key("TokenExpiredTime").MustString("1000"), 10, 64)
}
