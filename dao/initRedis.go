package dao

import (
	"Kube-CC/conf"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var RedisClient *redis.Client

func InitRedis() {
	//redis的配置
	redisOption := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.RedisHost, conf.RedisPort),
		DialTimeout:  time.Millisecond * 100,
		ReadTimeout:  time.Millisecond * 100,
		WriteTimeout: time.Millisecond * 200,
		PoolSize:     20,
		MinIdleConns: 3,
		MaxConnAge:   50,
	}
	RedisClient = redis.NewClient(redisOption)
}
