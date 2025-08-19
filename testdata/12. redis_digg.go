package main

import (
	"Blog_server/core"
	"Blog_server/flags"
	"Blog_server/global"
	"Blog_server/service/redis_service/redis_count"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	global.Redis = core.InitRedis()
	core.InitLogrus()
	redis_count.NewDigg().Set("weJTo5gBDXSP1jc1K3NU")

}
