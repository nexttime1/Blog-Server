package site_api

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/middleware"
	"Blog_server/service/log_service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type SiteApi struct {
}

type SiteUpdateRequest struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required" label:"年龄"`
}

type SiteInfoRequest struct {
	Name string `uri:"name"`
}

func (SiteApi) SiteInfoView(c *gin.Context) {
	fmt.Println("site_info_view")
	var cr SiteInfoRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	//就这一个不需要管理员身份 所以单独拿出来
	if cr.Name == "site" {
		res.OkWithData(c, global.Config.Site)
		return
	}
	middleware.AdminMiddleware(c)
	_, exists := c.Get("claims")
	if !exists {
		return
	}

	var data any
	switch cr.Name {
	case "Email":
		data = global.Config.Email
	case "QQ":
		data = global.Config.QQ
	case "QiNiu":
		data = global.Config.QiNiu
	case "Ai":
		data = global.Config.Ai
	default:
		res.FailWithErr(c, errors.New("未找到配置"))
		return
	}

	res.OkWithData(c, data)
	return
}

func (SiteApi) SiteUpdateView(c *gin.Context) {
	//先走的中间件  我这里获得的 就是中间件里的log
	log := log_service.GetLog(c)
	log.ShowRequest()
	log.ShowResponse()
	fmt.Println("SiteUpdateView")

	var req SiteUpdateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	fmt.Println("req", req)

	res.OkWithMessage(c, "更新成功")
	//c.JSON(200, gin.H{"code": 200, "msg": "站点信息"}) //调用 c.JSON() 方法时，它最终会自动调用你自定义的 ResponseWriter 的Write方法
	return
}
