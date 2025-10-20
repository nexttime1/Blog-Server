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
	r.PUT("big_model/auto_reply", middleware.AdminMiddleware, app.AutoReplyUpdateView)
	r.GET("big_model/auto_reply", middleware.AdminMiddleware, app.AutoReplyListView)
	r.DELETE("big_model/auto_reply", middleware.AdminMiddleware, app.AutoReplyDeleteView)
	r.PUT("big_model/tags", middleware.AdminMiddleware, app.BigModelTagUpdateView)
	r.GET("big_model/tags", middleware.AdminMiddleware, app.BigModelTagListView)
	r.DELETE("big_model/tags", middleware.AdminMiddleware, app.BigModelTagRemoveView)
	r.PUT("big_model/roles", middleware.AdminMiddleware, app.BigModelRoleUpdateView)
	r.GET("big_model/roles", middleware.AdminMiddleware, app.BigModelRoleListiew)
	r.DELETE("big_model/roles", middleware.AdminMiddleware, app.BigModelRoleRemoveView)
	r.GET("big_model/square", app.BigModelTagRoleListView)
	r.POST("big_model/session", middleware.AuthMiddleware, app.BigModelSessionCreateView) //创建 会话
	r.POST("big_model/chat", middleware.AuthMiddleware, app.BigModelChatCreateView)       //创建 对话
}
