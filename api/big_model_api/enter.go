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

type SettingType struct {
	conf.Setting
	Help string `json:"help"`
}

const FilePath = "uploads/doc"

var xtm = "big_model"

// BigModelOptionListView
// @Summary 获取可用大模型配置
// @Description 仅管理员可获取系统中可用的大模型配置列表
// @Tags 大模型管理
// @Accept application/json
// @Produce application/json
// @Security BearerAuth
// @Success 200 {object} res.Response{data=[]conf.Setting} "成功返回可用大模型配置列表"
// @Failure 401 {object} res.Response "未登录或非管理员权限"
// @Failure 500 {object} res.Response "服务内部错误"
// @Router /big_model/usable [get]
func (BigModelApi) BigModelOptionListView(c *gin.Context) {
	_, exist := c.Get("claims")
	if !exist {
		// 有问题 AuthMiddleware 已经res 返回了 这里不需要返回
		return
	}
	res.OkWithData(c, global.Config.BigModel)

}

// BigModelSettingView  只有管理员可以看到  而 游客和 普通用户就看一点
// @Summary 获取大模型配置
// @Description 登录用户可获取大模型配置（管理员可见完整配置，普通用户隐藏敏感字段）
// @Tags 大模型管理
// @Accept application/json
// @Produce application/json
// @Security BearerAuth
// @Success 200 {object} res.Response{data=SettingType} "成功返回配置及帮助信息"
// @Failure 401 {object} res.Response "未登录"
// @Failure 500 {object} res.Response "服务内部错误（如文件读取失败）"
// @Router /big_model/setting [get]
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

// BigModelUpdateView
// @Summary 修改大模型配置
// @Description 仅管理员可修改大模型基本配置（修改后会持久化到配置文件）
// @Tags 大模型管理
// @Accept application/json
// @Produce application/json
// @Security BearerAuth
// @Param body body conf.Setting true "大模型配置信息"
// @Success 200 {object} res.Response{msg=string} "修改成功"
// @Failure 400 {object} res.Response "请求参数格式错误"
// @Failure 401 {object} res.Response "未登录或非管理员权限"
// @Failure 500 {object} res.Response "配置文件写入失败"
// @Router /big_model/setting [put]
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

// BigModelSessionView
// @Summary 获取大模型会话配置
// @Description 登录用户可获取大模型会话相关配置（如最大长度、自动保存等）
// @Tags 大模型管理
// @Accept application/json
// @Produce application/json
// @Security BearerAuth
// @Success 200 {object} res.Response{data=conf.SessionSetting} "成功返回会话配置"
// @Failure 401 {object} res.Response "未登录"
// @Router /big_model/session_setting [get]
func (BigModelApi) BigModelSessionView(c *gin.Context) {
	_, exist := c.Get("claims")
	if !exist {
		return
	}
	res.OkWithData(c, global.Config.BigModel.SessionSetting)
}

// BigModelSessionUpdateView
// @Summary 修改大模型会话配置
// @Description 仅管理员可修改大模型会话配置（修改后会持久化到配置文件）
// @Tags 大模型管理
// @Accept application/json
// @Produce application/json
// @Security BearerAuth
// @Param body body conf.SessionSetting true "会话配置信息"
// @Success 200 {object} res.Response{msg=string} "修改成功"
// @Failure 400 {object} res.Response "请求参数格式错误"
// @Failure 401 {object} res.Response "未登录或非管理员权限"
// @Failure 500 {object} res.Response "配置文件写入失败"
// @Router /big_model/session_setting [put]
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

// UserScopeEnableView
// @Summary 查询用户是否可领取积分
// @Description 登录用户查询自己是否满足积分领取条件
// @Tags 用户积分
// @Accept application/json
// @Produce application/json
// @Security BearerAuth
// @Success 200 {object} res.Response{data=bool} "成功返回是否可领取（true=可领取，false=不可领取）"
// @Failure 401 {object} res.Response "未登录"
// @Failure 500 {object} res.Response "服务内部错误"
// @Router /big_model/user_scope_enable [get]
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

// UserScopeView
// @Summary 用户领取积分
// @Description 登录用户提交积分领取请求（需满足领取条件）
// @Tags 用户积分
// @Accept application/json
// @Produce application/json
// @Security BearerAuth
// @Param body body big_model_service.UserScopeRequest true "积分领取请求参数"
// @Success 200 {object} res.Response{data=string} "积分领取成功"
// @Failure 400 {object} res.Response "请求参数格式错误"
// @Failure 401 {object} res.Response "未登录"
// @Failure 500 {object} res.Response "领取失败（如不满足条件）"
// @Router /big_model/user_scope [post]
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

// AutoReplyUpdateView
// @Summary 新增或更新自动回复
// @Description 仅管理员可操作，ID存在则更新，不存在则新增
// @Tags 自动回复管理
// @Accept application/json
// @Produce application/json
// @Security BearerAuth
// @Param body body big_model_service.AutoReplyUpdateRequest true "自动回复信息"
// @Success 200 {object} res.Response{msg=string} "添加成功或更新成功"
// @Failure 400 {object} res.Response "请求参数格式错误"
// @Failure 401 {object} res.Response "未登录或非管理员权限"
// @Failure 500 {object} res.Response "操作失败"
// @Router /big_model/auto_reply [put]
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

// AutoReplyListView
// @Summary 获取自动回复列表
// @Description 仅管理员可获取，支持分页和名称模糊查询
// @Tags 自动回复管理
// @Accept application/json
// @Produce application/json
// @Security BearerAuth
// @Param page query int false "页码（默认1）"
// @Param page_size query int false "每页条数（默认10）"
// @Param name query string false "名称模糊查询"
// @Success 200 {object} res.Response{data=res.DataListResponse{List=[]models.AutoReplyModel,Count=int}} "成功返回列表和总条数"
// @Failure 401 {object} res.Response "未登录或非管理员权限"
// @Failure 500 {object} res.Response "查询失败"
// @Router /big_model/auto_reply [get]
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

// AutoReplyDeleteView
// @Summary 批量删除自动回复
// @Description 仅管理员可操作，接收要删除的自动回复ID列表，校验ID存在性后执行删除，返回删除数量
// @Tags 自动回复管理
// @Accept application/json
// @Produce application/json
// @Security BearerAuth
// @Param body body models.RemoveRequest true "删除请求参数（包含要删除的ID列表）"
// @Success 200 {object} res.Response{msg=string} "成功返回删除数量（如：共删除2个自动回复）"
// @Failure 400 {object} res.Response "请求参数错误（如ID列表为空或格式错误）"
// @Failure 401 {object} res.Response "未登录或非管理员权限"
// @Failure 404 {object} res.Response "部分或全部ID不存在（提示：删除的Id 不存在）"
// @Failure 500 {object} res.Response "数据库操作失败"
// @Router /big_model/auto_reply [delete]
func (BigModelApi) AutoReplyDeleteView(c *gin.Context) {
	_, exist := c.Get("claims")
	if !exist {
		return
	}
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	var AutoModels []models.AutoReplyModel

	err = global.DB.Take(&AutoModels, "id in ?", cr.IDList).Error
	if err != nil {
		logrus.Errorf("%v", err.Error())
		res.FailWithMsg(c, "删除的Id 不存在")
		return
	}
	global.DB.Delete(&AutoModels)

	res.FailWithMsg(c, fmt.Sprintf("共删除%d个自动恢复", len(AutoModels)))

}

// BigModelTagUpdateView
// @Summary 新增或修改大模型标签
// @Description 仅管理员可操作，根据是否传入ID判断：无ID则新增标签，有ID则修改已有标签
// @Tags 大模型标签管理
// @Accept application/json
// @Produce application/json
// @Security BearerAuth
// @Param body body big_model_service.TagUpdateRequest true "标签新增/修改参数"
// @Success 200 {object} res.Response{data=string} "操作成功（新增返回：标签添加成功；修改返回：标签修改成功）"
// @Failure 400 {object} res.Response "请求参数错误（如名称为空、长度超限等）"
// @Failure 401 {object} res.Response "未登录或非管理员权限"
// @Failure 500 {object} res.Response "服务内部错误（如数据库操作失败、标签名称重复等）"
// @Router /big_model/tags [put]
func (BigModelApi) BigModelTagUpdateView(c *gin.Context) {
	_, exist := c.Get("claims")
	if !exist {
		return
	}
	var cr big_model_service.TagUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	flag, err := big_model_service.BigModelTagUpdateService(cr)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	if flag == 1 {
		// 添加
		res.OkWithData(c, "标签添加成功")
		return
	}
	res.FailWithMsg(c, "标签修改成功")

}

// BigModelTagListView
// @Summary 获取大模型标签分页列表
// @Description 仅管理员可获取大模型标签的分页列表，支持基础分页参数
// @Tags 大模型标签管理
// @Accept application/json
// @Produce application/json
// @Security BearerAuth
// @Param page query int false "页码（默认值：1）" minimum(1)
// @Param page_size query int false "每页条数（默认值：10，最大值：50）" minimum(1) maximum(50)
// @Success 200 {object} res.Response{data=res.DataListResponse{List=[]big_model_service.TagListResponse,Count=int}} "成功返回标签列表及总条数"
// @Failure 400 {object} res.Response "分页参数错误（如页码小于1、每页条数超出限制等）"
// @Failure 401 {object} res.Response "未登录或非管理员权限"
// @Failure 500 {object} res.Response "服务内部错误（如数据库查询失败）"
// @Router /big_model/tags [get]
func (BigModelApi) BigModelTagListView(c *gin.Context) {
	_, exist := c.Get("claims")
	if !exist {
		return
	}
	var cr common.PageInfo
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	tagListResponseList, count, err := big_model_service.BigModelTagListService(cr)
	if err != nil {
		res.FailWithErr(c, err)
		return

	}
	res.OkWithList(c, tagListResponseList, count)

}

func (BigModelApi) BigModelTagRemoveView(c *gin.Context) {
	_, exist := c.Get("claims")
	if !exist {
		return
	}
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	err = big_model_service.BigModelTagRemoveService(cr)

	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	res.OkWithMessage(c, "删除成功")
}

func (BigModelApi) BigModelRoleUpdateView(c *gin.Context) {
	_, exist := c.Get("claims")
	if !exist {
		return
	}
	var cr big_model_service.RoleUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	flag, err := big_model_service.BigModelRoleUpdateService(cr)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	if flag == 1 {
		res.OkWithData(c, "角色添加成功")
		return
	}
	res.FailWithMsg(c, "角色更新成功")
}

func (BigModelApi) BigModelRoleListiew(c *gin.Context) {
	_, exist := c.Get("claims")
	if !exist {
		return
	}
	var cr common.PageInfo
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	list, count, err := common.ListQuery(models.BigModelRoleModel{}, common.Options{
		PageInfo: cr,
		Likes:    []string{"name"},
	})
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithList(c, list, count)
}

func (BigModelApi) BigModelRoleRemoveView(c *gin.Context) {
	_, exist := c.Get("claims")
	if !exist {
		return
	}
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
	}
	err, msg := big_model_service.BigModelRoleRemoveService(cr)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	res.OkWithMessage(c, msg)
}
