package site_api

import (
	"Blog_server/models/enum"
	"Blog_server/service/log_service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

type SiteApi struct {
}

type SiteUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}

func (SiteApi) SiteInfoView(c *gin.Context) {
	fmt.Println("site_info_view")
	log_service.NewLoginSuccess(c, enum.UserPwdLoginType)
	log_service.NewLoginFail(c, enum.UserPwdLoginType, "用户不存在", "xtm", "1234")
	c.JSON(200, gin.H{"code": 200, "msg": "站点信息"})

	return
}

func (SiteApi) SiteUpdateView(c *gin.Context) {
	log := log_service.GetLog(c)
	fmt.Println("SiteUpdateView")
	log.ShowRequest()
	log.ShowRequestHeader()
	log.ShowResponse()
	log.ShowResponseHeader()
	log.SetTitle("更新")
	log.SetItemInfo("请求时间", time.Now())
	log.SetImage("https://www.baidu.com")
	log.SetLink("学习地址", "https://www.bilibili.com/")
	c.Header("xxx", "xtm123")

	var req SiteUpdateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logrus.Errorf(err.Error())
		log.SetError("参数绑定失败", err)

	}
	fmt.Println("req", req)

	log.SetItemInfo("请求结构体", req)
	log.SetItemInfo("请求切片", []string{"1", "2", "3"})
	log.SetItemInfo("字符串", "你好")
	log.SetItemInfo("数字", 123)
	//log.Save()                                     //先调用
	c.JSON(200, gin.H{"code": 200, "msg": "站点信息"}) //调用 c.JSON() 方法时，它最终会自动调用你自定义的 ResponseWriter 的Write方法
	return
}
