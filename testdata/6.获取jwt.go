package main

import (
	"Blog_server/core"
	"Blog_server/flags"
	"Blog_server/global"
	"Blog_server/utils/jwts"
	"fmt"
)

func main() {

	flags.Parse() //绑定命令行参数
	global.Config = core.ReadConf()
	core.InitLogrus()
	token, err := jwts.GetToken(jwts.Claims{
		UserID: 2,
		Role:   2,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)

	username, err := jwts.ParseToken("dafartaeqwe1")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(username)

}
