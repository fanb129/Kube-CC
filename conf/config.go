package conf

import (
	"gopkg.in/ini.v1"
	"log"
	"strconv"
)

// 解析配置文件
var (
	AppMode          string    //服务器启动模式，默认debug模式
	Port             string    //服务启动端口
	JwtKey           string    //JWT签名
	DbType           string    //数据库类型
	DbHost           string    //数据库服务器主机
	DbPort           string    //数据库服务器端口
	DbUser           string    //数据库用户名
	DbPassword       string    //数据库密码
	BcryptCost       int       //bcrypt 生成密码时的cost
	DbName           string    //数据库名
	TokenExpiredTime int64     //JWT过期时间
	PageSize         int       // 分页大小
	KubeConfig       string    // kube.config文件位置
	SparkImage       string    // spark镜像
	LinuxImage       [2]string // linux镜像 1-centos 2-ubuntu
	ProjectName      string    // 项目名称，用于ingress域名
	SshPwd           string    // ssh默认密码
	HadoopImage      string    // hadoop镜像
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
	loadPwd(f)
}

// 加载服务器配置
func loadServer(file *ini.File) {
	s := file.Section("server")
	AppMode = s.Key("AppMode").MustString("debug")
	Port = s.Key("Port").MustString("8888")
	PageSize = s.Key("PageSize").MustInt(10)
	ProjectName = s.Key("ProjectName").MustString("")
}
func loadKubernetes(file *ini.File) {
	s := file.Section("kubernetes")
	KubeConfig = s.Key("KubeConfig").MustString("")
	SparkImage = s.Key("SparkImage").MustString("")
	HadoopImage = s.Key("HadoopImage").MustString("")
	LinuxImage[0] = s.Key("CentosImage").MustString("")
	LinuxImage[1] = s.Key("UbuntuImage").MustString("")
}

// 加载数据库相关配置
func loadDb(file *ini.File) {
	s := file.Section("database")
	DbType = s.Key("DbType").MustString("mysql")
	DbHost = s.Key("DbHost").MustString("")
	DbPort = s.Key("DbPort").MustString("3306")
	DbUser = s.Key("DbUser").MustString("root")
	DbPassword = s.Key("DbPassWord").MustString("")
	DbName = s.Key("DbName").MustString("k8s_deploy_gin")
}

// 加载JWT相关配置
func loadJWt(file *ini.File) {
	s := file.Section("jwt")
	JwtKey = s.Key("JwtKey").MustString("")
	TokenExpiredTime, _ = strconv.ParseInt(s.Key("TokenExpiredTime").MustString("1000"), 10, 64)
}

// 加载密码相关配置
func loadPwd(file *ini.File) {
	s := file.Section("password")
	BcryptCost, _ = strconv.Atoi(s.Key("bcryptCost").MustString("10"))
	SshPwd = s.Key("SshPwd").MustString("root123456")
}
