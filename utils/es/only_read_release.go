package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

// AutoFixReadOnlyIfNeeded 尝试解除 ES 索引的只读状态
func AutoFixReadOnlyIfNeeded() error {
	esURL := "http://localhost:9200" // 当前的 ES 地址
	body := []byte(`{"index.blocks.read_only_allow_delete": null}`)

	req, err := http.NewRequest("PUT", esURL+"/_all/_settings", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("构建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// ✅ 如果 ES 开启了用户名密码认证（像你用 curl -u elastic:123456）
	req.SetBasicAuth("elastic", "123456")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送解除只读请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("解除只读失败，HTTP 状态码: %d", resp.StatusCode)
	}

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		logrus.Warnf("解析解除只读返回失败: %v", err)
	}
	logrus.Infof("已尝试解除 ES 只读限制: %v", res)
	return nil
}

// SafeESUpdate 尝试更新 ES，如果发现只读则自动解除并重试
func SafeESUpdate(ctx context.Context, client *elastic.Client, index, id string, doc interface{}) (*elastic.UpdateResponse, error) {
	resp, err := client.Update().Index(index).Id(id).Doc(doc).Do(ctx)
	if err == nil {
		return resp, nil
	}

	if strings.Contains(err.Error(), "read-only-allow-delete") ||
		strings.Contains(err.Error(), "FORBIDDEN/12") ||
		strings.Contains(err.Error(), "TOO_MANY_REQUESTS") {

		logrus.Warn("[ES] 检测到索引只读，尝试自动解除并重试一次写入")
		if fixErr := AutoFixReadOnlyIfNeeded(); fixErr != nil {
			logrus.Errorf("[ES] 自动解除只读失败: %v", fixErr)
			return nil, fmt.Errorf("原始错误: %v；自动解除只读失败: %v", err, fixErr)
		}

		time.Sleep(500 * time.Millisecond)
		resp2, err2 := client.Update().Index(index).Id(id).Doc(doc).Do(ctx)
		if err2 != nil {
			return nil, fmt.Errorf("重试写入失败: %v (原始错误: %v)", err2, err)
		}
		return resp2, nil
	}

	return nil, err
}
