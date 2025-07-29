package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func AdvertRouter(c *gin.RouterGroup) {
	app := api.App.AdvertApi
	c.POST("adverts", app.AdvertAddView)
	c.GET("adverts", app.AdvertListView)
	c.PUT("adverts/:id", app.AdvertUpdateView)
	c.DELETE("adverts", app.AdvertDeleteView)

}
