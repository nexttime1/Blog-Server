package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func SiteRouter(r *gin.RouterGroup) {
	app := api.App.SiteApi

	r.GET("site", app.SiteInfoView)
	r.PUT("site", app.SiteUpdateView)
}
