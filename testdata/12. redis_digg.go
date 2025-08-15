package main

import (
	"Blog_server/core"
	"Blog_server/flags"
	"Blog_server/global"
	"Blog_server/service/redis_service/redis_digg"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	global.Redis = core.InitRedis()
	core.InitLogrus()
	redis_digg.Digging("weJTo5gBDXSP1jc1K3NU")

}
