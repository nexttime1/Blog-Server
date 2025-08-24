package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func DateRouter(c *gin.RouterGroup) {
	app := api.App.DateApi
	c.GET("data_login", app.SevenLoginView)
	c.GET("data_sum", app.DataSumView)
}
