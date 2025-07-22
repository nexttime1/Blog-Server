package jwts

import (
	"Blog_server/global"
	"Blog_server/models/enum"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type Claims struct {
	UserID   uint          `json:"userID"`
	Username string        `json:"username"`
	Role     enum.RoleType `json:"role"`
}

type MyClaims struct {
	Claims
	jwt.StandardClaims
}

// GetToken get token
func GetToken(claims Claims) (string, error) {
	cla := MyClaims{
		Claims: claims,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(global.Config.Jwt.Expire) * time.Hour).Unix(), // 过期时间
			Issuer:    global.Config.Jwt.Issuer,                                                   // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cla)
	return token.SignedString([]byte(global.Config.Jwt.Secret)) // 进行签名生成对应的token
}

// parse token
func ParseToken(tokenString string) (*MyClaims, error) {
	//如果 token 没有东西  说明应该去登录
	if tokenString == "" {
		return nil, errors.New("请登录")
	}

	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.Jwt.Secret), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, errors.New("token 过期")
		}
		if strings.Contains(err.Error(), "token is invalid") {
			return nil, errors.New("token 无效") //串改了
		}
		if strings.Contains(err.Error(), "token contains an invalid") {
			return nil, errors.New("token 非法") //串改了
		}
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// ParseTokenByGin 万一没在 请求头里
func ParseTokenByGin(c *gin.Context) (*MyClaims, error) {
	token := c.GetHeader("Token")
	if token == "" {
		token = c.Query("token")

	}

	return ParseToken(token)

}
