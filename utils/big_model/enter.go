package big_model

import (
	"Blog_server/global"
	"errors"
)

type BigModelInterface interface {
	Send(Content string) (msgChan chan string, err error)
}

func Send(SessionID uint, Content string) (msgChan chan string, err error) {
	switch global.Config.BigModel.Setting.Name {
	case "qwen":
		return QwenModel{SessionID: SessionID}.Send(Content)
	case "xtm":
		return XModel{}.Send(Content)
	default:
		return make(chan string), errors.New("不支持该大模型")
	}

}
