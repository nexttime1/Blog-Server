package log_api

import (
	"Blog_server/common"
	"Blog_server/common/res"
	"Blog_server/models"
	"Blog_server/models/enum"
	"github.com/gin-gonic/gin"
)

type LogApi struct {
}

type LogListView struct {
	common.PageInfo
	LogType     enum.LogType   `form:"logType"` //日志类型 1 2 3
	Level       enum.LevelType `form:"level"`   //日志级别  1 2 3
	UserID      uint           `form:"userID"`  //用户id   可以没有  没登录  设置为0
	IP          string         `form:"ip"`
	LoginStatus bool           `form:"loginStatus"` //登录状态
	ServiceName string         `form:"serviceName"`
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

	list, count, err := common.ListQuery(models.LogModel{ //前端没赋值  就相当于没用  Where  就是显示全部
		LogType:     cr.LogType,
		Level:       cr.Level,
		UserID:      cr.UserID,
		IP:          cr.IP,
		LoginStatus: cr.LoginStatus,
		ServiceName: cr.ServiceName,
	}, common.Options{
		PageInfo:     cr.PageInfo,
		Likes:        []string{"title"},
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
