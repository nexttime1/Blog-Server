package models

type LogModel struct {
	Model
	LogType   int8      `json:"logType"` //日志类型 1 2 3
	Title     string    `gorm:"size:64" json:"title"`
	Content   string    `json:"content"`
	Level     int8      `json:"level"`  //日志级别  1 2 3
	UserID    uint      `json:"userID"` //用户id   可以没有  没登录  设置为0
	UserModel UserModel `gorm:"ForeignKey:UserID" json:"-"`
	IP        string    `gorm:"size:32" json:"ip"`
	Addr      string    `gorm:"size:64" json:"addr"`
	IsRead    bool      `json:"isRead"` //是否读取

}
