package user_api

import (
	"Blog_server/common/res"
	"Blog_server/service/user_service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserApi struct {
}

func (u UserApi) UserEmailLogin(c *gin.Context) {
	var request user_service.EmailLoginRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	token, msg, err := user_service.UserEmailLoginService(request)
	if err != nil {
		logrus.Errorf("UserEmailLogin %s  err:  %s", msg, err)
		res.FailWithMsg(c, msg)
		return
	}
	res.OkWithData(c, token)
}
