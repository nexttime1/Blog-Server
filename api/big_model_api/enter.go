package big_model_api

import (
	"Blog_server/common"
	"Blog_server/common/res"
	"Blog_server/conf"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"Blog_server/service/big_model_service"
	"Blog_server/service/log_service"
	"Blog_server/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
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
// @Tags 大模型配置管理
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
	res.OkWithData(c, global.Config.BigModel.ModelList)

}

// BigModelSettingView  只有管理员可以看到  而 游客和 普通用户就看一点
// @Summary 获取大模型配置
// @Description 登录用户可获取大模型配置（管理员可见完整配置，普通用户隐藏敏感字段）
// @Tags 大模型配置管理
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
		Slogan:    "",
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
// @Tags 大模型配置管理
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
// @Tags 大模型配置管理
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
// @Tags 大模型配置管理
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
	_claim, exist := c.Get("claims")
	if !exist {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	var cr conf.SessionSetting
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log := log_service.GetLog(c)
	log.SetTitle("修改大模型会话配置")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
	global.Config.BigModel.SessionSetting = cr
	err = common.ToYAML(global.SettingYaml, global.Config)
	if err != nil {
		log.SetItemError("修改错误", err)
		res.FailWithErr(c, err)
		return
	}
	res.OkWithMessage(c, "修改成功")
}

// UserScopeEnableView
// @Summary 查询用户是否可领取积分
// @Description 登录用户查询自己是否满足积分领取条件
// @Tags 大模型用户积分管理
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
// @Tags 大模型用户积分管理
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
	log := log_service.GetLog(c)
	log.SetTitle("用户领取积分")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", Claims.Username)
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
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claim := _claims.(*jwts.MyClaims)
	var cr big_model_service.AutoReplyUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log := log_service.GetLog(c)
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
	flag, err := big_model_service.AutoReplyUpdateService(cr, log)
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
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claim := _claims.(*jwts.MyClaims)
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log := log_service.GetLog(c)
	log.SetTitle("删除大模型自动回复")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
	var AutoModels []models.AutoReplyModel

	err = global.DB.Take(&AutoModels, "id in ?", cr.IDList).Error
	if err != nil {
		logrus.Errorf("%v", err.Error())
		res.FailWithMsg(c, "删除的Id 不存在")
		return
	}
	global.DB.Delete(&AutoModels)

	res.OkWithMessage(c, fmt.Sprintf("共删除%d个自动恢复", len(AutoModels)))

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
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claim := _claims.(*jwts.MyClaims)
	var cr big_model_service.TagUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log := log_service.GetLog(c)

	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
	flag, err := big_model_service.BigModelTagUpdateService(cr, log)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	if flag == 1 {
		// 添加
		res.OkWithData(c, "标签添加成功")
		return
	}
	res.OkWithMessage(c, "标签修改成功")

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

// BigModelTagRemoveView 批量删除大模型标签
// @Summary 批量删除大模型标签
// @Description 仅管理员可调用，批量删除指定ID的大模型标签
// @Tags 大模型标签管理
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param data body models.RemoveRequest true "删除参数" remark "包含需要删除的标签ID列表（通常为`IDs []uint`）"
// @Success 200 {object} res.Response{code=int,data=map[string]interface{},msg=string} "删除成功"
// @Failure 1001 {object} res.Response{code=int,data=map[string]interface{},msg=string} "参数错误或删除失败"
// @Router /big_model/tags [delete]
func (BigModelApi) BigModelTagRemoveView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claim := _claims.(*jwts.MyClaims)
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log := log_service.GetLog(c)
	log.SetTitle("删除大模型标签")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
	err = big_model_service.BigModelTagRemoveService(cr, log)

	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	res.OkWithMessage(c, "删除成功")
}

// BigModelRoleUpdateView 新增或更新大模型角色
// @Summary 新增或更新大模型角色
// @Description 仅管理员可调用，ID为空时新增角色，ID存在时更新角色信息
// @Tags 大模型角色管理
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param data body big_model_service.RoleUpdateRequest true "角色信息"
// @Param data.name formData string true "角色名称"
// @Param data.enable formData bool true "是否启用"
// @Param data.icon formData string false "角色图标（可选，支持系统默认或上传）"
// @Param data.abstract formData string true "角色简介"
// @Param data.scope formData int true "消耗积分"
// @Param data.prologue formData string true "开场白"
// @Param data.prompt formData string true "设定词"
// @Param data.autoReply formData bool true "是否自动回复"
// @Param data.tagList formData []uint true "关联标签ID列表"
// @Success 200 {object} res.Response{code=int,data=string,msg=string} "成功（返回`角色添加成功`或`角色更新成功`）"
// @Failure 1001 {object} res.Response{code=int,data=map[string]interface{},msg=string} "参数错误或操作失败"
// @Router /big_model/roles [put]
func (BigModelApi) BigModelRoleUpdateView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claim := _claims.(*jwts.MyClaims)
	var cr big_model_service.RoleUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log := log_service.GetLog(c)
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
	flag, err := big_model_service.BigModelRoleUpdateService(cr, log)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	if flag == 1 {
		res.OkWithData(c, "角色添加成功")
		return
	}
	res.OkWithMessage(c, "角色更新成功")
}

// BigModelRoleListView 获取大模型角色分页列表
// @Summary 获取大模型角色分页列表
// @Description 仅管理员可调用，分页查询大模型角色列表，支持按名称模糊搜索
// @Tags 大模型角色管理
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param page query int false "页码（默认1）"
// @Param pageSize query int false "每页条数（默认10）"
// @Param name query string false "角色名称（模糊搜索）"
// @Success 200 {object} res.Response{code=int,data=res.DataListResponse{List=[]models.BigModelRoleModel,Count=int},msg=string} "成功（返回角色列表和总数）"
// @Failure 1001 {object} res.Response{code=int,data=map[string]interface{},msg=string} "参数错误或查询失败"
// @Router /big_model/roles [get]
func (BigModelApi) BigModelRoleListView(c *gin.Context) {
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
		Preload:  []string{"Tags"},
	})
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithList(c, list, count)
}

// BigModelRoleRemoveView 批量删除大模型角色
// @Summary 批量删除大模型角色
// @Description 仅管理员可调用，批量删除指定ID的大模型角色
// @Tags 大模型角色管理
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param data body models.RemoveRequest true "删除参数" remark "包含需要删除的角色ID列表（格式：{\"IDs\": [1,2,3]}）"
// @Success 200 {object} res.Response{code=int,data=map[string]interface{},msg=string} "删除成功（返回具体删除结果描述）"
// @Failure 1001 {object} res.Response{code=int,data=map[string]interface{},msg=string} "参数错误（如ID列表为空）"
// @Failure 1002 {object} res.Response{code=int,data=map[string]interface{},msg=string} "服务异常（如删除失败）"
// @Router /big_model/roles [delete]
func (BigModelApi) BigModelRoleRemoveView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claim := _claims.(*jwts.MyClaims)
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
	}
	log := log_service.GetLog(c)
	log.SetTitle("删除大模型角色")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
	err, msg := big_model_service.BigModelRoleRemoveService(cr)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	res.OkWithMessage(c, msg)
}

// BigModelTagRoleListView 获取大模型标签及关联角色列表
// @Summary 获取大模型标签及关联角色列表
// @Description 无需认证，返回所有标签及其关联的角色信息（用于角色广场展示）
// @Tags 大模型角色管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} res.Response{code=int,data=[]big_model_service.TagRoleListResponse,msg=string} "成功（返回标签列表，每个标签包含关联角色信息）"
// @Failure 1002 {object} res.Response{code=int,data=map[string]interface{},msg=string} "服务异常（如查询失败）"
// @Router /big_model/square [get]
func (BigModelApi) BigModelTagRoleListView(c *gin.Context) {
	err, response := big_model_service.BigModelTagRoleListService()
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithData(c, response)

}

// BigModelSessionCreateView 创建大模型会话
// @Summary 创建大模型会话
// @Description 需用户认证，基于指定角色创建新会话（用于后续对话交互）
// @Tags 大模型会话管理
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param data body big_model_service.SessionCreateRequest true "会话创建参数"
// @Param data.RoleID body uint true "角色ID（关联的大模型角色）"
// @Param data.Name body string false "会话名称（可选，不填则自动生成）"
// @Success 200 {object} res.Response{code=int,data=uint,msg=string} "成功（返回创建的会话ID，msg为“会话创建成功”）"
// @Failure 1001 {object} res.Response{code=int,data=map[string]interface{},msg=string} "参数错误（如RoleID为空）"
// @Failure 1002 {object} res.Response{code=int,data=map[string]interface{},msg=string} "服务异常（如创建失败）"
// @Router /big_model/session [post]
func (BigModelApi) BigModelSessionCreateView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)

	var cr big_model_service.SessionCreateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	sessionID, err, flag := big_model_service.BigModelSessionCreateService(claims, cr)
	if err != nil {
		if !flag {
			res.Ok(c, "已经存在新的会话", sessionID)
			return
		}
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}

	res.Ok(c, "会话创建成功", sessionID)

}

// BigModelChatCreateView 创建大模型对话（SSE流）
// @Summary 创建大模型对话（SSE流）
// @Description 需用户认证，通过SSE（Server-Sent Events）实时返回对话内容
// @Tags 大模型对话管理
// @Accept application/json
// @Produce text/event-stream
// @Security ApiKeyAuth
// @Param sessionID query uint true "会话ID（关联已创建的会话）"
// @Param content query string true "对话内容（用户输入的消息）"
// @Success 200 {string} string "SSE事件流（实时返回AI的响应内容，格式为event-stream）"
// @Failure 1001 {string} string "参数错误（如sessionID或content为空，通过SSE返回错误信息）"
// @Failure 1002 {string} string "服务异常（如对话生成失败，通过SSE返回错误信息）"
// @Router /big_model/chat_sse [get]
func (BigModelApi) BigModelChatCreateView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Header("Access-Control-Allow-Origin", "*")

	// 确保在设置头部后立即开始流式响应
	c.Writer.WriteHeader(http.StatusOK)
	var cr big_model_service.ChatCreateRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErrSSE(c, err)
		return
	}
	big_model_service.BigModelChatCreateService(c, claims, cr)

}

// BigModelSessionListView 获取大模型会话分页列表
// @Summary 获取大模型会话分页列表
// @Description 仅管理员可调用，分页查询所有用户的大模型会话记录
// @Tags 大模型会话管理
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param page query int false "页码（默认1）"
// @Param pageSize query int false "每页条数（默认10）"
// @Success 200 {object} res.Response{code=int,data=res.DataListResponse{List=[]big_model_service.SessionListResponse,Count=int},msg=string} "成功（返回会话列表和总数，List为会话详情数组）"
// @Failure 1001 {object} res.Response{code=int,data=map[string]interface{},msg=string} "参数错误（如分页参数格式错误）"
// @Failure 1002 {object} res.Response{code=int,data=map[string]interface{},msg=string} "服务异常（如查询失败）"
// @Router /big_model/session [get]
func (BigModelApi) BigModelSessionListView(c *gin.Context) {
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
	err, responses, count := big_model_service.BigModelSessionListService(cr)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	res.OkWithList(c, responses, count)
}

// BigModelUserUpdateNameView 用户修改会话名称
// @Summary 用户修改会话名称
// @Description 需用户认证，修改当前用户所属会话的名称
// @Tags 大模型会话管理
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param data body big_model_service.SessionUserUpdateNameRequest true "修改参数"
// @Param data.SessionID body uint true "会话ID（需属于当前用户）"
// @Param data.Name body string false "新会话名称（为空则不修改）"
// @Success 200 {object} res.Response{code=int,data=map[string]interface{},msg=string} "成功（msg为“修改成功”）"
// @Failure 1001 {object} res.Response{code=int,data=map[string]interface{},msg=string} "参数错误（如SessionID为空）"
// @Failure 1002 {object} res.Response{code=int,data=map[string]interface{},msg=string} "服务异常（如会话不存在或无权限）"
// @Router /big_model/session [put]
func (BigModelApi) BigModelUserUpdateNameView(c *gin.Context) {
	_Claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _Claims.(*jwts.MyClaims)
	var cr big_model_service.SessionUserUpdateNameRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	err = big_model_service.BigModelUserUpdateNameService(claims, cr)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	res.OkWithMessage(c, "修改成功")

}

// BigModelUserDeleteSessionView 用户删除单个会话
// @Summary 用户删除单个会话
// @Description 需用户认证，删除当前用户所属的指定会话（含关联对话记录）
// @Tags 大模型会话管理
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param id path uint true "会话ID（需属于当前用户）"
// @Success 200 {object} res.Response{code=int,data=map[string]interface{},msg=string} "成功（msg为“删除成功”）"
// @Failure 1001 {object} res.Response{code=int,data=map[string]interface{},msg=string} "参数错误（如ID格式错误）"
// @Failure 1002 {object} res.Response{code=int,data=map[string]interface{},msg=string} "服务异常（如会话不存在或无权限）"
// @Router /big_model/session/{id} [delete]
func (BigModelApi) BigModelUserDeleteSessionView(c *gin.Context) {
	_Claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _Claims.(*jwts.MyClaims)
	var request models.IDRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	err = big_model_service.BigModelUserDeleteSessionService(claims, request)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	res.OkWithMessage(c, "删除成功")

}

// BigModelAdminDeleteSessionView 管理员批量删除会话
// @Summary 管理员批量删除会话
// @Description 仅管理员可调用，批量删除指定ID的会话（支持跨用户删除）
// @Tags 大模型会话管理
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param data body models.RemoveRequest true "批量删除参数" remark "包含需要删除的会话ID列表（格式：{\"IDs\": [1,2,3]}）"
// @Success 200 {object} res.Response{code=int,data=map[string]interface{},msg=string} "成功（返回删除结果描述）"
// @Failure 1001 {object} res.Response{code=int,data=map[string]interface{},msg=string} "参数错误（如ID列表为空）"
// @Failure 1002 {object} res.Response{code=int,data=map[string]interface{},msg=string} "服务异常（如删除失败）"
// @Router /big_model/session [delete]
func (BigModelApi) BigModelAdminDeleteSessionView(c *gin.Context) {
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
	err, msg := big_model_service.BigModelAdminDeleteSessionService(cr)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	res.FailWithMsg(c, msg)
}

// BigModelRoleDetailView 获取大模型角色详情
// @Summary 获取大模型角色详情
// @Description 查询指定ID的大模型角色详细信息（含关联标签、聊天次数等）
// @Tags 大模型角色管理
// @Accept application/json
// @Produce application/json
// @Param id path uint true "角色ID"
// @Success 200 {object} res.Response{code=int,data=big_model_service.RoleDetailResponse,msg=string} "成功（返回角色详情，包含名称、图标、标签等信息）"
// @Failure 1001 {object} res.Response{code=int,data=map[string]interface{},msg=string} "参数错误（如ID格式错误）"
// @Failure 1002 {object} res.Response{code=int,data=map[string]interface{},msg=string} "服务异常（如角色不存在）"
// @Router /big_model/roles/{id} [get]
func (BigModelApi) BigModelRoleDetailView(c *gin.Context) {
	var cr models.IDRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	err, response := big_model_service.BigModelRoleDetailService(cr)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	res.OkWithData(c, response)
}

// BigModelUserRoleHistoryView 用户获取使用过的大模型聊天历史
// @Summary 用户获取使用过的大模型聊天历史
// @Description 需用户认证，返回当前用户与所有使用过的大模型角色的聊天历史记录
// @Tags 大模型对话管理
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Success 200 {object} res.Response{code=int,data=[]interface{},msg=string} "成功（返回历史记录列表）"
// @Failure 1002 {object} res.Response{code=int,data=map[string]interface{},msg=string} "服务异常（如查询失败）"
// @Router /big_model/roles_history [get]
func (BigModelApi) BigModelUserRoleHistoryView(c *gin.Context) {
	_Claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _Claims.(*jwts.MyClaims)
	list := big_model_service.BigModelUserRoleHistoryService(claims)

	res.OkWithData(c, list)
}

// BigModelChatListView 获取单个会话的聊天记录列表
// @Summary 获取单个会话的聊天记录列表
// @Description 需用户认证，分页查询指定会话的所有聊天记录（用户与AI的交互内容）
// @Tags 大模型对话管理
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param sessionID query uint true "会话ID（需属于当前用户）"
// @Param page query int false "页码（默认1）"
// @Param pageSize query int false "每页条数（默认10）"
// @Success 200 {object} res.Response{code=int,data=res.DataListResponse{List=[]big_model_service.ChatListResponse,Count=int},msg=string} "成功（返回聊天记录列表和总数）"
// @Failure 1001 {object} res.Response{code=int,data=map[string]interface{},msg=string} "参数错误（如sessionID为空或分页参数无效）"
// @Failure 1002 {object} res.Response{code=int,data=map[string]interface{},msg=string} "服务异常（如会话不存在或无权限）"
// @Router /big_model/chat [get]
func (BigModelApi) BigModelChatListView(c *gin.Context) {
	_Claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _Claims.(*jwts.MyClaims)
	var cr big_model_service.ChatListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	err, responses, count := big_model_service.BigModelChatListService(cr, claims)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	res.OkWithList(c, responses, count)

}

// BigModelUserChatDeleteView 用户删除单个对话记录
// @Summary 用户删除单个对话记录
// @Description 需用户认证，删除当前用户所属会话中的指定对话记录
// @Tags 大模型对话管理
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param id path uint true "对话记录ID（需属于当前用户）"
// @Success 200 {object} res.Response{code=int,data=map[string]interface{},msg=string} "成功（msg为“删除成功”）"
// @Failure 1001 {object} res.Response{code=int,data=map[string]interface{},msg=string} "参数错误（如ID格式错误）"
// @Failure 1002 {object} res.Response{code=int,data=map[string]interface{},msg=string} "服务异常（如对话记录不存在或无权限）"
// @Router /big_model/chat/{id} [delete]
func (BigModelApi) BigModelUserChatDeleteView(c *gin.Context) {
	_Claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _Claims.(*jwts.MyClaims)

	var cr models.IDRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	err = big_model_service.BigModelUserChatDeleteService(cr, claims)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	res.OkWithMessage(c, "删除成功")

}

// BigModelAdminChatDeleteView 管理员批量删除对话记录
// @Summary 管理员批量删除对话记录
// @Description 仅管理员可调用，批量删除指定ID的对话记录（支持跨用户删除）
// @Tags 大模型对话管理
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param data body models.RemoveRequest true "批量删除参数" remark "包含需要删除的对话记录ID列表（格式：{\"IDs\": [1,2,3]}）"
// @Success 200 {object} res.Response{code=int,data=map[string]interface{},msg=string} "成功（返回删除结果描述）"
// @Failure 1001 {object} res.Response{code=int,data=map[string]interface{},msg=string} "参数错误（如ID列表为空）"
// @Failure 1002 {object} res.Response{code=int,data=map[string]interface{},msg=string} "服务异常（如删除失败）"
// @Router /big_model/chat [delete]
func (BigModelApi) BigModelAdminChatDeleteView(c *gin.Context) {
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
	err, msg := big_model_service.BigModelAdMINChatDeleteService(cr)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	res.OkWithMessage(c, msg)

}

func (BigModelApi) BigModelRoleTagsListView(c *gin.Context) {

	_, exist := c.Get("claims")
	if !exist {
		return
	}
	var list []models.Options[uint]
	global.DB.Model(models.BigModelTagModel{}).Select("id as value", "title as label").Scan(&list)
	res.OkWithData(c, list)

}

func (BigModelApi) IconView(c *gin.Context) {
	dir, err := os.ReadDir("uploads/role_icons")
	if err != nil {
		logrus.Error(err)
		res.FailWithMsg(c, "目录不存在")
		return
	}
	var list []models.Options[string]
	for _, entry := range dir {
		key := "/" + path.Join("uploads/role_icons", entry.Name())
		list = append(list, models.Options[string]{
			Label: key,
			Value: key,
		})
	}
	res.OkWithData(c, list)
}

func (BigModelApi) BigModelRoleSessionListView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)

	var cr big_model_service.RoleSessionsRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	err, count, responses := big_model_service.BigModelRoleSessionListService(cr, claims)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
		return
	}
	res.OkWithList(c, responses, count)

}
