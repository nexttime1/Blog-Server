package comment_api

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/comment_service"
	"Blog_server/service/redis_service/redis_count"
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
	num := redis_count.NewUserDigg().Get(fmt.Sprintf("%d", claim.UserID))
	if num == 0 {
		// 说明这次是点赞
		redis_count.NewUserDigg().Set(fmt.Sprintf("%d", claim.UserID)) // 变成1
		//评论点赞数 + 1
		redis_count.NewCommentDigg().Set(fmt.Sprintf("%d", commentModel.ID))
		res.OkWithMessage(c, "评论点赞成功")
		return
	}
	// 这次 用户取消点赞
	redis_count.NewUserDigg().Sub(fmt.Sprintf("%d", claim.UserID)) // 变成0
	//评论点赞数 - 1
	redis_count.NewCommentDigg().Sub(fmt.Sprintf("%d", commentModel.ID))
	res.OkWithMessage(c, "取消点赞成功")

}
