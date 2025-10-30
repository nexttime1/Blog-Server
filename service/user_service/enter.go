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

type UserListResponse struct {
	models.Model
	Nickname       string            `json:"nick_name"`
	Username       string            `json:"user_name"`
	Avatar         string            `json:"avatar"`
	Email          string            `json:"email"`
	Tel            string            `json:"tel"`
	Addr           string            `json:"addr"`
	Token          string            `json:"token"`
	IP             string            `json:"ip"`
	Role           enum.RoleType     `json:"role"`
	RegisterSource enum.RegisterType `json:"sign_status"`
	Scope          int               `json:"scope"`
	Sign           string            `json:"sign"`
	Link           string            `json:"link"`
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

func UserInfoService(UserInfo UserInfoRequest) ([]UserListResponse, int, error) {
	list, count, err := common.ListQuery(models.UserModel{
		Username: UserInfo.Username,
	}, common.Options{
		PageInfo: UserInfo.PageInfo,
		Likes:    []string{"username", "nickname"},
	})
	if err != nil {
		return nil, 0, err
	}
	var UserModelList []UserListResponse
	for _, model := range list {
		if model.Role == enum.AdminRole {
			//隐藏
			model.Username = ""
		}
		model.Tel = desensitization.TelDesensitization(model.Tel)
		model.Email = desensitization.EmailDesensitization(model.Email)
		UserModelList = append(UserModelList, UserListResponse{
			Model:          model.Model,
			Nickname:       model.Nickname,
			Username:       model.Username,
			Avatar:         model.Avatar,
			Email:          model.Email,
			Tel:            model.Tel,
			Addr:           model.Addr,
			Token:          model.Token,
			IP:             model.IP,
			Role:           model.Role,
			RegisterSource: model.RegisterSource,
			Scope:          model.Scope,
			Sign:           model.Sign,
			Link:           model.Link,
		})
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
	//方便看看改没改  旧值
	nick_name := User.Nickname
	avatar := User.Avatar
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

	// 修改成功后消息列表也要修改
	{
		// 准备新的值（注意类型转换为 string）
		var (
			newNick       string
			newAvatar     string
			nickChanged   bool
			avatarChanged bool
		)

		if v, ok := toMap["nickname"]; ok {
			newNick = fmt.Sprintf("%v", v)
			if newNick != nick_name {
				nickChanged = true
			}
		}

		if v, ok := toMap["avatar"]; ok {
			newAvatar = fmt.Sprintf("%v", v)
			if newAvatar != avatar {
				avatarChanged = true
			}
		}

		// 如果没有任何变化则跳过更新
		if !nickChanged && !avatarChanged {
			// nothing to update in messages
		} else {
			// 分别为发送方和接收方构建 updates map，只包含实际变化的字段
			sendUpdates := map[string]interface{}{}
			revUpdates := map[string]interface{}{}

			if nickChanged {
				sendUpdates["send_user_nick_name"] = newNick
				revUpdates["rev_user_nick_name"] = newNick
			}
			if avatarChanged {
				sendUpdates["send_user_avatar"] = newAvatar
				revUpdates["rev_user_avatar"] = newAvatar
			}

			// 用事务包裹（可选，但推荐）以保证一致性
			tx := global.DB.Begin()
			if len(sendUpdates) > 0 {
				if err := tx.Model(&models.MessageModel{}).
					Where("send_user_id = ?", claim.UserID).
					Updates(sendUpdates).Error; err != nil {
					tx.Rollback()
					logrus.Errorf("更新消息发送方信息失败: %s", err.Error())
					// 不返回错误，仅记录（保持原函数行为），或根据你的需求返回 err
				}
			}

			if len(revUpdates) > 0 {
				if err := tx.Model(&models.MessageModel{}).
					Where("rev_user_id = ?", claim.UserID).
					Updates(revUpdates).Error; err != nil {
					tx.Rollback()
					logrus.Errorf("更新消息接收方信息失败: %s", err.Error())
					// 不返回错误，仅记录（保持原函数行为），或根据你的需求返回 err
				}
			}

			// 尝试提交事务（如果没有 rollback）
			if err := tx.Commit().Error; err != nil {
				logrus.Errorf("提交消息更新事务失败: %s", err.Error())
			}
		}
	}

	return nil

}
