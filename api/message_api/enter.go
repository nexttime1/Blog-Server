package message_api

import (
	"Blog_server/common"
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/message_service"
	"github.com/gin-gonic/gin"
)

type MessageApi struct {
}

// MessageAddView 发送消息
// @Summary 发送消息
// @Description 发送一个新的消息
// @Tags 消息管理
// @Accept json
// @Produce json
// @Param data body message_service.MessageAddRequest true "信息"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response "发送成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/messages [post]
func (MessageApi) MessageAddView(c *gin.Context) {
	var cr message_service.MessageAddRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	err = message_service.MessageAddService(cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	res.OkWithMessage(c, "发送消息成功")

}

// MessageListAllView 显示全部消息
// @Summary 显示全部消息
// @Description 显示全部消息
// @Tags 消息管理
// @Produce json
// @Param page query int false "页码，默认1" mininum(1)
// @Param limit query int false "每页条数，默认10" mininum(1) maxinum(100)
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response "操作成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/messages [get]
func (MessageApi) MessageListAllView(c *gin.Context) {
	var cr common.PageInfo
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	list, count, err := common.ListQuery(models.MessageModel{}, common.Options{
		PageInfo: cr,
	})
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithList(c, list, count)
}

// MessageRecordView 显示全部聊天记录
// @Summary 显示全部聊天记录
// @Description 显示全部聊天记录
// @Tags 消息管理
// @Accept json
// @Produce json
// @Param data body message_service.MessageRecordRequest true "输入要看谁的聊天记录"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response "操作成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/messages_record [get]
func (MessageApi) MessageRecordView(c *gin.Context) {
	var cr message_service.MessageRecordRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
	}
	var MessageModels []models.MessageModel
	global.DB.Where("send_user_id = ? or rev_user_id = ?", cr.UserID, cr.UserID).Find(&MessageModels)
	if len(MessageModels) == 0 {
		res.OkWithMessage(c, "并没有消息")
		return
	}
	res.OkWithData(c, MessageModels)
}
