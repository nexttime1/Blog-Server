package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func SiteRouter(r *gin.RouterGroup) {
	app := api.App.SiteApi

	r.GET("site/:name", app.SiteInfoView)
	r.PUT("site", middleware.AdminMiddleware, app.SiteUpdateView)
}
