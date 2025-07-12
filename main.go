package main

import (
	"Blog_server/core"
	"Blog_server/flags"
	"Blog_server/global"
)

func main() {
	flags.Parse() //绑定命令行参数
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.DB = core.InitDB()

	flags.Run()
}
