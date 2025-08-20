package redis_news

import (
	"Blog_server/global"
	"Blog_server/service/new_service"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

const NewPrefix = "news_index"

func SetNew(key string, NewData []new_service.NewData) {
	byteData, _ := json.Marshal(NewData)
	err := global.Redis.Set(fmt.Sprintf("%s_%s", NewPrefix, key), byteData, 1*time.Hour).Err()
	if err != nil {
		logrus.Errorf("Set 失败%s", err)
	}
}

func GetNews(key string) (newData []new_service.NewData, err error) {
	res := global.Redis.Get(fmt.Sprintf("%s_%s", NewPrefix, key)).Val()
	err = json.Unmarshal([]byte(res), &newData)
	return
}
