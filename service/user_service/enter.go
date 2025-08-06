package user_service

import (
	"Blog_server/common"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"Blog_server/utils/desensitization"
	"Blog_server/utils/jwts"
	"Blog_server/utils/pwd"
	"errors"
)

type EmailLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserInfoRequest struct {
	common.PageInfo
	Username string `json:"username"`
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

func UserInfoService(UserInfo UserInfoRequest) ([]models.UserModel, int, error) {
	list, count, err := common.ListQuery(models.UserModel{
		Username: UserInfo.Username,
	}, common.Options{
		PageInfo: UserInfo.PageInfo,
	})
	if err != nil {
		return nil, 0, err
	}
	var UserModelList []models.UserModel
	for _, model := range list {
		if model.Role == enum.AdminRole {
			//隐藏
			model.Username = ""
		}
		model.Tel = desensitization.TelDesensitization(model.Tel)
		model.Email = desensitization.EmailDesensitization(model.Email)
		UserModelList = append(UserModelList, model)
	}

	return UserModelList, count, nil

}
