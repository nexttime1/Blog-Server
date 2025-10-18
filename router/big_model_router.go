package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func BigModelRouter(r *gin.RouterGroup) {
	app := api.App.BigModelApi
	r.GET("/big_model/usable", middleware.AdminMiddleware, app.BigModelOptionListView)             // 可用大模型配置
	r.GET("/big_model/setting", middleware.AuthMiddleware, app.BigModelSettingView)                // 获得大模型配置
	r.PUT("/big_model/setting", middleware.AdminMiddleware, app.BigModelUpdateView)                // 修改大模型配置
	r.GET("/big_model/session_setting", middleware.AuthMiddleware, app.BigModelSessionView)        // 获取大模型会话配置
	r.PUT("/big_model/session_setting", middleware.AdminMiddleware, app.BigModelSessionUpdateView) // 修改大模型会话配置
	r.GET("big_model/user_scope_enable", middleware.AuthMiddleware, app.UserScopeEnableView)       // 用户是否可以领取积分
	r.POST("big_model/user_scope", middleware.AuthMiddleware, app.UserScopeView)                   // 用户 领取积分
}
