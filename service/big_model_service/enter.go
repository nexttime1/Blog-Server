package big_model_service

import (
	"Blog_server/common"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/utils/jwts"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"regexp"
)

// UserScopeEnableResponse 服务 UserScopeEnableService
type UserScopeEnableResponse struct {
	Enable bool `json:"enable"` // 用户能不能领取
	Scope  int  `json:"scope"`  // 能领取多少积分
}

// UserScopeRequest  服务  UserScopeService
type UserScopeRequest struct {
	Status bool `json:"status"`
}

// AutoReplyUpdateRequest 自动恢复 添加和修改的 请求 AutoReplyUpdateService
type AutoReplyUpdateRequest struct {
	ID           uint   `json:"id"`
	Name         string `json:"name" binding:"required"`               // 规则名称
	Mode         int    `json:"mode" binding:"required,oneof=1 2 3 4"` // 匹配模式 1 精确匹配，2 模糊匹配，3 前缀匹配，4 正则匹配
	Rule         string `json:"rule" binding:"required"`               // 匹配规则
	ReplyContent string `json:"replyContent" binding:"required"`       // 回复内容
}

// TagUpdateRequest 服务 BigModelTagUpdateService
type TagUpdateRequest struct {
	ID    uint   `json:"id"`                       // 更新使用
	Title string `json:"title" binding:"required"` // 名称
	Color string `json:"color" binding:"required"` // 颜色
}

// TagListResponse 服务于 BigModelTagListService
type TagListResponse struct {
	models.Model
	Title     string `json:"title"`     // 名称
	Color     string `json:"color"`     // 颜色
	RoleCount int    `json:"roleCount"` // 角色个数
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

	// 校验正则是否写错
	if cr.Mode == 4 {
		_, err := regexp.Compile(cr.Rule)
		if err != nil {
			logrus.Errorf(fmt.Sprintf("正则表达式错误 %s", err.Error()))
			return 0, errors.New("正则表达式错误")
		}
	}
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

func BigModelTagUpdateService(cr TagUpdateRequest) (int, error) {
	if cr.ID == 0 {
		//增加
		var model models.BigModelTagModel
		err := global.DB.Take(&model, "title = ?", cr.Title).Error
		if err == nil {
			// 找到了  重复了
			return 0, errors.New("与已有的标签名称重复")
		}
		global.DB.Create(&models.BigModelTagModel{
			Title: cr.Title,
			Color: cr.Color,
		})
		return 1, nil
	}
	// 修改
	var model models.BigModelTagModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		return 0, errors.New("标签不存在")
	}
	// 查一查 修改后的标签会不会重复
	var arm models.BigModelTagModel
	err = global.DB.Take(&arm, "title = ? and id <> ?", cr.Title, cr.ID).Error
	if err == nil {
		// 找到了
		return 0, errors.New("修改的名称重复")
	}
	err = global.DB.Model(&model).Updates(map[string]any{
		"title": cr.Title,
		"color": cr.Color,
	}).Error
	if err != nil {
		logrus.Errorf("%#v", err)
		return 0, errors.New("修改失败")
	}

	return 2, nil

}

func BigModelTagListService(cr common.PageInfo) ([]TagListResponse, int, error) {
	list, count, err := common.ListQuery(models.BigModelTagModel{}, common.Options{
		PageInfo: cr,
		Likes:    []string{"title"},
		Preload:  []string{"Roles"},
	})
	if err != nil {
		return []TagListResponse{}, 0, err
	}

	var response []TagListResponse
	for _, model := range list {
		response = append(response, TagListResponse{
			Model:     model.Model,
			Title:     model.Title,
			Color:     model.Color,
			RoleCount: len(model.Roles),
		})
	}
	return response, count, nil

}

func BigModelTagRemoveService(cr models.RemoveRequest) error {
	// 开启事务，保证操作原子性
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 删除中间表关联记录
	err := tx.Where("big_model_tag_model_id in ?", cr.IDList).Delete(&models.BigModelRoleTagModel{}).Error
	if err != nil {
		tx.Rollback()
		logrus.Errorf("中间表关联记录删除失败：%#v", err)
		return errors.New("关联表删除失败")
	}

	// 2. 检查标签是否存在
	var existCount int64
	err = tx.Model(&models.BigModelTagModel{}).Where("id in ?", cr.IDList).Count(&existCount).Error
	if err != nil {
		tx.Rollback()
		logrus.Errorf("标签存在性查询失败：%#v", err)
		return errors.New("查询标签失败")
	}
	// 判断 传入的 ID 数量与存在的数量不一致，说明有无效 ID
	if existCount != int64(len(cr.IDList)) {
		tx.Rollback()
		return errors.New("部分标签不存在，删除失败")
	}

	err = tx.Where("id in ?", cr.IDList).Delete(&models.BigModelTagModel{}).Error
	// 3. 批量删除标签表数据
	if err != nil {
		tx.Rollback()
		logrus.Errorf("标签删除失败：%#v", err)
		return errors.New("标签删除失败")
	}
	// 所有操作成功，提交事务
	return tx.Commit().Error
}
