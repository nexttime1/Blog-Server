package res

import (
	"Blog_server/utils/validate"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Response struct {
	Code Code        `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

var empty = map[string]interface{}{}

type Code int

const (
	SuccessCode      Code = 0
	FailValidCode    Code = 1001
	FailServiceCode  Code = 1002 //服务异常
	FailArgumentCode Code = 1003
	ArgumentError    Code = 1004
)

func (c Code) Message() string {
	switch c {
	case SuccessCode:
		return "成功"
	case FailValidCode:
		return "参数校验失败"
	case FailServiceCode:
		return "服务异常"
	case FailArgumentCode:
		return "参数错误"
	}
	return ""
}

type DataListResponse struct {
	List  any `json:"list"`
	Count int `json:"count"`
}

func (r Response) Json(c *gin.Context) {
	c.JSON(200, r)
}

func (r Response) ToJson() string {
	byteData, _ := json.Marshal(r)
	return string(byteData)

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

func OkWithList(c *gin.Context, List any, Count int) {
	Response{SuccessCode, DataListResponse{
		List:  List,
		Count: Count,
	}, "成功"}.Json(c)
}

// 辅助函数：设置SSE标准响应头
func setSSEHeaders(c *gin.Context) {
	// 只有在头部未设置时才设置
	if c.Writer.Header().Get("Content-Type") == "" {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("X-Accel-Buffering", "no")
		c.Header("Access-Control-Allow-Origin", "*")
	}
}

// 辅助函数：按SSE标准格式写入数据（data: {json}\n\n）
func writeSSEData(w io.Writer, dataJson string) {
	// 拼接格式：data: + JSON字符串 + 双换行
	sseFormat := "data: " + dataJson + "\n\n"
	_, _ = w.Write([]byte(sseFormat))
	// 强制刷新缓冲区，确保数据立即发送
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// OkWithMessageSSE 发送带消息的SSE（用于流式增量内容）
func OkWithMessageSSE(c *gin.Context, message string) {
	// 1. 设置SSE必须的响应头（每次调用都设置，避免被覆盖）
	setSSEHeaders(c)

	// 2. 构建响应体
	data := Response{
		Code: SuccessCode,
		Data: empty,
		Msg:  message, // 增量内容放在msg字段
	}.ToJson()

	// 3. 按SSE标准格式发送：data: {json}\n\n
	writeSSEData(c.Writer, data)
}

// OkWithDataSSE 发送带数据的SSE
func OkWithDataSSE(c *gin.Context, data interface{}) {
	setSSEHeaders(c)
	dataJson := Response{SuccessCode, data, "成功"}.ToJson()
	writeSSEData(c.Writer, dataJson)
}

// FailWithMsgSSE 发送错误消息的SSE
func FailWithMsgSSE(c *gin.Context, message string) {
	setSSEHeaders(c)
	data := Response{FailValidCode, empty, message}.ToJson()
	writeSSEData(c.Writer, data)
}

// FailWithDataSSE 发送带数据的错误SSE
func FailWithDataSSE(c *gin.Context, message string, data interface{}) {
	setSSEHeaders(c)
	dataJson := Response{FailServiceCode, data, message}.ToJson()
	writeSSEData(c.Writer, dataJson)
}

// FailWithErrSSE 发送错误的SSE
func FailWithErrSSE(c *gin.Context, err error) {
	data, msg := validate.ValidateErr(err)
	FailWithDataSSE(c, msg, data)
}

// OkWithSSE 发送带消息和数据的SSE（用于结束标识）
func OkWithSSE(c *gin.Context, message string, data interface{}) {
	setSSEHeaders(c)
	dataJson := Response{SuccessCode, data, message}.ToJson()
	writeSSEData(c.Writer, dataJson)
}
