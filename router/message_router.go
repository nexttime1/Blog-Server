package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func MessageRouter(c *gin.RouterGroup) {
	app := api.App.MessageApi
	c.POST("messages", middleware.AuthMiddleware, app.MessageAddView)
	c.GET("message_users/me", middleware.AuthMiddleware, app.MessageUserListByMeView)
	c.GET("messages", app.MessageListAllView)       // 静默
	c.GET("messages_record", app.MessageRecordView) //静默
	c.GET("message_users/record/me", middleware.AuthMiddleware, app.MessageUserRecordByMeView)
	c.DELETE("message_users", app.MessageRecordRemoveView)
	c.GET("message_users", app.MessageUserListView)

	c.GET("message_users/user", app.MessageUserListByUserView)
	c.GET("message_users/record", app.MessageUserRecordView)
}
