package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func RoleRouter(c *gin.RouterGroup) {
	app := api.App.RoleApi

	c.GET("role_ids", app.RoleIdListView)

}
