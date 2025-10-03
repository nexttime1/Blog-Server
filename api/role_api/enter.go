package role_api

import (
	"Blog_server/common/res"
	"github.com/gin-gonic/gin"
)

type RoleApi struct {
}

type OptionResponse struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

func (RoleApi) RoleIdListView(c *gin.Context) {
	res.OkWithData(c, []OptionResponse{
		{
			Label: "管理员",
			Value: 1,
		},
		{
			Label: "普通用户",
			Value: 2,
		},
		{
			Label: "游客",
			Value: 3,
		},
	})

}
