package log_service

import (
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"github.com/gin-gonic/gin"
)

func NewLoginSuccess(c *gin.Context, loginType enum.LoginType, model models.UserModel, addr string) {
	ip := c.ClientIP()

	global.DB.Create(&models.LogModel{
		LogType:     enum.LoginLogType,
		Level:       enum.LogInfoLevel,
		Title:       "用户登录",
		Content:     "用户登录成功",
		UserID:      model.ID,
		IP:          ip,
		Addr:        addr,
		LoginStatus: true,
		UserName:    model.Username,
		Pwd:         "-",
		LoginType:   loginType,
	})

}
func NewLoginFail(c *gin.Context, loginType enum.LoginType, msg string, username string, pwd string, addr string) {
	ip := c.ClientIP()

	global.DB.Create(&models.LogModel{
		LogType:     enum.LoginLogType,
		Title:       "用户登录失败",
		Level:       enum.LogErrLevel,
		Content:     msg,
		IP:          ip,
		Addr:        addr,
		LoginStatus: false,
		UserName:    username,
		Pwd:         pwd,
		LoginType:   loginType,
	})
}

func LogoutSuccess(c *gin.Context, model models.UserModel, addr string) {
	ip := c.ClientIP()

	global.DB.Create(&models.LogModel{
		LogType:     enum.LoginLogType,
		Title:       "用户注销",
		Level:       enum.LogInfoLevel,
		Content:     "用户注销成功",
		UserID:      model.ID,
		IP:          ip,
		Addr:        addr,
		LoginStatus: true,
		UserName:    model.Username,
		Pwd:         "-",
		LoginType:   enum.LogoutType,
	})

}
