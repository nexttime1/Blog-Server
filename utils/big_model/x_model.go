package big_model

import (
	"Blog_server/global"
	"Blog_server/models"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type XModel struct {
	SessionID uint `json:"session_id"`
}

func (x XModel) Send(Content string) (msgChan chan string, err error) {
	req := Request{
		Model: "qwen-go-custom", // 自定义模型标识
		Input: Input{
			Messages: []Message{},
		},
		Parameters: Parameters{
			IncrementalOutput: true,
		},
	}

	// 找到聊天记录（注意：不需要添加系统提示词，Python端会自动添加）
	var ChatList []models.BigModelChatModel
	global.DB.Where("session_id = ?", x.SessionID).Order("created_at asc").Find(&ChatList)
	for _, model := range ChatList {
		req.Input.Messages = append(req.Input.Messages, Message{
			Role:    "user",
			Content: model.Content,
		}, Message{
			Role:    "assistant",
			Content: model.BotContent,
		})
	}

	// 这次要问的内容
	req.Input.Messages = append(req.Input.Messages, Message{
		Role:    "user",
		Content: Content,
	})

	msgChan = make(chan string, 0)
	baseUrl := "http://172.21.13.113:8000/api/v1/services/aigc/text-generation/generation"
	byteData, _ := json.Marshal(req)
	buf := bytes.NewBuffer(byteData)

	request, err := http.NewRequest("POST", baseUrl, buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	scan := bufio.NewScanner(response.Body)
	scan.Split(bufio.ScanLines)

	go func() {
		defer response.Body.Close()
		for scan.Scan() {
			text := scan.Text()
			if text == "" ||
				strings.HasPrefix(text, "id") ||
				strings.HasPrefix(text, "event:") ||
				strings.HasPrefix(text, ":HTTP_STATUS") {
				continue
			}

			var res Response
			err = json.Unmarshal([]byte(text[5:]), &res)
			if err != nil {
				fmt.Println(err, text[5:])
				continue
			}

			msgChan <- res.Output.Text
			if res.Output.FinishReason == "stop" {
				close(msgChan)
				break
			}
		}
	}()

	return
}
