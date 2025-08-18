package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func CommentRouter(c *gin.RouterGroup) {
	app := api.App.CommentApi
	c.POST("comments", middleware.AuthMiddleware, app.CommentAddView)
	c.GET("comments", middleware.AuthMiddleware, app.CommentListView)
}
