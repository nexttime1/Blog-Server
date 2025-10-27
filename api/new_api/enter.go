package new_api

import (
	"Blog_server/common/res"
	"Blog_server/service/new_service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

type NewApi struct {
}

const newAPI = "https://api.codelife.cc/api/top/list"
const timeout = 2 * time.Second

func (NewApi) NewListView(c *gin.Context) {
	var cr new_service.Params
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		logrus.Errorf("接收参数问题：%s", err)
		res.FailWithErr(c, err)
		return
	}

	data, err := new_service.NewListService(cr, newAPI, timeout)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithData(c, data)

}
