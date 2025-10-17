package feedback_api

import (
	"Blog_server/common"
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"github.com/gin-gonic/gin"
)

type FeedBackApi struct {
}

type FeedBackAddRequest struct {
	Email   string `json:"email"`
	Content string `json:"content"`
}

func (FeedBackApi) FeedBackAddView(c *gin.Context) {

	var cr FeedBackAddRequest

	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithErr(c, err)
		return
	}

	err := global.DB.Create(&models.FeedbackModel{
		Email:   cr.Email,
		Content: cr.Content,
	}).Error
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithMessage(c, "提交成功")

}

func (FeedBackApi) FeedBackListView(c *gin.Context) {
	var cr common.PageInfo
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	list, count, err := common.ListQuery(models.FeedbackModel{}, common.Options{
		PageInfo: cr,
	})

	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithList(c, list, count)

}

func (FeedBackApi) FeedBackRemoveView(c *gin.Context) {
	var cr models.IDRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	err = global.DB.Delete(&models.FeedbackModel{}, "id = ?", cr.ID).Error
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithMessage(c, "删除成功")

}
