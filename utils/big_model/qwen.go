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

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Input struct {
	Messages []Message `json:"messages"`
}
type Parameters struct {
	IncrementalOutput bool `json:"incremental_output"` // 是否增量输出
}

type Request struct {
	Model      string     `json:"model"`
	Input      Input      `json:"input"`
	Parameters Parameters `json:"parameters"`
}

type Response struct {
	Output struct {
		FinishReason string `json:"finish_reason"`
		Text         string `json:"text"`
	} `json:"output"`
	Usage struct {
		TotalTokens  int `json:"total_tokens"`
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	RequestId string `json:"request_id"`
}

type QwenModel struct {
	SessionID uint `json:"session_id"`
}

func (q QwenModel) Send(Content string) (msgChan chan string, err error) {
	req := Request{
		Model: "qwen-turbo",
		Input: Input{
			Messages: []Message{},
		},
		Parameters: Parameters{
			IncrementalOutput: true,
		},
	}
	// 找到提示词
	var sessionModel models.BigModelSessionModel
	global.DB.Where("id = ?", q.SessionID).Preload("RoleModel").Find(&sessionModel)
	req.Input.Messages = append(req.Input.Messages, Message{
		Role:    "system",
		Content: sessionModel.RoleModel.Prompt,
	})

	// 找到聊天记录  sessionID 一定是对的  不然到不了这里
	var ChatList []models.BigModelChatModel
	global.DB.Where("session_id = ?", q.SessionID).Order("created_at asc").Find(&ChatList)
	for _, model := range ChatList {
		req.Input.Messages = append(req.Input.Messages, Message{
			Role:    "user",
			Content: model.Content,
		}, Message{
			Role:    "assistant",
			Content: model.BotContent,
		})
	}
	// 这是 我这次要问的东西
	req.Input.Messages = append(req.Input.Messages, Message{
		Role:    "user",
		Content: Content,
	})

	msgChan = make(chan string, 0)
	baseUrl := "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
	byteData, _ := json.Marshal(req)
	buf := bytes.NewBuffer(byteData)

	request, err := http.NewRequest("POST", baseUrl, buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	apiKey := "sk-e8ed4fd4fc82489f8494ef9b6317e9e8"
	request.Header.Add("Authorization", "Bearer "+apiKey)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-DashScope-SSE", "enable")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	scan := bufio.NewScanner(response.Body) // 分片读
	scan.Split(bufio.ScanLines)             // 按行读取
	go func() {
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
			}
		}
	}()
	return

}
