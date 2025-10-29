package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func ImageRouter(r *gin.RouterGroup) {
	app := api.App.ImageApi
	r.GET("images", app.ImageInfoView)
	r.GET("image_names", app.ImageNameListView)
	r.POST("image", app.ImageUploadView)
	r.DELETE("images", app.ImageRemoveView)
	r.PUT("images", app.ImageUpdateView)

}
