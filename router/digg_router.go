package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func DiggRouter(c *gin.RouterGroup) {
	app := api.App.DiggApi
	c.POST("digg/article", app.DiggArticleView)
}
