package models

import "Blog_server/models/enum"

type ChatModel struct {
	Model    `json:","`
	NickName string       `gorm:"size:15" json:"nick_name"`
	Avatar   string       `gorm:"size:128" json:"avatar"`
	Content  string       `gorm:"size:256" json:"content"`
	IP       string       `gorm:"size:32" json:"ip,omit(list)"`
	Addr     string       `gorm:"size:64" json:"addr,omit(list)"`
	IsGroup  bool         `json:"is_group"` // 是否是群组消息
	MsgType  enum.MsgType `gorm:"size:4" json:"msg_type"`
}
