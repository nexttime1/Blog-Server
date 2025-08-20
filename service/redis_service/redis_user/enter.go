package redis_user

import (
	"Blog_server/global"
	"github.com/sirupsen/logrus"
)

const UserDiggPrefix = "user_digg"

type UserDigg struct {
	Index string
}

func NewUserDigg() UserDigg {
	return UserDigg{
		Index: UserDiggPrefix,
	}
}

func (UserDigg) Add(UserID string, ArticleID ...interface{}) {
	global.Redis.RPush(UserID, ArticleID...)
}

func (UserDigg) Get(UserID string) []string {
	result, err := global.Redis.LRange(UserID, 0, -1).Result()
	if err != nil {
		logrus.Errorf("%v", err)
	}

	return result
}

func (UserDigg) Del(UserID string, ArticleID string) {
	err := global.Redis.LRem(UserID, 0, ArticleID).Err()
	if err != nil {
		logrus.Errorf("%v", err)
	}
}

func (UserDigg) DelAll(UserID string) {
	err := global.Redis.Del(UserID).Err()
	if err != nil {
		logrus.Errorf("%v", err)
	}
}

func Exist(list []string, ArticleID string) bool {
	// 判断 id 是否不在列表中
	notInList := false
	for _, item := range list {
		if item == ArticleID {
			notInList = true
			break // 找到匹配项，跳出循环
		}
	}
	return notInList
}
