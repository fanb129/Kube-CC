package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"k8s_deploy_gin/conf"
	"k8s_deploy_gin/models"
	"log"
	"time"
)

var mysqlDb *gorm.DB

func InitDB() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.DbUser,
		conf.DbPassword,
		conf.DbHost,
		conf.DbPort,
		conf.DbName)

	mysqlDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//DisableForeignKeyConstraintWhenMigrating: true, //关闭外键！！！
		//NamingStrategy: schema.NamingStrategy{
		//	SingularTable: false,       //默认在表的后面加s
		//	TablePrefix:   "", // 表名前缀
		//},
		//SkipDefaultTransaction: true, // 禁用默认事务
	})
	if err != nil {
		log.Println("数据库连接失败：", err)
		panic(err)
		return
	}
	err = mysqlDb.AutoMigrate(&models.User{}, &models.Spark{}) // 数据库自动迁移
	if err != nil {
		log.Println("数据库自动迁移失败，err:", err)
		panic(err)
		return
	}
	sqlDb, err := mysqlDb.DB()
	if err != nil {
		log.Fatal(err)
		return
	}
	sqlDb.SetMaxIdleConns(50)                   // 连接池中的最大闲置连接数
	sqlDb.SetMaxOpenConns(100)                  // 数据库的最大连接数量
	sqlDb.SetConnMaxLifetime(100 * time.Second) // 连接的最大可复用时间
	return
}
