package main

import (
	"Blog_server/core"
	"Blog_server/flags"
	"Blog_server/global"
	"Blog_server/models/enum"
	"Blog_server/service/redis_service/redis_jwt"
	"Blog_server/utils/jwts"
	"fmt"
)

func main() {
	flags.Parse()

	global.Config = core.ReadConf()
	global.DB = core.InitDB()
	core.InitLogrus()
	global.Redis = core.InitRedis()
	claims := jwts.Claims{
		UserID:   10,
		Username: "sss",
		Role:     enum.UserRole,
	}
	token, err := jwts.GetToken(claims)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println(token)

	//设置黑名单
	redis_jwt.TokenBlack(token, redis_jwt.UserBlackType)
	//查看是否在黑名单
	ok, blackType := redis_jwt.HasTokenBlack(token)
	fmt.Println(ok, blackType)

}
