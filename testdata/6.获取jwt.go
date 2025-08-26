package main

import (
	"Blog_server/core"
	"Blog_server/flags"
	"Blog_server/global"
	"Blog_server/models/enum"
	"Blog_server/utils/jwts"
	"fmt"
	"github.com/sirupsen/logrus"
)

func main() {

	flags.Parse() //绑定命令行参数

	global.Config = core.ReadConf()
	core.InitLogrus()
	token, err := jwts.GetToken(jwts.Claims{
		UserID: 4,
		Role:   enum.AdminRole,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)

	username, err := jwts.ParseToken("dafartaeqwe1")
	if err != nil {
		fmt.Println(err)
		logrus.Errorf("%s", err)
		return
	}
	fmt.Println(username)

}
