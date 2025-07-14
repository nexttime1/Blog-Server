package log_service

import (
	"Blog_server/core"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"fmt"
	"github.com/gin-gonic/gin"
)

func NewLoginSuccess(c *gin.Context, loginType enum.LoginType) {
	ip := c.ClientIP()
	addr := core.GetIpAddr(ip)

	token := c.GetHeader("token")
	fmt.Println("token:", token)
	//TODO 通过jwt 获取用户ID
	UserID := uint(1)
	UserName := "xxx"

	global.DB.Create(&models.LogModel{
		LogType:     enum.LoginLogType,
		Title:       "用户登录",
		Content:     "",
		UserID:      UserID,
		IP:          ip,
		Addr:        addr,
		LoginStatus: true,
		UserName:    UserName,
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
