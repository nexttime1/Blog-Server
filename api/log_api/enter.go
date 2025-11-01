package log_api

import (
	"Blog_server/common"
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"Blog_server/service/log_service"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type LogApi struct {
}

type LogListView struct {
	common.PageInfo
	LogType     enum.LogType   `form:"type"`   //日志类型 1 2 3
	Level       enum.LevelType `form:"level"`  //日志级别  1 2 3
	UserID      uint           `form:"userID"` //用户id   可以没有  没登录  设置为0
	IP          string         `form:"ip"`
	Addr        string         `form:"addr"`
	Date        string         `form:"date"`
	LoginStatus *bool          `form:"status"`                   //登录状态
	UserName    string         `json:"userName" form:"userName"` // 查用户名
}

type LogListResponse struct {
	models.LogModel
	UserNickName string `form:"userNickname"`
	UserAvatar   string `form:"userAvatar"`
}

func (LogApi) LogListNew(c *gin.Context) {
	var cr LogListView
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	var query = global.DB.Where("")
	if cr.Date != "" {
		_, dateTimeErr := time.Parse("2006-01-02", cr.Date)
		if dateTimeErr != nil {
			res.FailWithMsg(c, "时间格式错误")
			return
		}
		query.Where("date(created_at) = ?", cr.Date)
	}
	if cr.LoginStatus != nil {
		query.Where("login_status = ?", cr.LoginStatus)
	}

	list, count, err := common.ListQuery(models.LogModel{ //前端没赋值  就相当于没用  Where  就是显示全部
		LogType:  cr.LogType,
		Level:    cr.Level,
		UserID:   cr.UserID,
		IP:       cr.IP,
		UserName: cr.UserName,
		Addr:     cr.Addr,
	}, common.Options{
		PageInfo:     cr.PageInfo,
		Where:        query,
		Likes:        []string{"title", "user_name"},
		Preload:      []string{"UserModel"},
		Debug:        true,
		DefaultOrder: "created_at DESC", //写死了  默认降序排序  前端修改的话
	})

	var _list = make([]LogListResponse, 0)
	for _, v := range list {
		_list = append(_list, LogListResponse{
			LogModel:     v,
			UserNickName: v.UserModel.Nickname,
			UserAvatar:   v.UserModel.Avatar,
		})
	}

	res.OkWithList(c, _list, count)
}

func (LogApi) LogReadView(c *gin.Context) {
	var cr models.IDRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	var log models.LogModel
	err = global.DB.Take(&log, cr.ID).Error
	if err != nil {
		res.FailWithMsg(c, "不存在的日志")
		return
	}
	// 判断是否已经读取
	if !log.IsRead {
		global.DB.Debug().Model(&log).Update("is_read", true)
	}
	res.OkWithMessage(c, "读取成功")

}

func (LogApi) LogRemoveView(c *gin.Context) {
	var rc models.RemoveRequest
	err := c.ShouldBindJSON(&rc)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	//删除我 需要保存操作日志   先走中间件  中间件 设置了Set  Log  所以我这里可以获得 log  可以激活Save函数
	log := log_service.GetLog(c)
	log.ShowRequest()
	log.ShowResponse()

	var ModelList []models.LogModel

	ModelList = common.BatchRemove(ModelList, rc)

	msg := fmt.Sprintf("日志删除成功，共删除%d条数据", len(ModelList))
	res.OkWithMessage(c, msg)
}
