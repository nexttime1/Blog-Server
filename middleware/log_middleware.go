package middleware

import (
	"Blog_server/service/log_service"
	"fmt"
	"github.com/gin-gonic/gin"
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
	fmt.Println("中间件") //先走中间件
	log := log_service.NewActionLogByGin(c)

	log.SetRequest(c)
	c.Set("log", log)

	res := &ResponseWriter{
		ResponseWriter: c.Writer,
	}
	c.Writer = res

	c.Next()
	//响应部分
	fmt.Println("中间件响应") //走完 路由中的函数 再走这里
	fmt.Println("res", string(res.Body))
	log.SetResponse(res.Body)
	log.Save()

}
