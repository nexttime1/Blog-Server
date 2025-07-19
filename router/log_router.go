package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func LogRouter(r *gin.RouterGroup) {
	app := api.App.LogApi

	r.GET("logs", app.LogListNew)
}
