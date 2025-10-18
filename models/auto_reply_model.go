package models

type AutoReplyModel struct {
	Model
	Name         string `json:"name" gorm:"size:32"`           // 规则名称
	Mode         int    `json:"mode"`                          // 匹配模式 1 精确匹配，2 模糊匹配，3 前缀匹配，4 正则匹配
	Rule         string `json:"rule" gorm:"size:32"`           // 匹配规则
	ReplyContent string `json:"replyContent" gorm:"size:1024"` // 回复内容
}
