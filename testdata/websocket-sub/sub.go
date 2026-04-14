package main

// 比如服务端通过维护一个 hub 管理所有客户端活跃的链接。
type Hub struct {
	// Registered clients.
	clients map[*Client]struct{}
	// Inbound messages from the clients.
	broadcast chan []byte
	// Register requests from the clients.
	register chan *Client
	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]struct{}),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register: // 升级成 ws 协议之后，注册 client 到 hub
			h.clients[client] = struct{}{}
		case client := <-h.unregister: // 连接断开则注销 client 实例，释放资源
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast: // 一个阻塞的 channal，接收来自 client 的需要广播的消息
			for client := range h.clients { // 发送给所有的 client 实例
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
