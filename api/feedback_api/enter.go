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

// FeedBackAddView 提交反馈
// @Summary 提交反馈信息
// @Description 用户提交反馈内容，需包含邮箱和具体反馈信息
// @Tags 反馈管理
// @Accept application/json
// @Produce json
// @Param token header string true "用户认证令牌"
// @Param data body FeedBackAddRequest true "反馈信息（邮箱和内容）"
// @Success 200 {object} res.Response "提交成功"
// @Failure 400 {object} res.Response "请求参数错误（如邮箱格式错误、内容为空等）"
// @Failure 500 {object} res.Response "服务器错误（如数据库存储失败）"
// @Router /feedback [post]
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

// FeedBackListView 获取反馈列表
// @Summary 分页获取反馈列表
// @Description 分页查询所有用户提交的反馈记录，支持页码和每页条数参数
// @Tags 反馈管理
// @Accept application/json
// @Produce json
// @Param token header string true "用户认证令牌"
// @Param page query int false "页码（默认1）"
// @Param limit query int false "每页条数（默认10）"
// @Success 200 {object} res.Response{data=res.DataListResponse[]}
// @Failure 400 {object} res.Response "请求参数错误（如页码/条数为负数）"
// @Failure 500 {object} res.Response "服务器错误（如数据库查询失败）"
// @Router /feedback [get]
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

// FeedBackRemoveView 删除反馈
// @Summary 删除指定反馈
// @Description 根据反馈ID删除对应的反馈记录
// @Tags 反馈管理
// @Accept application/json
// @Produce json
// @Param token header string true "用户认证令牌"
// @Param data body models.IDRequest true "反馈记录ID"
// @Success 200 {object} res.Response "删除成功"
// @Failure 400 {object} res.Response "请求参数错误（如ID为空或格式错误）"
// @Failure 500 {object} res.Response "服务器错误（如数据库删除失败）"
// @Router /feedback [delete]
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
