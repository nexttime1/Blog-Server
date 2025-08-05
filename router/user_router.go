package router

import (
	"Blog_server/api"
	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {
	app := api.App.UserApi
	r.POST("email_login", app.UserEmailLogin)

}
