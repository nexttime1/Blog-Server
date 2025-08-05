package user_service

import (
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/utils/jwts"
	"Blog_server/utils/pwd"
	"errors"
)

type EmailLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func UserEmailLoginService(mr EmailLoginRequest) (token string, msg string, err error) {
	var userModel models.UserModel
	msg = "登录成功"
	token = ""
	err = global.DB.Where("username = ?", mr.Username).Take(&userModel).Error
	if err != nil {
		msg = "用户名或密码错误"
		return
	}
	//验证密码
	ok := pwd.CheckPwd(userModel.Password, mr.Password)
	if !ok {
		msg = "用户名或密码错误"
		err = errors.New("密码不正确")
		return
	}
	// 获得 token
	token, err = jwts.GetToken(jwts.Claims{
		UserID:   userModel.ID,
		Username: userModel.Username,
		Role:     userModel.Role,
	})
	if err != nil {
		msg = "token 获取失败"
		return
	}

	return
}
