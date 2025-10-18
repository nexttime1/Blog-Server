package user_service

import (
	"Blog_server/common"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"Blog_server/utils/desensitization"
	"Blog_server/utils/jwts"
	"Blog_server/utils/pwd"
	"Blog_server/utils/struct_to_map"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

type EmailLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserInfoRequest struct {
	common.PageInfo
	Username string `json:"user_name"`
}

type UserCreateRequest struct {
	NickName string        `json:"nickname" binding:"required" msg:"请输入昵称"`  // 昵称
	UserName string        `json:"username" binding:"required" msg:"请输入用户名"` // 用户名
	Password string        `json:"password" binding:"required" msg:"请输入密码"`  // 密码
	Role     enum.RoleType `json:"role" binding:"required" msg:"请选择权限"`      // 权限  1 管理员  2 普通用户  3 游客
}

func UserEmailLoginService(mr EmailLoginRequest) (userModel models.UserModel, token string, msg string, err error) {
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
		Likes:    []string{"username", "nickname"},
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

const Avatar = "uploads/avatars/default.png"

func UserCreateService(UserName, NickName, Password string, Role enum.RoleType, Email string, Ip string, address string) error {
	var model models.UserModel
	err := global.DB.Where("username = ?", UserName).Take(&model).Error
	if err == nil {
		//找到了  重复
		logrus.Errorf("用户名已存在")
		return errors.New("用户名已存在")
	}
	//对密码 进行哈希
	hashPwd := pwd.HashPwd(Password)

	//入库
	err = global.DB.Create(&models.UserModel{
		Nickname:       NickName,
		Username:       UserName,
		Password:       hashPwd,
		Email:          Email,
		Role:           Role,
		Avatar:         Avatar,
		IP:             Ip,
		Addr:           address,
		RegisterSource: enum.SignEmail,
	}).Error
	if err != nil {
		logrus.Errorf("创建用户失败")
		return fmt.Errorf("创建用户失败  %s", err.Error())
	}
	logrus.Infof("创建%s用户成功", UserName)
	return nil
}

type UserUpdateInfoRequest struct {
	NickName string `json:"nickname" structs:"nickname"`
	Sign     string `json:"sign" structs:"sign"`
	Link     string `json:"link" structs:"link"`
	Avatar   string `json:"avatar" structs:"avatar"`
}

func UserInfoPutService(cr UserUpdateInfoRequest, claim *jwts.MyClaims) error {
	toMap := struct_to_map.StructToMap(cr)
	var User models.UserModel
	err := global.DB.Where("id = ?", claim.UserID).Take(&User).Error

	if err != nil {
		logrus.Errorf("用户不存在")
		return errors.New(fmt.Sprintf("用户不存在  %s", err.Error()))
	}

	_, ok := toMap["avatar"]
	if ok && User.RegisterSource != enum.SignEmail {
		delete(toMap, "avatar")
	}
	err = global.DB.Model(&User).Updates(toMap).Error
	if err != nil {
		logrus.Errorf("用户修改错误 %s", err.Error())
		return errors.New("修改用户信息失败")
	}

	return nil

}
