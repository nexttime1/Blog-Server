package router

import (
	"Blog_server/global"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func Run() {
	gin.SetMode(global.Config.System.GinMode)
	r := gin.Default()
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	//静态路由   一般设置一样的  前面的也就是重命名的意思  127.0.0.1:8080/uploads/a.txt
	r.Static("/uploads", "uploads")

	//路由分组
	nr := r.Group("/api")
	nr.Use(middleware.LogMiddleware)
	FeedbackRouter(nr)
	RoleRouter(nr)
	DateRouter(nr)
	ChatRouter(nr)
	NewRouter(nr)
	CommentRouter(nr)
	CollectRouter(nr)
	DiggRouter(nr)
	ArticleRouter(nr)
	MessageRouter(nr)
	TagRouter(nr)
	UserRouter(nr)
	MenuRouter(nr)
	SiteRouter(nr)
	AdvertRouter(nr)
	ImageRouter(nr)
	SettingsRouter(nr)
	LogRouter(nr)

	addr := global.Config.System.GetAddr()
	r.Run(addr)
}
