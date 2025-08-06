package user_api

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"Blog_server/service/redis_service/redis_jwt"
	"Blog_server/service/user_service"
	"Blog_server/utils/jwts"
	"Blog_server/utils/pwd"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserApi struct {
}

type UserTypeUpdateRequest struct {
	Role     enum.RoleType `json:"role" binding:"required,oneof=1 2 3 4" gorm:"column:nickname;omitempty"`
	NickName string        `json:"nick_name"` // 防止用户昵称非法，管理员有能力修改
	UserID   uint          `json:"user_id" binding:"required"`
}

type UserPwdUpdateRequest struct {
	Pwd    string `json:"pwd" binding:"required"`
	NewPwd string `json:"new_pwd" binding:"required"`
}

// UserEmailLogin 邮箱登录
// @Summary 邮箱登录
// @Description 邮箱登录 输入 用户名和密码  都是必填
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body user_service.EmailLoginRequest true "登录信息"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response "登录成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/email_login [post]
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

// UserInfoView 用户列表
// @Summary 用户列表
// @Description 获取用户列表
// @Tags 用户管理
// @Produce json
// @Param page query int false "页码，默认1" mininum(1)
// @Param limit query int false "每页条数，默认10" mininum(1) maxinum(100)
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{data=res.DataListResponse}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/users [get]
func (u UserApi) UserInfoView(c *gin.Context) {
	var cr user_service.UserInfoRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	_, exist := c.Get("claims")
	if !exist {
		// 有问题 AuthMiddleware 已经res 返回了 这里不需要返回
		return
	}
	UserModelList, count, err := user_service.UserInfoService(cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithList(c, UserModelList, count)

}

// UserUpdateView 更新权限
// @Summary 更新权限
// @Description 更新权限
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body UserTypeUpdateRequest true "要更新的权限和用户id（必填）"
// @Param token header string true "token"
// @Success 200 {object} res.Response "更新成功"
// @Failure 400 {object} res.Response "请求参数错误（如ID格式错误）"
// @Failure 404 {object} res.Response "不存在"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/users_role [put]
func (u UserApi) UserUpdateView(c *gin.Context) {
	var cr UserTypeUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	_, ok := c.Get("claims")
	if !ok {
		return
	}
	var model models.UserModel
	err = global.DB.Where("id = ?", cr.UserID).Take(&model).Error
	if err != nil {
		res.FailWithMsg(c, "id 未找到")
		return
	}
	fmt.Println(cr.Role)
	err = global.DB.Model(&model).Updates(map[string]any{
		"nickname": cr.NickName,
		"role":     cr.Role, // 这里会自动调用 cr.Role.Value()
	}).Error
	if err != nil {
		res.FailWithMsg(c, "修改权限错误")
		return
	}
	res.OkWithMessage(c, "修改权限成功")

}

// UserPasswordView 更新密码
// @Summary 更新密码
// @Description 更新密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body UserPwdUpdateRequest true "要更新的密码 和 旧密码（必填）"
// @Param token header string true "token"
// @Success 200 {object} res.Response "更新成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 404 {object} res.Response "不存在"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/users_pwd [put]
func (u UserApi) UserPasswordView(c *gin.Context) {
	var cr UserPwdUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	//改自己的 密码
	_claim, exists := c.Get("claims")
	if !exists {
		res.FailWithMsg(c, "请登录")
		return
	}
	claim := _claim.(*jwts.MyClaims)
	var model models.UserModel
	err = global.DB.Where("id = ?", claim.UserID).Take(&model).Error
	if err != nil {
		res.FailWithMsg(c, "没有该用户")
		return
	}
	//核实旧密码
	ok := pwd.CheckPwd(model.Password, cr.Pwd)
	if !ok {
		res.FailWithMsg(c, "旧密码错误 无法修改")
		return
	}
	// 可以修改  先加密
	hashPwd := pwd.HashPwd(cr.NewPwd)
	global.DB.Model(&model).Update("password", hashPwd)

	//将token 加入黑名单  因为 修改密码了
	err = redis_jwt.TokenBlackByGin(c, redis_jwt.UserBlackType)
	if err != nil {
		res.FailWithMsg(c, "token 加入黑名单失败")
	}

	res.OkWithData(c, "修改密码成功")
}
