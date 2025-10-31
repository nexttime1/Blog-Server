package redis_count

import (
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/utils/es"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"

	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"strconv"
)

const (
	CommentPrefix     = "article_comment"
	DiggPrefix        = "article_digg"
	LookPrefix        = "article_look"
	CommentDiggPrefix = "comment_digg"
)

type CountDB struct {
	Index string
}

func NewDigg() CountDB {
	return CountDB{
		Index: DiggPrefix,
	}
}
func NewLook() CountDB {
	return CountDB{
		Index: LookPrefix,
	}
}
func NewComment() CountDB {
	return CountDB{
		Index: CommentPrefix,
	}
}

func NewCommentDigg() CountDB {
	return CountDB{
		Index: CommentDiggPrefix,
	}
}

func (c CountDB) Set(id string) {
	num, err := global.Redis.HGet(c.Index, id).Int()
	if err != nil {
		// 区分 "键不存在" 和其他错误（如连接失败）
		if errors.Is(err, redis.Nil) {
			logrus.Infof("键 %s 不存在，初始化为 0", id)
			num = 0 // 明确初始化为 0
		} else {
			logrus.Errorf("HGet 错误: %s", err)
			return // 非 "键不存在" 的错误，直接返回避免后续错误
		}
	}

	num++
	err = global.Redis.HSet(c.Index, id, num).Err()
	if err != nil {
		logrus.Errorf("HSet 失败: %s", err)
	} else {
		logrus.Infof("HSet 成功，id=%s, 新值=%d", id, num)
	}

}

func (c CountDB) Get(id string) int {
	num, err := global.Redis.HGet(c.Index, id).Int()
	if err != nil {
		logrus.Errorf("查找错误 %s", err)
	}

	return num
}

func (c CountDB) SetNum(id string, NewNum int) {
	OldNum, err := global.Redis.HGet(c.Index, id).Int()
	if err != nil {
		logrus.Errorf("查找错误 %s", err)
	}

	num := OldNum + NewNum
	global.Redis.HSet(c.Index, id, num)
}

func (c CountDB) GetInfo() map[string]int {
	var getInfo = make(map[string]int)
	data := global.Redis.HGetAll(c.Index).Val()
	for id, val := range data {
		num, _ := strconv.Atoi(val)
		getInfo[id] = num
	}
	return getInfo
}

func (c CountDB) Clear() {
	global.Redis.Del(c.Index)
}

func Update() {
	ctx := context.Background()
	result, err := global.Es.Search(models.ArticleModel{}.Index()).
		Query(elastic.NewMatchAllQuery()).
		Size(10000).
		Do(context.Background())
	if err != nil {
		logrus.Errorf("查询失败 %s", err)
	}
	DiggMap := NewDigg().GetInfo()
	LookMap := NewLook().GetInfo()
	CommentMap := NewComment().GetInfo()
	for _, hit := range result.Hits.Hits {
		var article models.ArticleModel
		var MapEnd = make(map[string]int)
		_ = json.Unmarshal(hit.Source, &article)
		// 将redis 存的值 放进去
		diggCount := DiggMap[article.ID]
		newDiggCount := article.DiggCount + diggCount // 也就是 没更新的 + reids 更新的   加完 redis 清空 这样 没必要占着 空间
		if diggCount != 0 {                           //也就是  缓存 不为0  也就需要更新
			MapEnd["digg_count"] = newDiggCount
		}

		lookCount := LookMap[article.ID]
		newLookCount := article.LookCount + lookCount
		if lookCount != 0 {
			MapEnd["look_count"] = newLookCount
		}
		commentCount := CommentMap[article.ID]
		newCommentCount := article.CommentCount + commentCount
		if commentCount != 0 {
			MapEnd["comment_count"] = newCommentCount
		}

		if len(MapEnd) == 0 {
			//logrus.Infof("%s 无变化", article.Title)
			continue
		}
		//需要更新
		_, err := es.SafeESUpdate(ctx, global.Es, models.ArticleModel{}.Index(), hit.Id, MapEnd)
		if err != nil {
			logrus.Errorf("文章《%s》ES 更新失败: %v", article.Title, err)
			continue
		}

		logrus.Infof("%s 更新成功", article.Title)
	}

	NewDigg().Clear()    //清除缓存
	NewLook().Clear()    //清除缓存
	NewComment().Clear() //清除缓存
}

func UpdateToDB() {
	infoList := NewCommentDigg().GetInfo()
	var modelList []models.CommentModel
	global.DB.Find(&modelList)
	for _, model := range modelList {
		RedisCount := infoList[fmt.Sprintf("%d", model.ID)]
		if RedisCount == 0 {
			//logrus.Infof("%s 不需要更新评论点赞", model.Content[:10])
			continue
		}
		global.DB.Model(&model).Update("digg_count", model.DiggCount+RedisCount)
		logrus.Infof("%s  该评论更新点赞成功", safeSubstr(model.Content, 10))

	}
	NewCommentDigg().Clear()
}

func safeSubstr(s string, length int) string {
	if len(s) < length {
		return s
	}
	return s[:length]
}
