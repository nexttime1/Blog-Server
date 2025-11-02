package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func LogRouter(r *gin.RouterGroup) {
	app := api.App.LogApi
	r.Use(middleware.AdminMiddleware)
	r.GET("logs", middleware.AdminMiddleware, app.LogListNew)
	r.GET("logs/read", middleware.AdminMiddleware, app.LogReadView)
	r.DELETE("logs", middleware.AdminMiddleware, app.LogRemoveView)
}
