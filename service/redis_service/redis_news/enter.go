package redis_news

import (
	"Blog_server/global"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

const NewPrefix = "news_index"

type NewData = any

func SetNew(key string, NewData []NewData) {
	byteData, _ := json.Marshal(NewData)
	err := global.Redis.Set(fmt.Sprintf("%s_%s", NewPrefix, key), byteData, 1*time.Hour).Err()
	if err != nil {
		logrus.Errorf("Set 失败%s", err)
	}
}

func GetNews(key string) (newData []NewData, err error) {
	res := global.Redis.Get(fmt.Sprintf("%s_%s", NewPrefix, key)).Val()
	err = json.Unmarshal([]byte(res), &newData)
	return
}
