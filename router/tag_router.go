package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func TagRouter(c *gin.RouterGroup) {
	app := api.App.TagApi
	c.GET("tags", middleware.AuthMiddleware, app.TagListView)
	c.POST("tags", middleware.AuthMiddleware, app.TagAddView)
	c.PUT("tags/:id", middleware.AuthMiddleware, app.TagUpdateView)
	c.DELETE("tags", middleware.AuthMiddleware, app.TagDeleteView)

}
