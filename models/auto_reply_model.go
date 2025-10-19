package models

import (
	"Blog_server/global"
	"regexp"
	"strings"
)

type AutoReplyModel struct {
	Model
	Name         string `json:"name" gorm:"size:32"`           // 规则名称
	Mode         int    `json:"mode"`                          // 匹配模式 1 精确匹配，2 模糊匹配，3 前缀匹配，4 正则匹配
	Rule         string `json:"rule" gorm:"size:32"`           // 匹配规则
	ReplyContent string `json:"replyContent" gorm:"size:1024"` // 回复内容
}

func (AutoReplyModel) AutoReplyBValid(content string) *AutoReplyModel {
	var list []AutoReplyModel

	// 查全部
	global.DB.Find(&list)
	for _, model := range list {
		switch model.Mode {
		case 1:
			// 精确匹配
			if model.Rule == content {
				return &model
			}
		case 2:
			// 包含
			if strings.Contains(content, model.Rule) {
				return &model
			}
		case 3:
			// 前缀
			if strings.HasPrefix(content, model.Rule) {
				return &model
			}
		case 4:
			// 正则
			regex, _ := regexp.Compile(model.Rule)
			if regex.MatchString(content) {
				return &model
			}

		}
	}
	return nil
}
