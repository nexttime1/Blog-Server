package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func ImageRouter(r *gin.RouterGroup) {
	app := api.App.ImageApi
	r.GET("images", app.ImageInfoView)
	r.GET("image_names", app.ImageNameListView)
	r.POST("images", middleware.AuthMiddleware, app.ImageUploadView)
	r.POST("image", middleware.AuthMiddleware, app.ImageUploadView)
	r.DELETE("images", middleware.AuthMiddleware, app.ImageRemoveView)
	r.PUT("images", middleware.AuthMiddleware, app.ImageUpdateView)

}
