package request

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"time"
)

func Get(urlString string, timeout time.Duration) (resp *http.Response, err error) {
	// 1. 解析原始URL，准备添加查询参数
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, err // 解析URL失败时返回错误
	}
	queryParams := parsedURL.Query()

	// 3. 添加需要的查询参数：lang=cn 和 id=KqndgxeLl9
	queryParams.Add("lang", "cn")
	queryParams.Add("id", "KqndgxeLl9")

	// 4. 将更新后的查询参数重新设置到URL中
	parsedURL.RawQuery = queryParams.Encode()

	// 5. 使用带查询参数的URL创建GET请求（注意第三个参数为nil，因为GET请求无body）
	httpReq, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		logrus.Errorf("NewRequest error:%s", err.Error())
		return nil, err

	}
	httpReq.Header.Set("User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 SLBrowser/9.0.6.8151 SLBChan/105 SLBVPV/64-bit")

	client := http.Client{
		Timeout: timeout,
	}
	//发送Get 得到结果
	httpResp, err := client.Do(httpReq)

	if err != nil {
		logrus.Errorf("client.Do error: %s", err.Error())
		return nil, err
	}

	// 检查 HTTP 状态码
	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		logrus.Errorf("HTTP status error: code=%d, response=%s", httpResp.StatusCode, string(body))
		return httpResp, fmt.Errorf("http status code: %d", httpResp.StatusCode)
	}

	return httpResp, err
}
