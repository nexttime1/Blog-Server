package main

import (
	"Blog_server/core"
	_ "Blog_server/docs"
	"Blog_server/flags"
	"Blog_server/global"
	"Blog_server/router"
)

// @title Blog_Server API文档
// @version 1.0
// @description Blog_Server API文档
// @host 127.0.0.01:8080
// @BasePath /
func main() {
	flags.Parse() //绑定命令行参数
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.DB = core.InitDB()
	global.Redis = core.InitRedis()
	flags.Run()
	router.Run()
}
