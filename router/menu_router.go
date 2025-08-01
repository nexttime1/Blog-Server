package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func MenuRouter(r *gin.RouterGroup) {
	app := api.App.MenuApi
	r.POST("menus", app.MenuCreateView)
	r.GET("menus", app.MenuListView)
	r.GET("menus_name", app.MenuNameListView)

}
