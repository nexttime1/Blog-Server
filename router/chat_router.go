package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func ChatRouter(c *gin.RouterGroup) {
	app := api.App.ChatApi

	c.GET("chat_group", app.ChatGroupView)
	c.GET("chat_groups_record", app.ChatListView)
	c.DELETE("chat_groups_records", app.ChatRemoveView)
}
