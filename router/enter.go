package router

import (
	"Blog_server/global"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()

	//路由分组
	nr := r.Group("/api")
	nr.Use(middleware.LogMiddleware)
	SiteRouter(nr)

	addr := global.Config.System.GetAddr()
	r.Run(addr)
}
