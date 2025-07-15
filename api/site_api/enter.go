package site_api

import (
	"Blog_server/models/enum"
	"Blog_server/service/log_service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SiteApi struct {
}

type SiteUpdateRequest struct {
	Name string `json:"name"`
}

func (SiteApi) SiteInfoView(c *gin.Context) {
	fmt.Println("site_info_view")
	log_service.NewLoginSuccess(c, enum.UserPwdLoginType)
	log_service.NewLoginFail(c, enum.UserPwdLoginType, "用户不存在", "xtm", "1234")
	c.JSON(200, gin.H{"code": 200, "msg": "站点信息"})

	return
}

func (SiteApi) SiteUpdateView(c *gin.Context) {
	fmt.Println("SiteUpdateView")

	var req SiteUpdateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logrus.Errorf(err.Error())
	}
	fmt.Println("req", req)

	log.Save()

	c.JSON(200, gin.H{"code": 200, "msg": "站点信息"})
	return
}
