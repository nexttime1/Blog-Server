package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func SettingsRouter(nr *gin.RouterGroup) {
	app := api.App.SettingApi
	nr.GET("/settings", app.SettingInfoView)
}
