package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func AdvertRouter(c *gin.RouterGroup) {
	app := api.App.AdvertApi
	c.POST("adverts", middleware.AuthMiddleware, app.AdvertAddView)
	c.GET("adverts", middleware.AuthMiddleware, app.AdvertListView)
	c.PUT("adverts/:id", middleware.AuthMiddleware, app.AdvertUpdateView)
	c.DELETE("adverts", middleware.AuthMiddleware, app.AdvertDeleteView)

}
