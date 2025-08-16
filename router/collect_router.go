package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRouter(c *gin.RouterGroup) {
	app := api.App.CollectApi
	c.POST("/articles/collects", middleware.AuthMiddleware, app.ArticleCollectView)
	c.GET("/articles/collects", middleware.AuthMiddleware, app.UserCollectListView)
	c.DELETE("/articles/collects", middleware.AuthMiddleware, app.UserCollectBatchDeleteView)
}
