package router

import (
	"Blog_server/global"
	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()

	//路由分组
	nr := r.Group("/api")
	SiteRouter(nr)

	addr := global.Config.System.GetAddr()
	r.Run(addr)
}
