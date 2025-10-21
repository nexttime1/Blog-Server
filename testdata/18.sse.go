package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

func StreamView(c *gin.Context) {
	var IntChan = make(chan int, 0)
	go func() {
		for i := 0; i < 10; i++ {
			IntChan <- i
			time.Sleep(time.Millisecond * 500)
		}
		close(IntChan) // 关闭通道
	}()
	/*
		func(w io.Writer) bool 是一个「回调函数」，它的作用是：
		每次被调用时，负责生成一段要发送给客户端的数据（通过 w 写入）。
		返回值 bool 决定是否继续流式传输：true 表示 “还有数据要发，继续调用我”；false 表示 “所有数据发送完毕，关闭流”。
		就跟个 函数一样
	*/
	c.Stream(func(w io.Writer) bool {
		if i, ok := <-IntChan; ok {

			/*
				SSE（Server-Sent Events）是一种基于 HTTP 的协议，专门用于 服务器主动向客户端持续推送数据（单向通信）。
				它有固定的数据格式要求（比如必须以 data: 开头，以换行符结束），客户端（如浏览器）可以通过 EventSource API 监听这种流，实时接收数据。
				c.SSEvent 的参数含义：
				第一个参数是「事件类型」（字符串），空字符串表示默认事件（客户端可通过 onmessage 监听）。
				第二个参数是要发送的数据（这里是 i），Gin 会自动将其序列化为符合 SSE 格式的字符串（例如 data: 0\n\n、data: 1\n\n 等）。
			*/
			c.SSEvent("", i)
			return true
		}
		return false
	})

}

func main() {

	r := gin.Default()

	// 用stream 流
	r.GET("/", StreamView)

	r.Run(":8080")

}
