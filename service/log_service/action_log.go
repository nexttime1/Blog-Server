package log_service

import (
	"Blog_server/core"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ActionLog struct {
	c       *gin.Context
	level   enum.LevelType
	title   string
	content string
}

func (ac *ActionLog) SetTitle(title string) {
	ac.title = title
}

func (ac *ActionLog) SetLevel(level enum.LevelType) {
	ac.level = level
}

func (ac *ActionLog) Save() {
	ip := ac.c.ClientIP()
	addr := core.GetIpAddr(ip)
	UserID := uint(0)
	err := global.DB.Create(&models.LogModel{
		LogType: enum.ActionLogType,
		Title:   ac.title,
		Content: "",
		Level:   ac.level,
		UserID:  UserID,
		IP:      ip,
		Addr:    addr,
	}).Error
	if err != nil {
		logrus.Errorf("日志创建失败 %s", err.Error())
		return
	}

}

func NewActionLogByGin(c *gin.Context) *ActionLog {
	return &ActionLog{c: c}
}
