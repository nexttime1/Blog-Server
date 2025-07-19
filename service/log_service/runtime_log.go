package log_service

import (
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"encoding/json"
	"fmt"
	e "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"time"
)

type RuntimeLog struct {
	level           enum.LevelType
	title           string
	ItemList        []string
	ServiceName     string
	RuntimeDateType RuntimeDateType
}

func (rc *RuntimeLog) SetError(label string, err error) {
	msg := e.WithStack(err)
	logrus.Errorf("%s  %s", label, err.Error())
	rc.ItemList = append(rc.ItemList, fmt.Sprintf("<div class=\"log_error\"><div class=\"line\"><div class=\"label\">%s</div><div class=\"value\">%s</div><div class=\"type\">%T</div></div><div class=\"stack\">%+v</div></div>\n",
		label, err, err, msg))

}

// SetLink 设置超链接
func (rc *RuntimeLog) SetLink(label string, href string) {
	rc.ItemList = append(rc.ItemList, fmt.Sprintf("<div class=\"log_item link\"><div class=\"log_item_label\">%s</div><div class=\"log_item_content\"><a href=\"%s\" target=\"_blank\">%s</a> </div></div>",
		label, href, href))

}

func (rc *RuntimeLog) SetImage(src string) {
	rc.ItemList = append(rc.ItemList, fmt.Sprintf("<div class=\"log_image\"><img src=\"%s\" alt=\"\" ></div>", src))
}
func (rc *RuntimeLog) setItem(label string, value any, levelType enum.LevelType) {
	var v string
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Map, reflect.Slice, reflect.Struct:
		ByteData, _ := json.Marshal(value)
		v = string(ByteData)
	default:
		v = fmt.Sprintf("%v", value)
	}

	rc.ItemList = append(rc.ItemList, fmt.Sprintf("<div class=\"log_item %s\"><div class=\"log_item_label\">%s</div><div class=\"log_item_content\">%s</div></div>",
		levelType, label, v))
}

func (rc *RuntimeLog) SetItem(label string, value any) {
	rc.setItem(label, value, enum.LogInfoLevel)
}

func (rc *RuntimeLog) SetItemInfo(label string, value any) {
	rc.setItem(label, value, enum.LogInfoLevel)
}
func (rc *RuntimeLog) SetItemWarn(label string, value any) {
	rc.setItem(label, value, enum.LogWainLevel)
}
func (rc *RuntimeLog) SetItemError(label string, value any) {
	rc.setItem(label, value, enum.LogErrLevel)

}
func (rc *RuntimeLog) SetTitle(title string) {
	rc.title = title
}

func (rc *RuntimeLog) SetLevel(level enum.LevelType) {
	rc.level = level
}

func (rc *RuntimeLog) SetNowTime() {
	rc.ItemList = append(rc.ItemList, fmt.Sprintf("<div class=\"log_time\">%s</div>", time.Now().Format("2006-01-02 15:04:05")))
}

func (rc *RuntimeLog) Save() {
	rc.SetNowTime()
	var log models.LogModel
	global.DB.Debug().Find(&log, fmt.Sprintf("service_name = ? and log_type = ? and created_at >= date_sub(now(), %s)", rc.RuntimeDateType.GetSqlTime()),
		rc.ServiceName, enum.RuntimeLogType)

	if log.ID != 0 { // 找到了 更新
		newContent := strings.Join(rc.ItemList, "\n")
		content := log.Content + "\n" + newContent
		global.DB.Model(&log).Updates(map[string]interface{}{"content": content})
		rc.ItemList = []string{}
		return
	}

	content := strings.Join(rc.ItemList, "\n")
	err := global.DB.Create(&models.LogModel{
		LogType:     enum.RuntimeLogType,
		Title:       rc.title,
		Content:     content,
		Level:       rc.level,
		ServiceName: rc.ServiceName,
	}).Error
	if err != nil {
		logrus.Errorf("创建运行日志错误 :%s", err)
		return
	}
	rc.ItemList = []string{}

}
func (rc RuntimeDateType) GetSqlTime() string {
	switch rc {
	case RuntimeDateTypeHour:
		return "interval 1 hour"
	case RuntimeDateDay:
		return "interval 1 day"
	case RuntimeDateTypeWeek:
		return "interval 1 week"
	case RuntimeDateTypeMonth:
		return "interval 1 month"
	}
	return "interval 1 day" //默认为 一天
}

type RuntimeDateType int8

const (
	RuntimeDateTypeHour  RuntimeDateType = 1
	RuntimeDateDay       RuntimeDateType = 2
	RuntimeDateTypeWeek  RuntimeDateType = 3
	RuntimeDateTypeMonth RuntimeDateType = 4
)

func NewRuntimeLog(serviceName string, dateType RuntimeDateType) *RuntimeLog {
	return &RuntimeLog{
		ServiceName:     serviceName,
		RuntimeDateType: dateType,
	}
}
