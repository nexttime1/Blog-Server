package res

import (
	"Blog_server/utils/validate"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code Code        `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

var empty = map[string]interface{}{}

type Code int

const (
	SuccessCode     Code = 0
	FailValidCode   Code = 1001
	FailServiceCode Code = 1002 //服务异常

)

func (c Code) Message() string {
	switch c {
	case SuccessCode:
		return "成功"
	case FailValidCode:
		return "参数校验失败"
	case FailServiceCode:
		return "服务异常"
	}
	return ""
}

func (r Response) Json(c *gin.Context) {
	c.JSON(200, r)
}

func Ok(c *gin.Context, message string, data interface{}) {
	Response{SuccessCode, data, message}.Json(c)
}
func OkWithMessage(c *gin.Context, message string) {
	Response{SuccessCode, empty, message}.Json(c)
}

func OkWithData(c *gin.Context, data interface{}) {
	Response{SuccessCode, data, "成功"}.Json(c)

}
func FailWithErr(c *gin.Context, err error) {
	data, msg := validate.ValidateErr(err)
	FailWithData(c, msg, data)
}

func FailWithMsg(c *gin.Context, message string) {
	Response{FailValidCode, empty, message}.Json(c)
}

func FailWithData(c *gin.Context, message string, data interface{}) {
	Response{FailServiceCode, data, message}.Json(c)
}
func FailWithCode(c *gin.Context, code Code) {
	Response{code, empty, code.Message()}.Json(c)
}

func OkWithList(c *gin.Context, List any, Count int64) {
	Response{SuccessCode, map[string]any{
		"List":  List,
		"Count": Count,
	}, "成功"}.Json(c)
}
