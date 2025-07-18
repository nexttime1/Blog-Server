package middleware

import (
	"Blog_server/service/log_service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResponseWriter  实现了gin.ResponseWriter 接口
/*
显式实现了 Write() 方法
嵌入机制自动提供了 Header() 和 WriteHeader() 方法
*/
type ResponseWriter struct {
	gin.ResponseWriter
	Body []byte
	Head http.Header
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	fmt.Println("response :", string(data))
	w.Body = append(w.Body, data...)
	return w.ResponseWriter.Write(data)
}

func (w *ResponseWriter) Header() http.Header {
	return w.Head
}

func LogMiddleware(c *gin.Context) {
	//请求部分
	fmt.Println("中间件") //先走中间件
	log := log_service.NewActionLogByGin(c)

	log.SetRequest(c)
	//"log"是设定的键，这是个字符串类型。  存到上下文中   c.get就开源拿出来
	//log则是要存储的值，它属于*ActionLog类型。
	c.Set("log", log)

	res := &ResponseWriter{
		ResponseWriter: c.Writer, //只要一个类型实现了接口的所有方法，就可以将该类型的值赋值给该接口变量 go 多态
		Head:           make(http.Header),
	}
	c.Writer = res

	c.Next()
	//响应部分
	fmt.Println("中间件响应") //走完 路由中的函数 再走这里
	fmt.Println("res", string(res.Body))
	log.SetResponse(res.Body)
	log.SetResponseHeader(res.Head)

	log.MiddleSave()

}
