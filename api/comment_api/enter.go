package comment_api

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/comment_service"
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
		fmt.Println(model)
		var subCommentModels []*models.CommentModel
		xxx(model, &subCommentModels)
		model.SubComments = subCommentModels
	}
	res.OkWithData(c, filter.Select("c", ParentsModels))
}

func xxx(model *models.CommentModel, subCommentModels *[]*models.CommentModel) {
	global.DB.Preload("SubComments.User").Take(model)
	fmt.Println(model.SubComments)
	*subCommentModels = append(*subCommentModels, model.SubComments...)

	for _, commentModel := range model.SubComments {
		xxx(commentModel, subCommentModels)
	}

}
