package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var store = cookie.NewStore([]byte("XIAOSONGAICHIBINGxtm666"))

func UserRouter(r *gin.RouterGroup) {
	app := api.App.UserApi
	r.Use(sessions.Sessions("sessionid", store))
	r.POST("email_login", app.UserEmailLogin)
	r.POST("qq_login", app.UserQQLogin)
	r.GET("users", middleware.AuthMiddleware, app.UserListInfoView)
	r.PUT("user_role", middleware.AdminMiddleware, app.UserUpdateView)
	r.PUT("users_pwd", middleware.AuthMiddleware, app.UserPasswordView)
	r.POST("users", middleware.AdminMiddleware, app.UserCreateView)
	//r.GET("user_info", middleware.AuthMiddleware, app.UserPersonInfoView)
	r.DELETE("users", middleware.AdminMiddleware, app.UserDeleteView)
	r.POST("user_bind_email", middleware.AdminMiddleware, app.UserBindEmailView)

	r.PUT("user_info", middleware.AuthMiddleware, app.UserUpdateInfoView)
	r.GET("user_info", middleware.AuthMiddleware, app.UserInfoView) //个人信息
	r.GET("qq_login_path", app.QQLoginLinkView)
	r.POST("logout", middleware.AuthMiddleware, app.UserLogoutView)
}
