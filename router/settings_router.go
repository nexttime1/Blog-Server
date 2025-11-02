package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func SettingsRouter(nr *gin.RouterGroup) {
	app := api.App.SettingApi
	nr.GET("settings/:name", middleware.AuthMiddleware, app.SettingInfoView)
	nr.PUT("settings/:name", middleware.AdminMiddleware, app.SettingInfoUpdateView)
}
