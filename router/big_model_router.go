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
	r.POST("big_model/session", middleware.AuthMiddleware, app.BigModelSessionCreateView)           //创建 会话
	r.POST("big_model/chat", middleware.AuthMiddleware, app.BigModelChatCreateView)                 //创建 对话
	r.GET("big_model/session", middleware.AdminMiddleware, app.BigModelSessionListView)             //会话列表
	r.PUT("big_model/session", middleware.AuthMiddleware, app.BigModelUserUpdateNameView)           // 修改会话名称
	r.DELETE("big_model/session/:id", middleware.AuthMiddleware, app.BigModelUserDeleteSessionView) //删除会话
	r.DELETE("big_model/session", middleware.AuthMiddleware, app.BigModelAdminDeleteSessionView)    //管理员批量删除会话
	r.GET("big_model/roles/:id", app.BigModelRoleDetailView)                                        // 角色细节
	r.GET("big_model/roles_history", middleware.AuthMiddleware, app.BigModelUserRoleHistoryView)    // 用户和使用过的大模型聊天记录
	r.GET("big_model/chat", middleware.AuthMiddleware, app.BigModelChatListView)                    //单个会话的聊天记录
	r.DELETE("big_model/chat/:id", middleware.AuthMiddleware, app.BigModelUserChatDeleteView)       // 用户删除对话
	r.DELETE("big_model/chat", middleware.AdminMiddleware, app.BigModelAdminChatDeleteView)         // 管理员批量删除对话
}
