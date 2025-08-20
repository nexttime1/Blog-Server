package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func NewRouter(c *gin.RouterGroup) {
	app := api.App.NewApi
	c.POST("news", app.NewListView)

}
