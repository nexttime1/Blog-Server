package big_model_service

import (
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/utils/jwts"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserScopeEnableResponse struct {
	Enable bool `json:"enable"` // 用户能不能领取
	Scope  int  `json:"scope"`  // 能领取多少积分
}

type UserScopeRequest struct {
	Status bool `json:"status"`
}

// AutoReplyUpdateRequest 自动恢复 添加和修改的 请求
type AutoReplyUpdateRequest struct {
	ID           uint   `json:"id"`
	Name         string `json:"name" binding:"required"`               // 规则名称
	Mode         int    `json:"mode" binding:"required,oneof=1 2 3 4"` // 匹配模式 1 精确匹配，2 模糊匹配，3 前缀匹配，4 正则匹配
	Rule         string `json:"rule" binding:"required"`               // 匹配规则
	ReplyContent string `json:"replyContent" binding:"required"`       // 回复内容
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

func AutoReplyUpdateService(cr AutoReplyUpdateRequest) (int, error) { // 1为添加  2 为修改   0 为错误
	if cr.ID == 0 {
		// 说明要田间
		var model models.AutoReplyModel
		err := global.DB.Take(&model, "name = ?", cr.Name).Error
		if err == nil {
			// 重复了
			return 0, errors.New("规则名称重复")
		}
		err = global.DB.Create(&models.AutoReplyModel{
			Name:         cr.Name,
			Mode:         cr.Mode,
			Rule:         cr.Rule,
			ReplyContent: cr.ReplyContent,
		}).Error

		if err != nil {
			logrus.Errorf("%#v", err)
			return 0, errors.New("添加失败")

		}
		return 1, nil
	}
	var model models.AutoReplyModel

	err := global.DB.Take(&model, "id = ?", cr.ID).Error
	if err != nil {
		return 0, errors.New("要修改的id 不存在")
	}

	var arm models.AutoReplyModel
	err = global.DB.Take(&arm, "name = ? and id <> ?", cr.Name, cr.ID).Error
	if err == nil {
		//说明 你改的name 重复了
		return 0, errors.New("修改的规则名称重复")
	}

	err = global.DB.Model(&model).Updates(map[string]any{
		"name":          cr.Name,
		"mode":          cr.Mode,
		"rule":          cr.Rule,
		"reply_content": cr.ReplyContent,
	}).Error

	if err != nil {
		logrus.Errorf("%#v", err)
		return 0, errors.New("修改失败")
	}
	return 2, nil
}
