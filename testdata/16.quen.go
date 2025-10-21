package main

//type Message struct {
//	Role    string `json:"role"`
//	Content string `json:"content"`
//}
//type RequestBody struct {
//	Model    string    `json:"model"`
//	Messages []Message `json:"messages"`
//}
//
//type Choice struct {
//	Message      Message
//	FinishReason string      `json:"finish_reason"`
//	Index        int         `json:"index"`
//	Logprobs     interface{} `json:"logprobs"`
//}
//
//type Usage struct {
//	PromptTokens        int `json:"prompt_tokens"`
//	CompletionTokens    int `json:"completion_tokens"`
//	TotalTokens         int `json:"total_tokens"`
//	PromptTokensDetails struct {
//		CachedTokens int `json:"cached_tokens"`
//	} `json:"prompt_tokens_details"`
//}
//type QuenResponse struct {
//	Choices           []Choice `json:"choices"`
//	Object            string   `json:"object"`
//	Usage             `json:"usage"`
//	Created           int         `json:"created"`
//	SystemFingerprint interface{} `json:"system_fingerprint"`
//	Model             string      `json:"model"`
//	Id                string      `json:"id"`
//}
//
//func Send(requestBody RequestBody) (error, chan string) {
//	jsonData, err := json.Marshal(requestBody)
//	if err != nil {
//		log.Fatal(err)
//	}
//	client := &http.Client{}
//	TextChan := make(chan string, 0)
//	// 创建 POST 请求
//	// 以下是北京地域base_url，如果使用新加坡地域的模型，需要将base_url替换为：https://dashscope-intl.aliyuncs.com/compatible-mode/v1/chat/completions
//	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions", bytes.NewBuffer(jsonData))
//	if err != nil {
//		return err, TextChan
//	}
//	// 设置请求头
//	// 若没有配置环境变量，请用阿里云百炼API Key将下行替换为：apiKey := "sk-xxx"
//	// 新加坡和北京地域的API Key不同。获取API Key：https://help.aliyun.com/zh/model-studio/get-api-key
//	apiKey := os.Getenv("DASHSCOPE_API_KEY")
//	apiKey = "sk-e8ed4fd4fc82489f8494ef9b6317e9e8"
//	req.Header.Set("Authorization", "Bearer "+apiKey)
//	req.Header.Set("Content-Type", "application/json")
//	// 发送请求
//	resp, err := client.Do(req)
//	if err != nil {
//		return err, TextChan
//	}
//	defer resp.Body.Close()
//	// 读取响应体
//	bodyText, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return err, TextChan
//	}
//	// 打印响应内容
//	fmt.Printf("%s\n", bodyText)
//	var response QuenResponse
//	err = json.Unmarshal(bodyText, &response)
//	if err != nil {
//		return err, TextChan
//
//	}
//	go func() {
//		for _, choice := range response.Choices {
//			TextChan <- choice.Message.Content
//			if choice.FinishReason == "stop" {
//				close(TextChan)
//			}
//		}
//	}()
//	return nil, TextChan
//}
//
//func main() {
//	// 创建 HTTP 客户端
//
//	// 构建请求体
//	requestBody := RequestBody{
//		// 模型列表：https://help.aliyun.com/zh/model-studio/getting-started/models
//		Model: "qwen-plus",
//		Messages: []Message{
//			{
//				Role:    "system",
//				Content: "You are a helpful assistant.",
//			},
//			{
//				Role:    "user",
//				Content: "我希望你用10个字夸我一下，不能多？",
//			},
//			{
//				Role:    "user",
//				Content: "简短介绍一下北京工业大学",
//			},
//		},
//	}
//	err, TextChan := Send(requestBody)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for text := range TextChan {
//		fmt.Println(text)
//	}
//}
