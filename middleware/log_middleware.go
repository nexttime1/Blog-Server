package middleware

import (
	"Blog_server/service/log_service"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
)

type ResponseWriter struct {
	gin.ResponseWriter
	Body []byte
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	fmt.Println("response :", string(data))
	w.Body = append(w.Body, data...)
	return w.ResponseWriter.Write(data)
}

func LogMiddleware(c *gin.Context) {
	//请求部分
	log := log_service.NewActionLogByGin(c)
	ByteData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Errorf(err.Error())
	}
	fmt.Println("Body", string(ByteData))

	c.Request.Body = io.NopCloser(bytes.NewReader(ByteData))
	res := &ResponseWriter{
		ResponseWriter: c.Writer,
	}
	c.Writer = res

	c.Next()
	//响应部分
	fmt.Println("res", string(res.Body))

}
