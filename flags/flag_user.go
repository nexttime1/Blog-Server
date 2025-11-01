package flags

import (
	"Blog_server/models/enum"
	"Blog_server/service/user_service"
	"fmt"
	"github.com/sirupsen/logrus"
)

func FlagUser(permission string) {
	var (
		UserName   string
		NickName   string
		Password   string
		RePassword string
		Email      string
	)
	fmt.Print("请输入用户名 （必填）")
	fmt.Scan(&UserName)
	fmt.Print("请输入昵称 （必填）")
	fmt.Scan(&NickName)
	fmt.Print("请输入密码 （必填）")
	fmt.Scan(&Password)
	fmt.Print("请再次输入密码 （必填）")
	fmt.Scan(&RePassword)
	fmt.Print("请输入邮箱 （必填）")
	fmt.Scan(&Email)

	// 判断 密码
	if Password != RePassword {
		logrus.Errorf("两次密码不一致，请重新输入")
		return
	}
	role := enum.UserRole
	if permission == "admin" {
		role = enum.AdminRole
	}

	err := user_service.UserCreateService(UserName, NickName, Password, role, Email, "内网", "内网ip")
	if err != nil {
		return
	}

}
