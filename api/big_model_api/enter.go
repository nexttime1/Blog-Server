package big_model_api

import (
	"Blog_server/common"
	"Blog_server/common/res"
	"Blog_server/conf"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"Blog_server/service/big_model_service"
	"Blog_server/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

type BigModelApi struct {
}

func (BigModelApi) BigModelOptionListView(c *gin.Context) {
	_, exist := c.Get("claims")
	if !exist {
		// 有问题 AuthMiddleware 已经res 返回了 这里不需要返回
		return
	}
	res.OkWithData(c, global.Config.BigModel)

}

type SettingType struct {
	conf.Setting
	Help string `json:"help"`
}

const FilePath = "uploads/doc"

var xtm = "big_model"

//BigModelSettingView  只有管理员可以看到  而 游客和 普通用户就看一点

func (BigModelApi) BigModelSettingView(c *gin.Context) {
	flag := false
	MdPath := path.Join(FilePath, fmt.Sprintf("%s.md", xtm))
	byteData, err := os.ReadFile(MdPath)
	if err != nil {
		logrus.Errorf("read file %s error %v", MdPath, err)
		return
	}
	_claims, exist := c.Get("claims")
	if !exist {
		// 游客

	} else {
		Claims := _claims.(*jwts.MyClaims)
		if Claims.Role == enum.AdminRole {
			flag = true
		}
	}
	if flag {
		var st = SettingType{
			Setting: global.Config.BigModel.Setting,
			Help:    string(byteData),
		}

		res.OkWithData(c, st)
		return
	}

	UserSetting := conf.Setting{
		Name:      "",
		Enable:    true,
		ApiKey:    "",
		ApiSecret: "",
		Title:     "",
		Prompt:    "",
	}

	var UserSt = SettingType{
		Setting: UserSetting,
		Help:    string(byteData),
	}

	res.OkWithData(c, UserSt)

}

func (BigModelApi) BigModelUpdateView(c *gin.Context) {
	_, exist := c.Get("claims")
	if !exist {

		return
	}
	var setting conf.Setting

	err := c.ShouldBindJSON(&setting)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	global.Config.BigModel.Setting = setting

	err = common.ToYAML(global.SettingYaml, global.Config)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.FailWithMsg(c, "修改成功")

}

func (BigModelApi) BigModelSessionView(c *gin.Context) {
	_, exist := c.Get("claims")
	if !exist {
		return
	}
	res.OkWithData(c, global.Config.BigModel.SessionSetting)
}

func (BigModelApi) BigModelSessionUpdateView(c *gin.Context) {
	_, exist := c.Get("claims")
	if !exist {
		return
	}
	var cr conf.SessionSetting
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	global.Config.BigModel.SessionSetting = cr
	err = common.ToYAML(global.SettingYaml, global.Config)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.FailWithMsg(c, "修改成功")
}

func (BigModelApi) UserScopeEnableView(c *gin.Context) {
	_Claims, exist := c.Get("claims")
	if !exist {
		return
	}
	Claims := _Claims.(*jwts.MyClaims)

	response, err := big_model_service.UserScopeEnableService(Claims)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}

	res.OkWithData(c, response)

}

func (BigModelApi) UserScopeView(c *gin.Context) {
	_Claims, exist := c.Get("claims")
	if !exist {
		return
	}
	Claims := _Claims.(*jwts.MyClaims)

	var cr big_model_service.UserScopeRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	err = big_model_service.UserScopeService(cr, Claims)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	res.OkWithData(c, "积分领取成功")

}

func (BigModelApi) AutoReplyUpdateView(c *gin.Context) {
	_, exist := c.Get("claims")
	if !exist {
		return
	}
	var cr big_model_service.AutoReplyUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	flag, err := big_model_service.AutoReplyUpdateService(cr)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	if flag == 1 {
		res.OkWithMessage(c, "自动回复添加成功")
	} else {
		res.OkWithMessage(c, "自动回复更新成功")
	}

}

func (BigModelApi) AutoReplyListView(c *gin.Context) {
	_, exist := c.Get("claims")
	if !exist {
		return
	}
	var cr common.PageInfo
	c.ShouldBindQuery(&cr) // 有默认值 没事
	list, count, err := common.ListQuery(models.AutoReplyModel{}, common.Options{
		PageInfo: cr,
		Likes:    []string{"name"},
	})
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithList(c, list, count)

}
