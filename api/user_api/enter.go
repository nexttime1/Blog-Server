package user_api

import (
	"Blog_server/common"
	"Blog_server/common/res"
	"Blog_server/core"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"Blog_server/plugins/email"
	"Blog_server/plugins/qq"
	"Blog_server/service/log_service"
	"Blog_server/service/redis_service/redis_jwt"
	"Blog_server/service/user_service"
	"Blog_server/utils/jwts"
	"Blog_server/utils/pwd"
	"Blog_server/utils/random"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/liu-cn/json-filter/filter"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserApi struct {
}

type UserTypeUpdateRequest struct {
	Role     enum.RoleType `json:"role" binding:"required,oneof=1 2 3 4" gorm:"omitempty"`
	NickName string        `json:"nick_name"` // 防止用户昵称非法
	UserID   uint          `json:"user_id" binding:"required"`
}

type UserPwdUpdateRequest struct {
	Pwd    string `json:"pwd" binding:"required"`
	NewPwd string `json:"new_pwd" binding:"required"`
}

type UserBindEmailRequest struct {
	Email    string  `json:"email" binding:"required,email"`
	Code     *string `json:"code"`
	Password string  `json:"password"`
}

type UserInfoResponse struct {
	Username string        `json:"user_name"`
	Nickname string        `json:"nick_name"`
	Role     enum.RoleType `json:"role"`   // 1 为管理员  2 为 普通用户
	Avatar   string        `json:"avatar"` //头像
	Sign     string        `json:"sign"`
	Link     string        `json:"link"`
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
	model, token, msg, err := user_service.UserEmailLoginService(request)
	if err != nil {
		logrus.Errorf("UserEmailLogin %s  err:  %s", msg, err)
		//记录登录失败的日志
		log_service.NewLoginFail(c, enum.EmailLoginType, msg, request.Username, request.Password)
		res.FailWithMsg(c, msg)
		return
	}
	//记录登录日志
	log_service.NewLoginSuccess(c, enum.EmailLoginType, model)

	//记录用户登录数据
	global.DB.Create(&models.LoginDataModel{
		UserID:    model.ID,
		IP:        c.ClientIP(),
		NickName:  model.Nickname,
		Token:     token,
		Device:    "",
		Addr:      core.GetIpAddr(c.ClientIP()),
		LoginType: enum.EmailLoginType,
	})

	res.OkWithData(c, token)
}

// UserListInfoView 用户列表
// @Summary 用户列表
// @Description 获取用户列表
// @Tags 用户管理
// @Produce json
// @Param page query int false "页码，默认1" mininum(1)
// @Param limit query int false "每页条数，默认10" mininum(1) maxinum(100)
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{data= user_service.UserListResponse}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/users [get]
func (u UserApi) UserListInfoView(c *gin.Context) {
	_, exist := c.Get("claims")
	if !exist {
		// 有问题 AuthMiddleware 已经res 返回了 这里不需要返回
		return
	}
	var cr user_service.UserInfoRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
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
// @Router /api/user_role [put]
func (u UserApi) UserUpdateView(c *gin.Context) {
	_, ok := c.Get("claims")
	if !ok {
		return
	}
	var cr UserTypeUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
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
	//改自己的 密码
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	var cr UserPwdUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
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

// UserDeleteView 批量删除用户
// @Summary 批量删除用户
// @Description 批量删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body models.RemoveRequest true "用户id列表"
// @Param token header string true "token"
// @Success 200 {object} res.Response "删除成功"
// @Failure 400 {object} res.Response "请求参数错误（如id格式错误、列表为空）"
// @Failure 404 {object} res.Response "部分用户id不存在"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/users [delete]
func (u UserApi) UserDeleteView(c *gin.Context) {
	//这个是管理员 才能进的
	_, exists := c.Get("claims")
	if !exists {
		return
	}
	var mr models.RemoveRequest
	err := c.ShouldBindJSON(&mr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	var Count int
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		var model []models.UserModel
		ModelList := common.BatchRemove(model, mr)
		Count = len(ModelList)
		//todo: 删除关联表

		return nil
	})
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithMessage(c, fmt.Sprintf("删除成功，共删除%d个用户", Count))

}

// UserBindEmailView 用户绑定邮箱
// @Summary 用户绑定邮箱
// @Description 用户绑定邮箱 输入 邮箱 密码 验证码  第一次 只填 邮箱 必填  第二次 都是必填
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body user_service.EmailLoginRequest true "登录信息"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response "登录成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/email_login [post]
func (u UserApi) UserBindEmailView(c *gin.Context) {
	_claims, exists := c.Get("claims")
	if !exists {
		return
	}
	claims := _claims.(*jwts.MyClaims)

	var mr UserBindEmailRequest
	err := c.ShouldBindJSON(&mr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	session := sessions.Default(c)
	if mr.Code == nil {
		//第一次发送
		code := random.DigitCode()
		err = email.NewCode().Send(mr.Email, "您的验证码是: "+code)
		if err != nil {
			res.FailWithMsg(c, fmt.Sprintf("邮件发送失败  %v", err))
			return
		}
		// 设置session过期时间  10分钟
		session.Options(sessions.Options{MaxAge: 600})
		session.Set("valid_code", code)
		session.Set("email", mr.Email)
		fmt.Println("保存的 code 和 email ", code, mr.Email)
		err = session.Save()
		if err != nil {
			res.FailWithMsg(c, fmt.Sprintf("session 保存错误 %v", err))
			return
		}
		res.OkWithMessage(c, "验证码发送成功 请注意查收")
		return
	}
	//发送验证码后 再一次调用
	validCode := session.Get("valid_code")
	validEmail := session.Get("email")
	fmt.Println(validEmail)
	fmt.Println(mr.Email)
	if len(mr.Password) < 4 {
		res.FailWithMsg(c, "密码强度过低")
		return
	}
	//必须保证 第一次发的 邮箱 和 第二次发的邮箱 一样
	if validEmail != mr.Email {
		res.FailWithMsg(c, "邮箱绑定失败 请使用接收验证码的邮箱")
		return
	}

	if validCode != *mr.Code {
		res.FailWithMsg(c, "验证码错误")
		return
	}
	var model models.UserModel
	err = global.DB.Where("id = ?", claims.UserID).Take(&model).Error
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("未找到该用户 %v", err))
		return
	}

	hashPwd := pwd.HashPwd(mr.Password)

	err = global.DB.Model(&model).Updates(map[string]any{
		"email":    mr.Email,
		"password": hashPwd,
	}).Error
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("用户绑定邮箱失败 %v", err))
		return
	}
	// 绑定成功后清理session
	session.Delete("valid_code")
	session.Delete("email")
	session.Save() // 保存修改
	res.OkWithMessage(c, "绑定成功")
}

// UserQQLogin 用户QQ登录
// @Summary 用户QQ登录
// @Description 通过QQ授权码进行登录，支持新用户自动注册
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param code query string true "QQ授权返回的code参数"
// @Success 200 {object} res.Response{data=string} "登录成功，返回token"
// @Failure 400 {object} res.Response "请求参数错误（如缺少code）"
// @Failure 500 {object} res.Response "服务器内部错误（如QQ接口调用失败、数据库错误等）"
// @Router /api/qq_login [post]
func (u UserApi) UserQQLogin(c *gin.Context) {
	code := c.Query("code")
	qqInfo, err := qq.NewQQLogin(code)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	var model models.UserModel
	err = global.DB.Where("token = ?", qqInfo.OpenID).Take(&model).Error
	if err != nil {
		// 没找到  注册
		model = models.UserModel{
			Username:       qqInfo.OpenID,
			Nickname:       qqInfo.Nickname,
			Password:       random.GenerateRandomString(8),
			Avatar:         qqInfo.Avatar,
			RegisterSource: enum.SignQQ,
			Addr:           core.GetIpAddr(c.ClientIP()),
			Token:          qqInfo.OpenID,
			IP:             c.ClientIP(),
			Role:           enum.UserRole,
		}
		err = global.DB.Create(&model).Error
		if err != nil {
			res.FailWithMsg(c, fmt.Sprintf("注册用户失败 %v", err))
			//失败日志
			log_service.NewLoginFail(c, enum.QQLoginType, fmt.Sprintf("注册用户失败 %v", err), model.Username, model.Password)
			return
		}
	}
	//登录操作
	token, err := jwts.GetToken(jwts.Claims{
		UserID:   model.ID,
		Username: model.Username,
		Role:     model.Role,
	})
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("token 申请失败 %v", err))
		//失败日志
		log_service.NewLoginFail(c, enum.QQLoginType, fmt.Sprintf("token 申请失败 %v", err), model.Username, model.Password)
		return
	}
	// 登录成功日志
	log_service.NewLoginSuccess(c, enum.QQLoginType, model)

	//记录用户登录数据
	global.DB.Create(&models.LoginDataModel{
		UserID:    model.ID,
		IP:        c.ClientIP(),
		NickName:  model.Nickname,
		Token:     token,
		Device:    "",
		Addr:      core.GetIpAddr(c.ClientIP()),
		LoginType: enum.QQLoginType,
	})
	res.OkWithData(c, token)
}

// UserPersonInfoView 用户个人信息
// @Tags 用户管理
// @Summary 用户个人信息
// @Description 用户个人信息
// @Router /api/user_info [get]
// @Param token header string  true  "token"
// @Produce json
// @Success 200 {object} res.Response{data=models.UserModel}
func (UserApi) UserPersonInfoView(c *gin.Context) {
	_claims, exists := c.Get("claims")
	if !exists {
		return
	}
	claims := _claims.(*jwts.MyClaims)
	var userInfo models.UserModel
	err := global.DB.Where("id = ?", claims.UserID).Take(&userInfo).Error
	if err != nil {
		res.FailWithMsg(c, "用户不存在")
	}
	res.OkWithData(c, filter.Select("info", userInfo))

}

// UserCreateView 创建用户
// @Summary 创建用户
// @Description 创建用户 输入 昵称 用户名 密码 权限
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body user_service.UserCreateRequest true "创建信息"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response "创建成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/users [post]
func (u UserApi) UserCreateView(c *gin.Context) {
	_, exists := c.Get("claims")
	if !exists {
		return
	}
	var mr user_service.UserCreateRequest
	err := c.ShouldBindJSON(&mr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	err = user_service.UserCreateService(mr.UserName, mr.NickName, mr.Password, mr.Role, "", c.ClientIP(), core.GetIpAddr(c.ClientIP()))
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("创造用户失败 %v", err))
		return
	}

	res.OkWithMessage(c, fmt.Sprintf("创建%s用户成功", mr.UserName))
}

// UserUpdateInfoView 修改当前登录人的昵称，签名，链接
// @Tags 用户管理
// @Summary 修改当前登录人的昵称，签名，链接
// @Description 修改当前登录人的昵称，签名，链接
// @Router /api/user_info [put]
// @Param token header string  true  "token"
// @Param data body user_service.UserUpdateInfoRequest  true  "昵称，签名，链接"
// @Produce json
// @Success 200 {object} res.Response{}
func (UserApi) UserUpdateInfoView(c *gin.Context) {
	_Claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claims := _Claim.(*jwts.MyClaims)

	var cr user_service.UserUpdateInfoRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	err = user_service.UserInfoPutService(cr, claims)
	if err != nil {
		res.FailWithMsg(c, err.Error())
		return
	}

	res.OkWithMessage(c, "修改个人信息成功")

}

// UserInfoView 用户信息
// @Tags 用户管理
// @Summary 用户信息
// @Description 用户信息
// @Router /api/user_info [get]
// @Param token header string  true  "token"
// @Produce json
// @Success 200 {object} res.Response{data= UserInfoResponse }
func (UserApi) UserInfoView(c *gin.Context) {

	_claims, exist := c.Get("claims")
	if !exist {
		logrus.Errorf("没登陆")
		return
	}
	claims := _claims.(*jwts.MyClaims)

	var userInfo models.UserModel
	err := global.DB.Take(&userInfo, claims.UserID).Error
	if err != nil {
		res.FailWithMsg(c, "用户不存在")
		return
	}
	response := UserInfoResponse{
		Username: userInfo.Username,
		Nickname: userInfo.Nickname,
		Role:     userInfo.Role,
		Avatar:   userInfo.Avatar,
		Sign:     userInfo.Sign,
		Link:     userInfo.Link,
	}

	res.OkWithData(c, response)

}

// QQLoginLinkView 获取qq登录的跳转链接
// @Tags 用户管理
// @Summary 获取qq登录的跳转链接
// @Description 获取qq登录的跳转链接,data就是qq的跳转地址
// @Router /api/qq_login_path [get]
// @Produce json
// @Success 200 {object} res.Response{}
func (UserApi) QQLoginLinkView(c *gin.Context) {
	path := global.Config.QQ.GetPath()
	if path == "" {
		res.FailWithMsg(c, "未配置qq登录地址")
		return
	}
	res.OkWithData(c, path)
	return
}

func (UserApi) UserLogoutView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)

	var user models.UserModel
	err := global.DB.Where("id = ?", claim.UserID).Take(&user).Error
	if err != nil {
		res.FailWithMsg(c, "没有该用户")
		return
	}

	err = redis_jwt.TokenBlackByGin(c, redis_jwt.UserBlackType)
	if err != nil {
		logrus.Errorf(" %v", err)
		res.FailWithMsg(c, "注销失败")
		return
	}
	log_service.LogoutSuccess(c, user)

	logrus.Info(fmt.Sprintf("用户 %s 注销登录", claim.Username))

	res.OkWithMessage(c, "注销成功")

}
