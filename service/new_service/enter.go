package new_service

import (
	"Blog_server/service/redis_service/redis_news"
	"Blog_server/utils/request"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

type Params struct {
	ID   string `json:"id"`
	Size int    `json:"size"`
}

type Header struct {
	Signaturekey string `form:"signaturekey" structs:"signaturekey"`
	Version      string `form:"version" structs:"version"`
	UserAgent    string `form:"User-Agent" structs:"User-Agent"`
}

type NewResponse struct {
	Code int                  `json:"code"`
	Data []redis_news.NewData `json:"data"`
	Msg  string               `json:"msg"`
}

func NewListService(cr Params, newAPI string, timeout time.Duration) ([]redis_news.NewData, error) {
	if cr.Size == 0 {
		cr.Size = 1
	}
	key := fmt.Sprintf("%s-%d", cr.ID, cr.Size)
	data, _ := redis_news.GetNews(key)
	if len(data) != 0 {
		//缓存有
		return data, nil

	}
	httpResponse, err := request.Get(newAPI, timeout)
	if err != nil {
		logrus.Errorf("get 请求错误：%s", err)
		return []redis_news.NewData{}, err
	}
	defer httpResponse.Body.Close()
	var response NewResponse

	byteData, err := io.ReadAll(httpResponse.Body)
	err = json.Unmarshal(byteData, &response)
	if err != nil {
		logrus.Errorf("json解析：%s", err)
		return []redis_news.NewData{}, err
	}
	if response.Code != 200 {
		logrus.Errorf("状态码错误：%d", response.Code)
		return []redis_news.NewData{}, errors.New(response.Msg)
	}
	redis_news.SetNew(key, response.Data)
	return response.Data, nil
}
