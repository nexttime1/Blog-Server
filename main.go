package main

import (
	"Blog_server/core"
	"Blog_server/flags"
)

func main() {
	flags.Parse()
	core.ReadConf()

}
