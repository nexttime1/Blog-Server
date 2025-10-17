package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func CommentRouter(c *gin.RouterGroup) {
	app := api.App.CommentApi
	c.GET("comments/articles", app.CommentByArticleListView)
	c.POST("comments", middleware.AuthMiddleware, app.CommentAddView)
	c.GET("comments/:id", middleware.AuthMiddleware, app.CommentListView)
	c.GET("comments/digg/:id", middleware.AuthMiddleware, app.CommentDiggView)
	c.DELETE("comments/:id", middleware.AuthMiddleware, app.CommentDeleteView)

}
