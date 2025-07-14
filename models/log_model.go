package models

import "Blog_server/models/enum"

type LogModel struct {
	Model
	LogType     enum.LogType   `json:"logType"` //日志类型 1 2 3
	Title       string         `gorm:"size:64" json:"title"`
	Content     string         `json:"content"`
	Level       enum.LevelType `json:"level"`  //日志级别  1 2 3
	UserID      uint           `json:"userID"` //用户id   可以没有  没登录  设置为0
	UserModel   UserModel      `gorm:"ForeignKey:UserID" json:"-"`
	IP          string         `gorm:"size:32" json:"ip"`
	Addr        string         `gorm:"size:64" json:"addr"`
	IsRead      bool           `json:"isRead"`      //是否读取
	LoginStatus bool           `json:"loginStatus"` //登录状态
	UserName    string         `gorm:"size:32" json:"userName"`
	Pwd         string         `gorm:"size:32" json:"pwd"`
	LoginType   enum.LoginType `json:"loginType"`
}
