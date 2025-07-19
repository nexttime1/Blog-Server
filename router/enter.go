package router

import (
	"Blog_server/global"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func Run() {
	gin.SetMode(global.Config.System.GinMode)
	r := gin.Default()

	//静态路由   一般设置一样的  前面的也就是重命名的意思  127.0.0.1:8080/uploads/a.txt
	r.Static("/uploads", "uploads")

	//路由分组
	nr := r.Group("/api")
	nr.Use(middleware.LogMiddleware)
	SiteRouter(nr)
	LogRouter(nr)

	addr := global.Config.System.GetAddr()
	r.Run(addr)
}
