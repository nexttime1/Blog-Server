package log_api

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"github.com/gin-gonic/gin"
)

type LogApi struct {
}

type LogListView struct {
	Limit       int            `form:"limit"`
	Page        int            `form:"page"`
	Key         string         `form:"key"`
	LogType     enum.LogType   `form:"logType"` //日志类型 1 2 3
	Level       enum.LevelType `form:"level"`   //日志级别  1 2 3
	UserID      uint           `form:"userID"`  //用户id   可以没有  没登录  设置为0
	IP          string         `form:"ip"`
	LoginStatus bool           `form:"loginStatus"` //登录状态
	ServiceName string         `form:"serviceName"`
}

func (LogApi) LogListNew(c *gin.Context) {
	var cr LogListView
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	var list []models.LogModel
	//添加限制  增加默认
	if cr.Page <= 0 || cr.Page >= 20 {
		cr.Page = 1
	}
	if cr.Limit <= 0 || cr.Limit >= 50 {
		cr.Limit = 10
	}

	offset := (cr.Page - 1) * cr.Limit
	model := models.LogModel{ //前端没赋值  就相当于没用  Where  就是显示全部
		LogType:     cr.LogType,
		Level:       cr.Level,
		UserID:      cr.UserID,
		IP:          cr.IP,
		LoginStatus: cr.LoginStatus,
		ServiceName: cr.ServiceName,
	}

	global.DB.Debug().Where(model).Limit(cr.Limit).Offset(offset).Find(&list)
	// 总数
	var count int64
	global.DB.Debug().Where(model).Model(&models.LogModel{}).Count(&count)

	res.OkWithList(c, list, count)
}
