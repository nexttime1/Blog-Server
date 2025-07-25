package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func ImageRouter(r *gin.RouterGroup) {
	app := api.App.ImageApi
	r.POST("images", app.ImageUploadView)

}
