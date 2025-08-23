package comment_api

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/comment_service"
	"Blog_server/service/redis_service/redis_count"
	"Blog_server/service/redis_service/redis_user"
	"Blog_server/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liu-cn/json-filter/filter"
	"github.com/sirupsen/logrus"
)

type CommentApi struct {
}
type CommentRequest struct {
	ArticleID string `form:"article_id"`
}
type CommentDiggRequest struct {
	ID int `form:"id" uri:"id"`
}

// CommentAddView 添加评论
// @Summary 添加评论
// @Description 创建一个新的评论，包含文章，内容和父评论Id
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param data body comment_service.CommentAddRequest true "评论信息"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response "创建评论成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/comments [post]
func (CommentApi) CommentAddView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)

	var cr comment_service.CommentAddRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	err = comment_service.CommentAddService(cr, claim)
	if err != nil {
		logrus.Errorf("%s", err)
		res.FailWithMsg(c, fmt.Sprintf("%s", err))
		return
	}
	res.OkWithMessage(c, "创建评论成功")
}

// CommentListView 评论列表
// @Summary 评论列表
// @Description 获取某个文章的评论列表
// @Tags 评论管理
// @Produce json
// @Param article_id query string true "输入文章ID"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{data=res.DataListResponse}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/comments [get]
func (CommentApi) CommentListView(c *gin.Context) {
	var cr CommentRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	fmt.Println(cr.ArticleID)
	//所有根评论
	var ParentsModels []*models.CommentModel

	global.DB.Preload("User").Where("article_id = ? and parent_comment_id is null", cr.ArticleID).Find(&ParentsModels)
	fmt.Println(len(ParentsModels))
	for _, model := range ParentsModels {
		var subCommentModels []*models.CommentModel
		Recursion(model, &subCommentModels)
		model.SubComments = subCommentModels
	}
	res.OkWithData(c, filter.Select("c", ParentsModels))
}

func Recursion(model *models.CommentModel, subCommentModels *[]*models.CommentModel) {
	CommentDiggList := redis_count.NewCommentDigg().GetInfo()

	global.DB.Preload("SubComments.User").Take(model)

	model.DiggCount = model.DiggCount + CommentDiggList[fmt.Sprintf("%d", model.ID)]

	for _, commentModel := range model.SubComments {
		*subCommentModels = append(*subCommentModels, commentModel)
		Recursion(commentModel, subCommentModels)
		commentModel.SubComments = nil
	}

}

// CommentDiggView 用户点赞评论
// @Summary 用户点赞评论
// @Description 用户点赞评论
// @Tags 评论管理
// @Produce json
// @Param id path int true "id"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/comments/{id} [get]
func (CommentApi) CommentDiggView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	var cr CommentDiggRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	//查看评论存不存咋
	var commentModel models.CommentModel
	err = global.DB.Where("id = ?", cr.ID).Take(&commentModel).Error
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("评论不存在 %s", err))
	}
	//判断该用户点没点赞
	list := redis_user.NewUserDigg().Get(fmt.Sprintf("%d", claim.UserID))

	if !redis_user.Exist(list, commentModel.ArticleID) {
		// 说明这次是点赞  不存在
		redis_user.NewUserDigg().Add(fmt.Sprintf("%d", claim.UserID), commentModel.ArticleID) // 在列表中增加
		//评论点赞数 + 1
		redis_count.NewCommentDigg().Set(fmt.Sprintf("%d", commentModel.ID))
		res.OkWithMessage(c, "评论点赞成功")
		return
	}
	// 这次 用户取消点赞
	redis_user.NewUserDigg().Del(fmt.Sprintf("%d", claim.UserID), commentModel.ArticleID) //  在列表中减去
	//评论点赞数 - 1
	redis_count.NewCommentDigg().SetNum(fmt.Sprintf("%d", commentModel.ID), -1)
	res.OkWithMessage(c, "取消点赞成功")

}

// CommentDeleteView 评论删除
// @Summary 删除评论
// @Description 删除评论
// @Tags 评论管理
// @Produce json
// @Param id path string true "输入删除的评论id"
// @Param token header string true "token"
// @Success 200 {object} res.Response "删除成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 404 {object} res.Response "评论id不存在"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/comments/{id} [delete]
func (CommentApi) CommentDeleteView(c *gin.Context) {
	_, exists := c.Get("claims")
	if !exists {
		return
	}
	var cr comment_service.CommentDeleteRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithErr(c, err)
	}
	count, err := comment_service.CommentDeleteService(cr)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%s", err))
	}
	res.OkWithData(c, fmt.Sprintf("共删除%d条评论 ", count))
}
