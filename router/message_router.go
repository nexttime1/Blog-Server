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

}
