package chat_api

import (
	"Blog_server/common"
	"Blog_server/common/res"
	"Blog_server/core"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"encoding/json"
	"fmt"
	"github.com/DanPlayer/randomname"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/liu-cn/json-filter/filter"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type ChatApi struct {
}

const (
	InRoomMsg  enum.MsgType = 1
	TextMsg    enum.MsgType = 2
	ImageMsg   enum.MsgType = 3
	voiceMsg   enum.MsgType = 4
	videoMsg   enum.MsgType = 5
	SystemMsg  enum.MsgType = 6
	OutRoomMsg enum.MsgType = 7
)

type ChatUser struct {
	Conn     *websocket.Conn
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
}
type GroupRequest struct {
	Content string       `json:"content"`  // 聊天的内容
	MsgType enum.MsgType `json:"msg_type"` // 聊天类型
}
type GroupResponse struct {
	Content     string       `json:"content"`      // 聊天的内容
	MsgType     enum.MsgType `json:"msg_type"`     // 聊天类型
	NickName    string       `json:"nick_name"`    // 前端自己生成
	Avatar      string       `json:"avatar"`       // 头像
	OnlineCount int          `json:"online_count"` //在线人数
	Date        time.Time    `json:"date"`         // 消息的时间
}

var ConnGroupMap = make(map[string]ChatUser)

func (ChatApi) ChatGroupView(c *gin.Context) {
	// 定义一个 websocket.Upgrader 实例，用于将 HTTP 连接升级为 WebSocket 连接
	var upGrader = websocket.Upgrader{
		// 定义跨域检查函数，决定是否允许客户端连接
		CheckOrigin: func(r *http.Request) bool {
			// 此处直接返回 true，表示不做鉴权，允许所有来源的连接
			// 实际生产环境中应根据业务需求实现鉴权逻辑（如验证 token 等）
			return true
		},
	}

	// 调用 Upgrader.Upgrade 方法将 HTTP 连接升级为 WebSocket 连接
	// 参数说明：
	// c.Writer：HTTP 响应写入器（用于发送升级响应）
	// c.Request：HTTP 请求对象（包含客户端的连接信息）
	// nil：额外的响应头（此处不需要，传 nil）
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	// 打印升级过程中可能出现的错误（如协议不支持、鉴权失败等）
	fmt.Println(err)
	// 如果升级失败（err 不为 nil），返回参数错误响应
	if err != nil {
		// res.FailWithCode 是自定义的响应工具函数，用于返回错误码
		// res.FailArgumentCode 表示参数错误（通常对应 HTTP 400）
		res.FailWithCode(c, res.FailArgumentCode)
		return // 终止函数执行
	}
	nickName := randomname.GenerateName()
	nickNameFirst := string([]rune(nickName)[0])
	avatar := fmt.Sprintf("uploads/chat_avatar/%s.png", nickNameFirst)
	var addr = conn.RemoteAddr().String()
	var User = ChatUser{
		Conn:     conn,
		NickName: nickName,
		Avatar:   avatar,
	}
	// map 记录客户地址 对应的 websocket  有了之后 可以群聊
	ConnGroupMap[addr] = User

	logrus.Infof("%s 连接成功", addr)
	// 循环读取客户端发送的消息（WebSocket 连接是长连接，需要持续监听）
	for {
		// 调用 conn.ReadMessage 读取客户端消息
		// 返回值说明：
		// 第一个返回值：消息类型（如 websocket.TextMessage 文本消息、BinaryMessage 二进制消息等）
		// p：消息内容（字节数组）
		// err：读取过程中的错误（如客户端断开连接会返回错误）
		_, p, err := conn.ReadMessage()

		// 如果读取错误（如客户端断开连接），退出循环
		if err != nil {
			// 用户断开聊天（日志说明）
			SendGroupMsg(conn, GroupResponse{
				MsgType:     OutRoomMsg,
				NickName:    User.NickName,
				Avatar:      User.Avatar,
				Content:     fmt.Sprintf("%s 离开聊天室", addr),
				OnlineCount: len(ConnGroupMap) - 1,
				Date:        time.Now(),
			})

			break
		}
		var request GroupRequest
		json.Unmarshal(p, &request)

		switch request.MsgType {
		case TextMsg:
			if strings.TrimSpace(request.Content) == "" {
				SendMsg(addr, GroupResponse{
					Avatar:      User.Avatar,
					MsgType:     SystemMsg,
					Content:     "消息不能为空",
					NickName:    User.NickName,
					OnlineCount: len(ConnGroupMap),
					Date:        time.Now(),
				})
				continue
			}
			// 将接收到的消息字节数组转换为字符串并发送给人
			SendGroupMsg(conn, GroupResponse{
				Content:     request.Content,
				MsgType:     TextMsg,
				NickName:    User.NickName,
				Avatar:      User.Avatar,
				OnlineCount: len(ConnGroupMap),
				Date:        time.Now(),
			})
		case InRoomMsg:
			SendGroupMsg(conn, GroupResponse{
				NickName:    User.NickName,
				Avatar:      User.Avatar,
				Content:     fmt.Sprintf("%s 进入聊天室", User.NickName),
				OnlineCount: len(ConnGroupMap),
				Date:        time.Now(),
			})
		default:
			SendMsg(addr, GroupResponse{
				Content:     "格式错误",
				MsgType:     SystemMsg,
				NickName:    User.NickName,
				Avatar:      User.Avatar,
				OnlineCount: len(ConnGroupMap),
				Date:        time.Now(),
			})
		}

		// 向客户端发送消息
		// 参数说明：
		// 第一个参数：消息类型（websocket.TextMessage 表示文本消息）
		// 第二个参数：消息内容（字节数组，这里发送固定字符串 "xxx"）
		//conn.WriteMessage(websocket.TextMessage, []byte("xxx"))
	}

	// 延迟关闭 WebSocket 连接（在函数退出前执行）
	// 注意：由于上面的 for 循环退出后函数会结束，因此 defer 会在此处执行
	defer conn.Close()
	delete(ConnGroupMap, addr)
}

func SendGroupMsg(conn *websocket.Conn, response GroupResponse) {
	byteData, _ := json.Marshal(response)
	addr := conn.RemoteAddr().String()
	ip, address := GetIpAndAddr(addr)
	global.DB.Create(&models.ChatModel{
		NickName: response.NickName,
		Avatar:   response.Avatar,
		Content:  response.Content,
		IP:       ip,
		Addr:     address,
		IsGroup:  true,
		MsgType:  response.MsgType,
	})
	for _, User := range ConnGroupMap {
		// 服务器 向客户端发送消息
		User.Conn.WriteMessage(websocket.TextMessage, byteData)
	}

}

func SendMsg(addr string, response GroupResponse) {
	byteData, _ := json.Marshal(response)
	user := ConnGroupMap[addr]
	ip, address := GetIpAndAddr(addr)
	user.Conn.WriteMessage(websocket.TextMessage, byteData)
	global.DB.Create(&models.ChatModel{
		NickName: response.NickName,
		Avatar:   response.Avatar,
		Content:  response.Content,
		IP:       ip,
		Addr:     address,
		IsGroup:  false,
		MsgType:  response.MsgType,
	})

}

func GetIpAndAddr(addr string) (string, string) {
	addrList := strings.Split(addr, ":")
	address := core.GetIpAddr(addrList[0])
	return addrList[0], address
}

func (ChatApi) ChatListView(c *gin.Context) {
	var cr common.PageInfo
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	list, count, err := common.ListQuery(models.ChatModel{
		IsGroup: true,
	}, common.Options{
		PageInfo:     cr,
		DefaultOrder: "created_at desc", //默认
	})

	data := filter.Omit("list", list)
	_list, _ := data.(filter.Filter)
	if string(_list.MustMarshalJSON()) == "{}" {
		list := make([]models.AdvertModel, 0)
		res.OkWithList(c, list, 0)

	}
	res.OkWithList(c, data, count)

}
