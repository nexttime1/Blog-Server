package redis_jwt

import (
	"Blog_server/global"
	"Blog_server/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

type BlackType int8

const (
	UserBlackType   BlackType = 1 // 用户主动注销登录
	AdminBlackType  BlackType = 2 //管理员 强制注销登录
	DeviceBlackType BlackType = 3 // 不能几个设备登录
)

func ParseBlack(b string) BlackType {
	switch b {
	case "1":
		return UserBlackType
	case "2":
		return AdminBlackType
	case "3":
		return DeviceBlackType
	}
	return UserBlackType
}

func (b BlackType) Msg() string {
	switch b {
	case UserBlackType:
		return "已注销"
	case AdminBlackType:
		return "禁止登录"
	case DeviceBlackType:
		return "设备已下线"
	}
	return "已注销"
}

func (b BlackType) String() string {
	return fmt.Sprintf("%d", b)
}

func TokenBlack(token string, value BlackType) error {
	key := fmt.Sprintf("xtm_token_%s", token)
	Chains, err := jwts.ParseToken(token)
	if err != nil {
		logrus.Errorf("token解析失败 %s", err.Error())
		return err
	}
	second := Chains.ExpiresAt - time.Now().Unix()

	_, err = global.Redis.Set(key, value.String(), time.Duration(second)*time.Second).Result()
	if err != nil {
		logrus.Errorf("redis黑名单加载失败 %s", err.Error())
		return err
	}
	return nil
}

func HasTokenBlack(token string) (bool, BlackType) {
	key := fmt.Sprintf("xtm_token_%s", token)
	result, err := global.Redis.Get(key).Result()
	if err != nil {
		return false, UserBlackType
	}
	blackType := ParseBlack(result)
	return true, blackType
}

func HasTokenBlackByGin(c *gin.Context) (bool, BlackType) {
	token := c.GetHeader("Token")
	if token == "" {
		token = c.Query("token")

	}
	return HasTokenBlack(token)
}

func TokenBlackByGin(c *gin.Context, value BlackType) error {
	token := c.GetHeader("Token")
	if token == "" {
		token = c.Query("token")

	}
	return TokenBlack(token, value)
}
