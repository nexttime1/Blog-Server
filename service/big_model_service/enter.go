package big_model_service

import (
	"Blog_server/common"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/utils/jwts"
	"Blog_server/utils/struct_to_map"
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

// RoleUpdateRequest 服务  BigModelRoleUpdateService
type RoleUpdateRequest struct {
	ID        uint   `json:"id" structs:"-"`
	Name      string `json:"name" structs:"name"`            // 角色名称
	Enable    bool   `json:"enable" structs:"enable"`        // 是否启用
	Icon      string `json:"icon" structs:"icon"`            // 可以选择系统默认的一些，也可以图片上传
	Abstract  string `json:"abstract" structs:"abstract"`    // 简介
	Scope     int    `json:"scope" structs:"scope"`          // 消耗的积分
	Prologue  string `json:"prologue" structs:"prologue"`    // 开场白
	Prompt    string `json:"prompt" structs:"prompt"`        // 设定词
	AutoReply bool   `json:"autoReply" structs:"auto_reply"` // 自动回复
	TagList   []uint `json:"tagList" structs:"-"`            // 标签的id列表
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

func BigModelRoleUpdateService(cr RoleUpdateRequest) (int, error) {
	// 开启事务，保证操作原子性
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var tags []models.BigModelTagModel
	count := tx.Where("id in ?", cr.TagList).Find(&tags).RowsAffected
	if count != int64(len(cr.TagList)) {
		tx.Rollback()
		return 0, errors.New("部分标签不存在")
	}

	if cr.ID == 0 {
		// 添加
		var model models.BigModelRoleModel
		count := tx.Where("name = ?", cr.Name).Find(&model).RowsAffected
		if count != 0 {
			return 0, errors.New("名称已经存在")
		}
		role := models.BigModelRoleModel{
			Name:      cr.Name,
			Enable:    cr.Enable,
			Icon:      cr.Icon,
			Abstract:  cr.Abstract,
			Scope:     cr.Scope,
			Prologue:  cr.Prologue,
			Prompt:    cr.Prompt,
			AutoReply: cr.AutoReply,
			Tags:      tags,
		}
		err := tx.Create(&role).Error
		if err != nil {
			tx.Rollback()
			logrus.Errorf("%#v", err)
			return 0, errors.New("添加角色失败")
		}
		// 添加第三张表    已经删除了

		return 1, tx.Commit().Error // 显式提交事务
	}

	//修改
	var model models.BigModelRoleModel
	err := tx.Where("id = ?", cr.ID).Take(&model).Error
	if err != nil {
		tx.Rollback()
		logrus.Errorf("%#v", err)
		return 0, errors.New("角色不存在")
	}
	// 查看 修改后的 name  会不会重复
	var arm models.BigModelRoleModel
	err = tx.Where("name = ? and id <> ?", cr.Name, cr.ID).Take(&arm).Error
	if err == nil {
		return 0, errors.New("修改的名称已经存在")
	}

	toMap := struct_to_map.StructToMap(cr) //只修改 传入的
	// 先修改 角色表
	err = tx.Model(model).Updates(toMap).Error
	if err != nil {
		tx.Rollback()
		logrus.Errorf("%#v", err)
		return 0, errors.New("修改角色失败")
	}

	// 修改中间表  全删除  在新建  便是 修改
	var ThreeList []models.BigModelRoleTagModel
	tx.Where("big_model_role_model_id = ?", cr.ID).Find(&ThreeList)
	if len(ThreeList) > 0 {
		// 有 要删除
		err := tx.Delete(&ThreeList).Error
		if err != nil {
			tx.Rollback()
			logrus.Errorf("%#v", err)
			return 0, errors.New("修改失败")
		}
	}
	// 重新添加  由于刚开始已经搜索到了 tags
	for _, tagId := range cr.TagList {
		err = tx.Create(&models.BigModelRoleTagModel{
			BigModelTagModelId:  tagId,
			BigModelRoleModelId: cr.ID,
		}).Error
		if err != nil {
			tx.Rollback()
			logrus.Errorf("%#v", err)
			return 0, errors.New("关联表添加失败")
		}
	}
	return 2, tx.Commit().Error // 显式提交事务
}

func BigModelRoleRemoveService(cr models.RemoveRequest) (error, string) {
	// 开启事务，保证操作原子性
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//先看看有没有这个人
	var userModels []models.BigModelRoleModel
	tx.Where("id in ?", cr.IDList).Find(&userModels)
	if len(userModels) != len(cr.IDList) {
		tx.Rollback()
		logrus.Errorf("只查到了%d个角色", len(userModels))
		return errors.New("部分角色不存在"), ""
	}
	// 去找对应的 关联表
	err := tx.Where("big_model_role_model_id in ?", cr.IDList).Delete(&models.BigModelRoleTagModel{}).Error
	if err != nil {
		tx.Rollback()
		logrus.Errorf("%v", err)
		return errors.New("删除失败"), ""
	}
	//在删除人
	err = tx.Delete(&userModels).Error

	if err != nil {
		tx.Rollback()
		logrus.Errorf("%v", err)
		return errors.New("删除失败"), ""
	}
	return tx.Commit().Error, fmt.Sprintf("成功删除%d个角色", len(userModels))
}
