package flags

import (
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"Blog_server/utils/pwd"
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
	var model models.UserModel
	err := global.DB.Where("username = ?", UserName).Take(&model).Error
	if err == nil {
		//找到了  重复
		logrus.Errorf("用户名已存在")
		return
	}
	role := enum.UserRole
	if permission == "admin" {
		role = enum.AdminRole
	}

	//对密码 进行哈希
	hashPwd := pwd.HashPwd(Password)

	avatar := "uploads/avatars/default.png"
	//入库
	err = global.DB.Create(&models.UserModel{
		Nickname:       NickName,
		Username:       UserName,
		Password:       hashPwd,
		Email:          Email,
		Role:           role,
		Avatar:         avatar,
		IP:             "127.0.0.1",
		Addr:           "内网地址",
		RegisterSource: enum.SignEmail,
	}).Error
	if err != nil {
		logrus.Errorf("创建用户失败")
		return
	}
	logrus.Infof("创建%s用户成功", UserName)
}
