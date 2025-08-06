package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {
	app := api.App.UserApi
	r.POST("email_login", app.UserEmailLogin)
	r.GET("users", middleware.AuthMiddleware, app.UserInfoView)
	r.PUT("users_role", middleware.AdminMiddleware, app.UserUpdateView)
	r.PUT("users_pwd", middleware.AuthMiddleware, app.UserPasswordView)

}
