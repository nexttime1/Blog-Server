package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func AdvertRouter(c *gin.RouterGroup) {
	app := api.App.AdvertApi
	c.POST("adverts", app.AdvertAddView)

}
