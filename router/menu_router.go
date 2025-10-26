package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func MenuRouter(r *gin.RouterGroup) {
	app := api.App.MenuApi
	r.POST("menus", app.MenuCreateView)
	r.GET("menus", app.MenuListView)
	r.GET("menus_names", app.MenuNameListView)
	r.PUT("menus/:id", app.MenuUpdateView)
	r.DELETE("menus", app.MenuDeleteView)
	r.GET("menus/detail", app.MenuDetailView)

}
