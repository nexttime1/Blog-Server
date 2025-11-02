package settings_api

import (
	"Blog_server/common/res"
	"Blog_server/conf"
	"Blog_server/global"
	"Blog_server/service/log_service"
	"Blog_server/utils/jwts"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
)

type SettingApi struct {
}

type SettingsResponse struct {
	Name string `uri:"name"`
}

func (SettingApi) SettingInfoView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	var s SettingsResponse
	err := c.ShouldBindUri(&s)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log := log_service.GetLog(c)
	log.SetTitle("查看系统配置")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
	if s.Name == "" {
		fmt.Println("没接收到东西")
	}
	switch s.Name {
	case "site":
		res.OkWithData(c, global.Config.SiteInfo)
	case "qq":
		qq := global.Config.QQ
		qq.Key = "******"
		res.OkWithData(c, qq)
	case "email":
		email := global.Config.Email
		email.Password = "******"
		res.OkWithData(c, email)
	case "qiniu":
		qiniu := global.Config.QiNiu
		qiniu.AccessKey = "******"
		res.OkWithData(c, qiniu)

	case "jwt":
		jwt := global.Config.Jwt
		jwt.Secret = "******"
		res.OkWithData(c, jwt)
	case "chat_group":
		res.OkWithData(c, global.Config.ChatGroup)
	case "gaode":
		res.OkWithData(c, global.Config.Gaode)
	default:
		res.FailWithErr(c, errors.New("配置信息未找到"))
	}
	return

}

func (SettingApi) SettingInfoUpdateView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	var s SettingsResponse
	err := c.ShouldBindUri(&s)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log := log_service.GetLog(c)
	log.SetTitle("修改系统配置")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
	switch s.Name {
	case "site":
		var siteData conf.SiteInfo
		err := c.ShouldBindJSON(&siteData)
		if err != nil {
			logrus.Error("接收参数错误", err)
			res.FailWithCode(c, res.ArgumentError)
			return
		}
		global.Config.SiteInfo = siteData
	case "email":
		var emailData conf.Email
		err := c.ShouldBindJSON(&emailData)
		if err != nil {
			logrus.Error("接收参数错误", err)
			res.FailWithCode(c, res.ArgumentError)
			return
		}
		password := global.Config.Email.Password
		global.Config.Email = emailData
		if emailData.Password == "******" {
			global.Config.Email.Password = password
		}

	case "qiniu":
		var QiNiuData conf.QiNiu
		err := c.ShouldBindJSON(&QiNiuData)
		if err != nil {
			logrus.Error("接收参数错误", err)
			res.FailWithCode(c, res.ArgumentError)
			return
		}
		AccessKey := global.Config.QiNiu.AccessKey
		global.Config.QiNiu = QiNiuData
		if QiNiuData.AccessKey == "******" {
			global.Config.QiNiu.AccessKey = AccessKey
		}

	case "jwt":
		var jwtData conf.Jwt
		err := c.ShouldBindJSON(&jwtData)
		if err != nil {
			logrus.Error("接收参数错误", err)
			res.FailWithCode(c, res.ArgumentError)
			return
		}
		secret := global.Config.Jwt.Secret
		global.Config.Jwt = jwtData
		if jwtData.Secret == "******" {
			global.Config.Jwt.Secret = secret
		}

	case "qq":
		var QQData conf.QQ
		err := c.ShouldBindJSON(&QQData)
		if err != nil {
			logrus.Error("接收参数错误", err)
			res.FailWithCode(c, res.ArgumentError)
			return
		}
		Key := global.Config.QQ.Key
		global.Config.QQ = QQData
		if QQData.Key == "******" {
			global.Config.QQ.Key = Key
		}
	case "chat_group":
		var info conf.ChatGroup
		err = c.ShouldBindJSON(&info)
		if err != nil {
			res.FailWithCode(c, res.ArgumentError)
			return
		}
		global.Config.ChatGroup = info
	case "gaode":
		var info conf.Gaode
		err = c.ShouldBindJSON(&info)
		if err != nil {
			res.FailWithCode(c, res.ArgumentError)
			return
		}
		global.Config.Gaode = info
	default:
		res.FailWithErr(c, errors.New("配置信息未找到"))
	}

	// 3. 把更新后的 global.Config 完整写入 yaml 文件（永久保存）
	filePath := global.SettingYaml
	if err := saveConfigToYaml(filePath, global.Config); err != nil { // 假设 ConfigPath 是 yaml 路径
		res.FailWithErr(c, errors.New("配置保存失败："+err.Error()))
		return
	}

	res.OkWithMessage(c, "配置更新成功")
}

// 保存配置到 yaml 的工具函数（确保序列化正确）
func saveConfigToYaml(filePath string, config *conf.Config) error {
	// 用 yaml 包序列化（注意导入正确的包，如 gopkg.in/yaml.v3）
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败：%w", err)
	}
	// 写入文件（权限 0644，确保目录可写）
	if err := os.WriteFile(filePath, yamlData, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败：%w", err)
	}
	return nil

}
