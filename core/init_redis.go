package core

import (
	"Blog_server/global"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func InitRedis() *redis.Client {
	r := global.Config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	})
	err := client.Ping().Err()
	if err != nil {
		fmt.Println("redis 连接失败")
		logrus.Errorf("redis 连接失败  %s", err.Error())
	}
	logrus.Infof("redis 连接成功")
	return client
}
