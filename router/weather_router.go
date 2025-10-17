package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func WeatherRouter(r *gin.RouterGroup) {
	app := api.App.GaodeApi
	r.GET("gaode/weather", app.WeatherInfoView)
}
