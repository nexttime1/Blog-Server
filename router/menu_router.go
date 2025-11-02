package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func MenuRouter(r *gin.RouterGroup) {
	app := api.App.MenuApi
	r.POST("menus", middleware.AuthMiddleware, app.MenuCreateView)
	r.GET("menus", app.MenuListView)
	r.GET("menus_names", app.MenuNameListView)
	r.PUT("menus/:id", middleware.AuthMiddleware, app.MenuUpdateView)
	r.DELETE("menus", middleware.AuthMiddleware, app.MenuDeleteView)
	r.GET("menus/detail", app.MenuDetailView)

}
