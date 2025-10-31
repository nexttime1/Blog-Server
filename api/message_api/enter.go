package message_api

import (
	"Blog_server/common"
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/message_service"
	"Blog_server/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MessageApi struct {
}

type MessageByMeRequest struct {
	common.PageInfo
	UserID string `json:"userID" form:"userID" uri:"userID"`
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
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	logrus.Infof("claim = %v", claim)
	var cr message_service.MessageAddRequest
	err := c.ShouldBindJSON(&cr)
	cr.SendUserID = claim.UserID
	logrus.Infof("cr: %v", cr)
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
// @Param page query int false "页码，默认1" mininum(1)
// @Param limit query int false "每页条数，默认10" mininum(1) maxinum(100)
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

// MessageUserRecordByMeView 显示自己和指定用户的聊天记录
// @Summary 显示自己和指定用户的聊天记录
// @Description 获取当前登录用户与指定用户之间的聊天记录，支持分页和搜索
// @Tags 消息管理
// @Produce json
// @Param userID query string true "目标用户ID"
// @Param page query int false "页码，默认1" default(1)
// @Param limit query int false "每页条数，默认10" default(10)
// @Param key query string false "搜索关键词"
// @Param order query string false "排序方式，如id desc"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{data=res.DataListResponse[list = []models.MessageModel]}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/message_users/record/me [get]
func (MessageApi) MessageUserRecordByMeView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	var cr MessageByMeRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	list, count, err := common.ListQuery(models.MessageModel{}, common.Options{
		PageInfo: cr.PageInfo,
		Where:    global.DB.Where("(send_user_id = ? and rev_user_id = ?) or (rev_user_id = ? and send_user_id = ?)", claim.UserID, cr.UserID, claim.UserID, cr.UserID),
		Debug:    true,
	})
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithList(c, list, count)

}

// MessageRecordRemoveView 删除用户的消息记录
// @Tags 消息管理
// @Summary 删除用户的消息记录
// @Description 删除用户的消息记录
// @Router /api/message_users [delete]
// @Param token header string  true  "token"
// @Param data body models.RemoveRequest   true  "查询参数"
// @Produce json
// @Success 200 {object} res.Response{]}
func (MessageApi) MessageRecordRemoveView(c *gin.Context) {
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithCode(c, res.ArgumentError)
		return
	}

	var messageList []models.MessageModel
	global.DB.Find(&messageList, cr.IDList)

	if len(messageList) > 0 {
		err = global.DB.Delete(&messageList).Error
		if err != nil {
			res.FailWithMsg(c, "消息记录删除失败")
			return
		}
	}

	res.OkWithMessage(c, fmt.Sprintf("共删除记录%d条", len(messageList)))
}

type MessageUserListRequest struct {
	common.PageInfo
	NickName string `json:"nickName" form:"nickName"`
}

type MessageUserListResponse struct {
	UserName string `json:"userName"`
	NickName string `json:"nickName"`
	UserID   uint   `json:"userID"`
	Avatar   string `json:"avatar"`
	Count    int    `json:"count"`
}

// MessageUserListView 有消息的用户列表
// @Tags 消息管理
// @Summary 有消息的用户列表
// @Description 有消息的用户列表  不传这个nickname  查全部  也就是 页面中的 左中右  中的左
// @Router /api/message_users [get]
// @Param token header string  true  "token"
// @Param data query MessageUserListRequest   false  "查询参数"
// @Produce json
// @Success 200 {object} res.Response{data=res.DataListResponse[list = MessageUserListResponse[]]}
func (MessageApi) MessageUserListView(c *gin.Context) {
	var cr MessageUserListRequest
	err := c.ShouldBindQuery(&cr)

	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	var count int64

	global.DB.Model(models.MessageModel{}).Where(models.MessageModel{SendUserNickName: cr.NickName}).
		Group("send_user_id").Count(&count)

	type resType struct {
		SendUserID uint
		Count      int // 发送人的个数2
	}
	offset := (cr.Page - 1) * cr.Limit

	var _list []resType
	global.DB.Model(models.MessageModel{}).Where(models.MessageModel{SendUserNickName: cr.NickName}).
		Group("send_user_id").Limit(cr.Limit).Offset(offset).Select("send_user_id", "count(distinct rev_user_id) as count").Scan(&_list)

	var userMessageMap = map[uint]int{}

	for _, r := range _list {
		userMessageMap[r.SendUserID] = r.Count
	}
	var userIDList []uint
	for uid, _ := range userMessageMap {
		userIDList = append(userIDList, uid)
	}
	var userList []models.UserModel
	global.DB.Find(&userList, userIDList)
	var userMap = map[uint]models.UserModel{}
	for _, model := range userList {
		userMap[model.ID] = model
	}

	var list = make([]MessageUserListResponse, 0)
	for uid, count := range userMessageMap {
		user := userMap[uid]
		list = append(list, MessageUserListResponse{
			UserName: user.Username,
			NickName: user.Nickname,
			UserID:   user.ID,
			Avatar:   user.Avatar,
			Count:    count,
		})
	}
	res.OkWithList(c, list, int(count))
}

// MessageUserListByMeView 我与其他用户的聊天列表
// @Tags 消息管理
// @Summary 我与其他用户的聊天列表
// @Description 我与其他用户的聊天列表
// @Router /api/message_users/me [get]
// @Param token header string  true  "token"
// @Produce json
// @Success 200 {object} res.Response{data=res.DataListResponse[list = MessageUserListResponse[]]}
func (m MessageApi) MessageUserListByMeView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claims := _claim.(*jwts.MyClaims)
	c.Request.URL.RawQuery = fmt.Sprintf("userID=%d", claims.UserID)
	m.MessageUserListByUserView(c)
}

type MessageUserListByUserRequest struct {
	UserID uint `json:"userID" form:"userID" binding:"required"`
}

// MessageUserListByUserView 某个用户的聊天列表
// @Tags 消息管理
// @Summary 某个用户的聊天列表
// @Description 某个用户的聊天列表
// @Router /api/message_users/user [get]
// @Param token header string  true  "token"
// @Param data query MessageUserListByUserRequest   false  "查询参数"
// @Produce json
// @Success 200 {object} res.Response{data=res.DataListResponse[list = MessageUserListResponse[]]}
func (MessageApi) MessageUserListByUserView(c *gin.Context) {
	var cr MessageUserListByUserRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithMsg(c, "参数错误")
		return
	}

	type resType struct {
		SendUserID uint
		RevUserID  uint
		Count      int
	}

	var _list []resType
	global.DB.Model(models.MessageModel{}).Where("send_user_id = ? or rev_user_id = ?", cr.UserID, cr.UserID).
		Group("send_user_id").
		Group("rev_user_id").Select("send_user_id", "rev_user_id", "count(id) as count").Scan(&_list)

	var userMessageMap = map[uint]int{}

	for _, r := range _list {
		sendVal, ok1 := userMessageMap[r.SendUserID]
		if !ok1 && cr.UserID != r.SendUserID {
			userMessageMap[r.SendUserID] = r.Count
		} else {
			if cr.UserID != r.SendUserID {
				userMessageMap[r.SendUserID] = r.Count + sendVal
			}
		}
		revVal, ok2 := userMessageMap[r.RevUserID]
		if !ok2 && cr.UserID != r.RevUserID {
			userMessageMap[r.RevUserID] = r.Count
		} else {
			if cr.UserID != r.RevUserID {
				userMessageMap[r.RevUserID] = r.Count + revVal
			}
		}
	}
	var userIDList []uint
	for uid, _ := range userMessageMap {
		userIDList = append(userIDList, uid)
	}
	var userList []models.UserModel
	global.DB.Find(&userList, userIDList)
	var userMap = map[uint]models.UserModel{}
	for _, model := range userList {
		userMap[model.ID] = model
	}

	var list = make([]MessageUserListResponse, 0)
	for uid, count := range userMessageMap {
		user := userMap[uid]
		list = append(list, MessageUserListResponse{
			UserName: user.Username,
			NickName: user.Nickname,
			UserID:   user.ID,
			Avatar:   user.Avatar,
			Count:    count,
		})
	}

	res.OkWithList(c, list, len(list))
}

type MessageUserRecordRequest struct {
	common.PageInfo
	SendUserID uint `json:"sendUserID" form:"sendUserID" binding:"required"`
	RevUserID  uint `json:"revUserID" form:"revUserID" binding:"required"`
}

// MessageUserRecordView 两个用户之间的聊天记录
// @Tags 消息管理
// @Summary 两个用户之间的聊天记录
// @Description 两个用户之间的聊天记录
// @Router /api/message_users/record [get]
// @Param token header string  true  "token"
// @Param data query MessageUserRecordRequest   false  "查询参数"
// @Produce json
// @Success 200 {object} res.Response{data=res.DataListResponse[list = models.MessageModel[]]}
func (MessageApi) MessageUserRecordView(c *gin.Context) {
	var cr MessageUserRecordRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithMsg(c, "参数错误")
		return
	}

	list, count, _ := common.ListQuery(models.MessageModel{}, common.Options{
		PageInfo: cr.PageInfo,
		Where:    global.DB.Where("(send_user_id = ? and rev_user_id = ? ) or ( rev_user_id = ? and send_user_id = ? )", cr.SendUserID, cr.RevUserID, cr.SendUserID, cr.RevUserID),
	})

	res.OkWithList(c, list, count)
}
