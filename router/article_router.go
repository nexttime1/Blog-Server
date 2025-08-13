package router

import (
	"Blog_server/api"
	"Blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func ArticleRouter(c *gin.RouterGroup) {
	app := api.App.ArticleApi
	c.POST("articles", middleware.AuthMiddleware, app.ArticleCreateView)
	c.GET("articles", app.ArticleListView)
	c.GET("articles/:id", app.ArticleDetailByIdView)
	c.GET("articles_detail_title", app.ArticleDetailByTitleView)
	c.GET("articles/calendar", app.ArticleCalendarView)
	c.GET("articles/tags", app.ArticleTagListView)

}
