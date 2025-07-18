package log_service

import (
	"Blog_server/core"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	e "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type ActionLog struct {
	c                  *gin.Context
	level              enum.LevelType
	title              string
	content            string
	requestBody        []byte
	responseBody       []byte
	log                *models.LogModel //备份
	showRequestHeader  bool
	showRequest        bool
	showResponseHeader bool
	showResponse       bool
	itemList           []string
	ResponseHeader     http.Header
	IsMiddenWare       bool
}

func (ac *ActionLog) SetError(label string, err error) {
	msg := e.WithStack(err)
	logrus.Errorf("%s  %s", label, err.Error())
	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_error\"><div class=\"line\"><div class=\"label\">%s</div><div class=\"value\">%s</div><div class=\"type\">%T</div></div><div class=\"stack\">%+v</div></div>\n",
		label, err, err, msg))

}

// ShowRequestHeader 显示请求头
func (ac *ActionLog) ShowRequestHeader() {
	ac.showRequestHeader = true
}

// ShowResponseHeader 显示响应头
func (ac *ActionLog) ShowResponseHeader() {
	ac.showResponseHeader = true
}

// SetLink 设置超链接
func (ac *ActionLog) SetLink(label string, href string) {
	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_item link\"><div class=\"log_item_label\">%s</div><div class=\"log_item_content\"><a href=\"%s\" target=\"_blank\">%s</a> </div></div>",
		label, href, href))

}

func (ac *ActionLog) SetImage(src string) {
	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_image\"><img src=\"%s\" alt=\"\" ></div>", src))
}

func (ac *ActionLog) setItem(label string, value any, levelType enum.LevelType) {
	var v string
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Map, reflect.Slice, reflect.Struct:
		ByteData, _ := json.Marshal(value)
		v = string(ByteData)
	default:
		v = fmt.Sprintf("%v", value)
	}

	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_item %s\"><div class=\"log_item_label\">%s</div><div class=\"log_item_content\">%s</div></div>",
		levelType, label, v))
}

func (ac *ActionLog) SetItem(label string, value any) {
	ac.setItem(label, value, enum.LogInfoLevel)
}

func (ac *ActionLog) SetItemInfo(label string, value any) {
	ac.setItem(label, value, enum.LogInfoLevel)
}
func (ac *ActionLog) SetItemWarn(label string, value any) {
	ac.setItem(label, value, enum.LogWainLevel)
}
func (ac *ActionLog) SetItemError(label string, value any) {
	ac.setItem(label, value, enum.LogErrLevel)

}

func (ac *ActionLog) ShowRequest() {
	ac.showRequest = true
}

func (ac *ActionLog) ShowResponse() {
	ac.showResponse = true
}

func (ac *ActionLog) SetTitle(title string) {
	ac.title = title
}

func (ac *ActionLog) SetLevel(level enum.LevelType) {
	ac.level = level
}

func (ac *ActionLog) SetRequest(c *gin.Context) {
	ByteData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Errorf(err.Error())
	}
	fmt.Println("Body", string(ByteData))
	c.Request.Body = io.NopCloser(bytes.NewReader(ByteData))
	ac.requestBody = ByteData
}

func (ac *ActionLog) SetResponse(data []byte) {
	ac.responseBody = data
}

func (ac *ActionLog) SetResponseHeader(header http.Header) {
	ac.ResponseHeader = header
}

func (ac *ActionLog) MiddleSave() {
	value, _ := ac.c.Get("SaveLog")
	b, _ := value.(bool)
	if !b { //如果 b是false 说明 我没有调用 该Save的函数  就不要保存 操作日志
		return
	}
	if ac.log == nil {
		ac.IsMiddenWare = true
		ac.Save()
	}
	//如果 调用 中间件响应  之前 save 过
	//响应头
	if ac.showResponseHeader {

		ByteData, _ := json.Marshal(ac.ResponseHeader)
		fmt.Println("showResponseHeader  ", string(ByteData))
		ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_response_header\"><pre class=\"log_json_body\">%s</pre></div>", string(ByteData)))
	}

	//设置响应
	if ac.showResponse {
		ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_response\"><pre class=\"log_json_body\">%s</pre></div>",
			string(ac.responseBody),
		))
	}
	ac.Save()
}

func (ac *ActionLog) Save() uint {
	if ac.log != nil {
		//说明已经存在  更新一下
		NewItemList := strings.Join(ac.itemList, "\n")
		content := ac.log.Content + "\n" + NewItemList

		global.DB.Model(ac.log).Updates(map[string]interface{}{
			"content": content,
		})
		return ac.log.ID
	}
	var newItemList []string
	//请求头
	if ac.showRequestHeader {
		//fmt.Println("ac.itemList\n", ac.itemList)   //先 中间件请求部分   再是Api Api最后才调用Save方法   在是中间件响应
		ByteData, _ := json.Marshal(ac.c.Request.Header)
		newItemList = append(newItemList, fmt.Sprintf("<div class=\"log_request_header\"><pre class=\"log_json_body\">%s</pre></div>", string(ByteData)))
	}

	//设置请求
	if ac.showRequest {
		newItemList = append(newItemList, fmt.Sprintf("<div class=\"log_request\"><div class=\"log_request_head\"><span class=\"log_request_method %s\">%s</span><span class=\"log_request_path\">%s</span></div><div class=\"log_request_body\"><pre class=\"log_json_body\">%s</pre></div></div>",
			strings.ToLower(ac.c.Request.Method),
			ac.c.Request.Method,
			ac.c.Request.URL.String(),
			string(ac.requestBody),
		))
	}

	//设置 content
	newItemList = append(newItemList, ac.itemList...)

	if ac.IsMiddenWare { //只有是 走到中间件响应才走这个 要不 没有
		//响应头
		if ac.showResponseHeader {

			ByteData, _ := json.Marshal(ac.ResponseHeader)
			fmt.Println("showResponseHeader  ", string(ByteData))
			newItemList = append(newItemList, fmt.Sprintf("<div class=\"log_response_header\"><pre class=\"log_json_body\">%s</pre></div>", string(ByteData)))
		}

		//设置响应
		if ac.showResponse {
			newItemList = append(newItemList, fmt.Sprintf("<div class=\"log_response\"><pre class=\"log_json_body\">%s</pre></div>",
				string(ac.responseBody),
			))
		}
		//清空  如果第二次调用Save   ac.itemList我只希望有尾部
		ac.itemList = []string{}
	}

	ip := ac.c.ClientIP()
	addr := core.GetIpAddr(ip)
	UserID := uint(0)
	log := models.LogModel{
		LogType: enum.ActionLogType,
		Title:   ac.title,
		Content: strings.Join(newItemList, "\n"), //按照换行去合并
		Level:   ac.level,
		UserID:  UserID,
		IP:      ip,
		Addr:    addr,
	}
	err := global.DB.Create(&log).Error
	if err != nil {
		logrus.Errorf("日志创建失败 %s", err.Error())
		return 0
	}
	ac.log = &log
	return ac.log.ID
}

func NewActionLogByGin(c *gin.Context) *ActionLog {
	return &ActionLog{c: c}
}

func GetLog(c *gin.Context) *ActionLog {
	//c.Get的功能是从当前请求的上下文里获取之前存进去的值。
	_log, exists := c.Get("log")
	if !exists {
		return NewActionLogByGin(c)
	}
	log, ok := _log.(*ActionLog)
	if !ok {
		return NewActionLogByGin(c)
	}
	c.Set("SaveLog", true) //这样如果没有调用该Save的函数  就可以不去走 Save方法  因为中间件必走 最后会调用Save
	return log

}
