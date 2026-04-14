package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// 单条消息最大限制
	maxMessageSize = 512
	// 等待客户端 pong 响应超时
	pongWait = 60 * time.Second
	// 服务端定时 ping 间隔
	pingPeriod = (pongWait * 9) / 10
	// 写入消息超时
	writeWait = 10 * time.Second
)

var (
	addr     = flag.String("addr", "localhost:8080", "http service address")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome) // HTTP 请求
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r) // HTTP 协议升级成 ws 协议
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // 劫持 conn，并升级为 ws 协议的函数，这部分是 websocket 库实现的
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	// 极简前端HTML，开箱即用
	html := `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>WebSocket 测试</title>
</head>
<body>
<input id="msg" placeholder="输入消息">
<button onclick="send()">发送</button>
<div id="content"></div>

<script>
let ws = new WebSocket("ws://"+location.host+"/ws");
ws.onmessage = function(e){
    let div = document.createElement("div");
    div.innerText = e.data;
    document.getElementById("content").appendChild(div);
}
function send(){
    let val = document.getElementById("msg").value;
    ws.send(val);
}
</script>
</body>
</html>
`
	_, _ = w.Write([]byte(html))
}
