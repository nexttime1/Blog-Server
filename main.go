package main

import (
	"Blog_server/core"
	"Blog_server/flags"
	"Blog_server/global"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()

}
