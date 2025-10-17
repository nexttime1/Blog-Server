package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func FeedbackRouter(r *gin.RouterGroup) {
	app := api.App.FeedBackApi
	r.POST("feedback", app.FeedBackAddView)
	r.GET("feedback", app.FeedBackListView)
	r.DELETE("feedback", app.FeedBackRemoveView)
}
