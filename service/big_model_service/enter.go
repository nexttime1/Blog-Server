package big_model_service

import (
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/utils/jwts"
	"errors"
	"gorm.io/gorm"
)

type UserScopeEnableResponse struct {
	Enable bool `json:"enable"` // 用户能不能领取
	Scope  int  `json:"scope"`  // 能领取多少积分
}

type UserScopeRequest struct {
	Status bool `json:"status"`
}

func UserScopeEnableService(claim *jwts.MyClaims) (UserScopeEnableResponse, error) {

	// 查这个用户，今天能不能领取这个积分
	var userScopeModel models.UserScopeModel
	err := global.DB.Take(&userScopeModel, "user_id = ? and to_days(created_at)=to_days(now())", claim.UserID).Error
	var response UserScopeEnableResponse
	if err == nil {
		// 查到了
		return response, errors.New("今日已领取积分了")
	}
	response.Enable = true
	response.Scope = global.Config.BigModel.SessionSetting.DayScope
	return response, nil
}

func UserScopeService(cr UserScopeRequest, claim *jwts.MyClaims) error {

	var userScopeModel models.UserScopeModel
	err := global.DB.Take(&userScopeModel, "user_id = ? and to_days(created_at)=to_days(now())", claim.UserID).Error
	if err == nil {
		// 查到了
		return errors.New("今日已领取积分了")
	}

	// 没查到  那就可以 给user + 积分
	var user models.UserModel
	err = global.DB.Take(&user, "id = ?", claim.UserID).Error
	if err != nil {
		return errors.New("用户不存在")
	}

	scope := global.Config.BigModel.SessionSetting.DayScope
	err = global.DB.Model(&user).Update("scope", gorm.Expr("scope + ?", scope)).Error
	if err != nil {
		return errors.New("添加用户积分失败")
	}

	// 给用户加积分
	// 加数据
	global.DB.Create(&models.UserScopeModel{
		UserID: claim.UserID,
		Scope:  scope,
		Status: cr.Status,
	})

	return nil

}
