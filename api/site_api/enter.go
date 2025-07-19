package site_api

import (
	"Blog_server/common/res"
	"fmt"
	"github.com/gin-gonic/gin"
)

type SiteApi struct {
}

type SiteUpdateRequest struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required" label:"年龄"`
}

func (SiteApi) SiteInfoView(c *gin.Context) {
	fmt.Println("site_info_view")
	res.OkWithData(c, "xx")
	return
}

func (SiteApi) SiteUpdateView(c *gin.Context) {
	//先走的中间件  我这里获得的 就是中间件里的log
	//log := log_service.GetLog(c)
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
