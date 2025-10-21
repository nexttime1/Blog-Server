package middleware

import (
	"Blog_server/common/res"
	"Blog_server/models/enum"
	"Blog_server/service/redis_service/redis_jwt"
	"Blog_server/utils/jwts"
	"errors"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		res.FailWithErr(c, err)
		c.Abort()
		return
	}
	//判断是否在黑名单里
	ok, blackType := redis_jwt.HasTokenBlackByGin(c)
	if ok { //ok = true  的话 在黑名单 不能再走了
		res.FailWithMsg(c, blackType.Msg())
		c.Abort() //后面请求响应都不走  但	c.Set  要走  所以要return
		return
	}
	c.Set("claims", claims)
	c.Next()
	return
}

func AdminMiddleware(c *gin.Context) {
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		res.FailWithErr(c, err)
		c.Abort()
		return
	}
	//判断是否在黑名单里
	ok, blackType := redis_jwt.HasTokenBlackByGin(c)
	if ok { //ok = true  的话 在黑名单 不能再走了
		res.FailWithMsg(c, blackType.Msg())
		c.Abort() //后面请求响应都不走  但	c.Set  要走  所以要return
		return
	}

	if claims.Role != enum.AdminRole {
		res.FailWithErr(c, errors.New("权限错误"))
		c.Abort()
		return
	}
	c.Set("claims", claims)
	c.Next()
}

func AuthSSEMiddleware(c *gin.Context) {
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		res.FailWithErrSSE(c, err)
		c.Abort()
		return
	}
	//判断是否在黑名单里
	ok, blackType := redis_jwt.HasTokenBlackByGin(c)
	if ok { //ok = true  的话 在黑名单 不能再走了
		res.FailWithMsgSSE(c, blackType.Msg())
		c.Abort() //后面请求响应都不走  但	c.Set  要走  所以要return
		return
	}
	c.Set("claims", claims)
	c.Next()
	return
}
