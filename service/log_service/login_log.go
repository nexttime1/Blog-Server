package log_service

import (
	"Blog_server/core"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"github.com/gin-gonic/gin"
)

func NewLoginSuccess(c *gin.Context, loginType enum.LoginType, model models.UserModel) {
	ip := c.ClientIP()
	addr := core.GetIpAddr(ip)

	global.DB.Create(&models.LogModel{
		LogType:     enum.LoginLogType,
		Title:       "用户登录",
		Content:     "",
		UserID:      model.ID,
		IP:          ip,
		Addr:        addr,
		LoginStatus: true,
		UserName:    model.Username,
		Pwd:         "-",
		LoginType:   loginType,
	})

}
func NewLoginFail(c *gin.Context, loginType enum.LoginType, msg string, username string, pwd string) {
	ip := c.ClientIP()
	addr := core.GetIpAddr(ip)

	global.DB.Create(&models.LogModel{
		LogType:     enum.LoginLogType,
		Title:       "用户登录失败",
		Content:     msg,
		IP:          ip,
		Addr:        addr,
		LoginStatus: false,
		UserName:    username,
		Pwd:         pwd,
		LoginType:   loginType,
	})
}
