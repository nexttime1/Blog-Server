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
	"github.com/sirupsen/logrus"
	"io"
	"reflect"
	"strings"
)

type ActionLog struct {
	c            *gin.Context
	level        enum.LevelType
	title        string
	content      string
	requestBody  []byte
	responseBody []byte
	log          *models.LogModel //备份
	showRequest  bool
	showResponse bool
	itemList     []string
}

// SetLink 设置超链接
func (ac *ActionLog) SetLink(label string, href string) string {
	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_item link\"><div class=\"log_item_label\">%s</div><div class=\"log_item_content\"><a href=\"%s\" target=\"_blank\">%s</a> </div></div>",
		label, href, href))
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

func (ac *ActionLog) Save() {
	if ac.log != nil {
		//说明已经存在  更新一下
		global.DB.Model(ac.log).Updates(map[string]interface{}{
			"title": "更新",
		})
		return
	}

	var newItemList []string

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

	//设置响应
	if ac.showResponse {
		newItemList = append(newItemList, fmt.Sprintf("<div class=\"log_response\"><pre class=\"log_json_body\">%s</pre></div>",
			string(ac.responseBody),
		))
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
		return
	}
	ac.log = &log

}

func NewActionLogByGin(c *gin.Context) *ActionLog {
	return &ActionLog{c: c}
}

func GetLog(c *gin.Context) *ActionLog {
	_log, exists := c.Get("log")
	if !exists {
		return NewActionLogByGin(c)
	}
	log, ok := _log.(*ActionLog)
	if !ok {
		return NewActionLogByGin(c)
	}
	return log

}
