package settings_api

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type SettingApi struct {
}

type SettingsResponse struct {
	Name string `uri:"name"`
}

func (SettingApi) SettingInfoView(c *gin.Context) {
	fmt.Println("SettingInfoView")
	var s SettingsResponse
	err := c.ShouldBindUri(&s)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	fmt.Println("打印打印", s.Name)
	if s.Name == "" {
		fmt.Println("没接收到东西")
	}
	switch s.Name {
	case "site":
		res.OkWithData(c, global.Config.SiteInfo)
	case "qq":
		res.OkWithData(c, global.Config.QQ)
	case "email":
		res.OkWithData(c, global.Config.Email)
	case "qiniu":
		res.OkWithData(c, global.Config.QiNiu)
	default:
		res.FailWithErr(c, errors.New("配置信息未找到"))
	}
	return

}
