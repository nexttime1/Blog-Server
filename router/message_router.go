package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func MessageRouter(c *gin.RouterGroup) {
	app := api.App.MessageApi
	c.POST("messages", app.MessageAddView)
	c.GET("messages", app.MessageListAllView)
	c.GET("messages_record", app.MessageRecordView)
	c.GET("message_users/record/me", app.MessageUserRecordByMeView)
	c.DELETE("message_users", app.MessageRecordRemoveView)
	c.GET("message_users", app.MessageUserListView)
	c.GET("message_users/me", app.MessageUserListByMeView)
	c.GET("message_users/user", app.MessageUserListByUserView)
	c.GET("message_users/record", app.MessageUserRecordView)
}
