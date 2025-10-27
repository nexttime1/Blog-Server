package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func BigModelRouter(r *gin.RouterGroup) {
	app := api.App.BigModelApi
	// 大模型相关
	{
		r.GET("/big_model/usable", middleware.AdminMiddleware, app.BigModelOptionListView)             // 可用大模型列表
		r.GET("/big_model/setting", middleware.AuthMiddleware, app.BigModelSettingView)                // 获得大模型配置
		r.PUT("/big_model/setting", middleware.AdminMiddleware, app.BigModelUpdateView)                // 修改大模型配置
		r.GET("/big_model/session_setting", middleware.AuthMiddleware, app.BigModelSessionView)        // 获取大模型会话配置
		r.PUT("/big_model/session_setting", middleware.AdminMiddleware, app.BigModelSessionUpdateView) // 修改大模型会话配置
		r.PUT("big_model/auto_reply", middleware.AdminMiddleware, app.AutoReplyUpdateView)             //  新增或更新大模型自动回复
		r.GET("big_model/auto_reply", middleware.AdminMiddleware, app.AutoReplyListView)               // 获取大模型自动回复列表
		r.DELETE("big_model/auto_reply", middleware.AdminMiddleware, app.AutoReplyDeleteView)          //批量删除大模型自动回复
		r.PUT("big_model/tags", middleware.AdminMiddleware, app.BigModelTagUpdateView)                 //新增或修改大模型标签
		r.GET("big_model/tags", middleware.AdminMiddleware, app.BigModelTagListView)                   //获取大模型标签分页列表
		r.DELETE("big_model/tags", middleware.AdminMiddleware, app.BigModelTagRemoveView)              //批量删除大模型标签
		r.GET("big_model/tags/options", middleware.AdminMiddleware, app.BigModelRoleTagsListView)      // 角色标签id列表
	}
	// 用户相关
	{
		r.GET("big_model/user_scope_enable", middleware.AuthMiddleware, app.UserScopeEnableView)     // 用户是否可以领取积分
		r.POST("big_model/user_scope", middleware.AuthMiddleware, app.UserScopeView)                 // 用户 领取积分
		r.GET("big_model/role_sessions", middleware.AuthMiddleware, app.BigModelRoleSessionListView) //当前用户查询当前角色的会话列表
	}

	{
		r.GET("big_model/icons/options", app.IconView) //角色可选 图标
	}

	//大模型角色相关
	{
		r.PUT("big_model/roles", middleware.AdminMiddleware, app.BigModelRoleUpdateView)    // 创建或者修改大模型角色
		r.GET("big_model/roles", middleware.AdminMiddleware, app.BigModelRoleListView)      // 大模型角色列表
		r.DELETE("big_model/roles", middleware.AdminMiddleware, app.BigModelRoleRemoveView) //
		r.GET("big_model/roles/:id", app.BigModelRoleDetailView)                            // 角色细节
		r.GET("big_model/square", app.BigModelTagRoleListView)                              //大模型角色广场

	}
	//大模型会话相关
	{
		r.POST("big_model/session", middleware.AuthMiddleware, app.BigModelSessionCreateView)           //创建 会话
		r.GET("big_model/session", middleware.AdminMiddleware, app.BigModelSessionListView)             //会话列表
		r.PUT("big_model/session", middleware.AuthMiddleware, app.BigModelUserUpdateNameView)           // 修改会话名称
		r.DELETE("big_model/session/:id", middleware.AuthMiddleware, app.BigModelUserDeleteSessionView) //删除会话
		r.DELETE("big_model/session", middleware.AuthMiddleware, app.BigModelAdminDeleteSessionView)    //管理员批量删除会话
	}
	//对话相关
	{
		r.GET("big_model/chat_sse", middleware.AuthSSEMiddleware, app.BigModelChatCreateView)        //创建 对话
		r.GET("big_model/roles_history", middleware.AuthMiddleware, app.BigModelUserRoleHistoryView) // 用户和使用过的大模型聊天记录
		r.GET("big_model/chat", middleware.AuthMiddleware, app.BigModelChatListView)                 //单个会话的聊天记录
		r.DELETE("big_model/chat/:id", middleware.AuthMiddleware, app.BigModelUserChatDeleteView)    // 用户删除对话
		r.DELETE("big_model/chat", middleware.AdminMiddleware, app.BigModelAdminChatDeleteView)      // 管理员批量删除对话
	}

}
