package main

import (
	"Blog_server/core"
	"Blog_server/flags"
	"Blog_server/global"
	"Blog_server/service/log_service"
)

func main() {
	flags.Parse() //绑定命令行参数
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.DB = core.InitDB()

	log := log_service.NewRuntimeLog("运行日志1111", log_service.RuntimeDateTypeHour)
	log.SetItem("xtm", 999)
	log.Save()

	log.SetItem("ttt", "love")
	log.Save()
}
