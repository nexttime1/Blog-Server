package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func DiggRouter(c *gin.RouterGroup) {
	app := api.App.DiggApi
	c.POST("articles/digg", middleware.AuthMiddleware, app.DiggArticleView)
}
