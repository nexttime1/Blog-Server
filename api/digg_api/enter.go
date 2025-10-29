package digg_api

import (
	"Blog_server/common/res"
	"Blog_server/models"
	"Blog_server/service/redis_service/redis_count"
	"github.com/gin-gonic/gin"
)

type DiggApi struct {
}

// DiggArticleView 文章点赞
// @Summary 点赞文章
// @Description 点赞文章
// @Tags 文章点赞管理
// @Accept json
// @Produce json
// @Param data body models.EsIdQuest false "文章id"
// @Param token header string true "token"
// @Success 200 {object} res.Response "操作成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/articles/digg [post]
func (DiggApi) DiggArticleView(c *gin.Context) {
	_, exists := c.Get("claims")
	if !exists {
		return
	}
	var cr models.EsIdQuest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	redis_count.NewDigg().Set(cr.ID)
	res.OkWithMessage(c, "文章点赞成功")

}
