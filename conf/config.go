package conf

import (
	"strconv"
	"time"

	"gopkg.in/ini.v1"
)

type NodeConfig struct {
	Host     string `json:"host" forms:"host"`
	Port     int    `json:"port" forms:"port"`
	User     string `json:"user" forms:"user"`
	Password string `json:"password" forms:"password"`
}

// 解析配置文件
var (
	AppMode string //服务器启动模式，默认debug模式
	Port    string //服务启动端口

	JwtKey           string //JWT签名
	DbHost           string //数据库服务器主机
	DbPort           int    //数据库服务器端口
	DbUser           string //数据库用户名
	DbPassword       string //数据库密码
	DbName           string //数据库名
	TokenExpiredTime int64  //JWT过期时间

	KubeConfig  string    // kube.config文件位置
	SparkImage  string    // spark镜像
	LinuxImage  [2]string // linux镜像 1-centos 2-ubuntu
	SshPwd      string    // ssh默认密码
	HadoopImage string    // hadoop镜像

	RedisHost string // redis服务器
	RedisPort int

	MasterInfo NodeConfig // master的ssh信息

	Cpu             string        //Cpu配额
	Memory          string        //内存配额
	Storage         string        //存储配额
	Pvcstorage      string        //持久化存储配额
	Gpu             string        //Gpu配额
	UserExpiredTime time.Duration // 用户默认使用时间

	EmailAddr     string
	EmailUser     string
	EmailPassword string
	EmailHost     string
)

func init() {
	f, err := ini.Load("./conf/config.ini")
	if err != nil {
		panic("配置文件初始化失败")
	}
	loadEmail(f)
	loadServer(f)
	loadKubernetes(f)
	loadDb(f)
	loadJWt(f)
	loadPwd(f)
	loadMaster(f)
	loadResource(f)
}

// 加载服务器配置
func loadResource(file *ini.File) {
	s := file.Section("resource")
	Cpu = s.Key("Cpu").MustString("1")
	Memory = s.Key("Memory").MustString("1Gi")
	Storage = s.Key("Storage").MustString("10Gi")
	Pvcstorage = s.Key("Pvcstorage").MustString("10Gi")
	Gpu = s.Key("Gpu").MustString("1Gi")
	UserExpiredTime = time.Duration(s.Key("UserExpiredTime").MustInt(30) * 24 * int(time.Hour.Nanoseconds()))
}
func loadServer(file *ini.File) {
	s := file.Section("server")
	AppMode = s.Key("AppMode").MustString("debug")
	Port = s.Key("Port").MustString("8888")
	//PageSize = s.Key("PageSize").MustInt(10)
	//ProjectName = s.Key("ProjectName").MustString("")
}
func loadKubernetes(file *ini.File) {
	s := file.Section("kubernetes")
	KubeConfig = s.Key("KubeConfig").MustString("")
	SparkImage = s.Key("SparkImage").MustString("")
	HadoopImage = s.Key("HadoopImage").MustString("")
	LinuxImage[0] = s.Key("CentosImage").MustString("")
	LinuxImage[1] = s.Key("UbuntuImage").MustString("")
	//DockerHost = s.Key("DockerHost").MustString("")
}

// 加载数据库相关配置
func loadDb(file *ini.File) {
	s := file.Section("database")
	//DbType = s.Key("DbType").MustString("mysql")
	DbHost = s.Key("DbHost").MustString("")
	DbPort = s.Key("DbPort").MustInt(3306)
	DbUser = s.Key("DbUser").MustString("root")
	DbPassword = s.Key("DbPassWord").MustString("")
	DbName = s.Key("DbName").MustString("")
	RedisHost = s.Key("RedisHost").MustString("")
	RedisPort = s.Key("RedisPort").MustInt(6379)
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
	//BcryptCost, _ = strconv.Atoi(s.Key("bcryptCost").MustString("10"))
	SshPwd = s.Key("SshPwd").MustString("root123456")
}

// 加载master的ssh配置
func loadMaster(file *ini.File) {
	s := file.Section("master")
	MasterInfo.Host = s.Key("Host").MustString("")
	MasterInfo.Port = s.Key("Port").MustInt(22)
	MasterInfo.User = s.Key("User").MustString("root")
	MasterInfo.Password = s.Key("Password").MustString("")
}

func loadEmail(file *ini.File) {
	s := file.Section("email")
	EmailAddr = s.Key("EmailAddr").MustString("")
	EmailUser = s.Key("EmailUser").MustString("")
	EmailPassword = s.Key("EmailPassword").MustString("")
	EmailHost = s.Key("EmailHost").MustString("")
}
