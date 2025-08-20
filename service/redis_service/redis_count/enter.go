package redis_count

import (
	"Blog_server/global"
	"Blog_server/models"
	"context"
	"encoding/json"
	"fmt"
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
		logrus.Errorf("查找错误 %s", err)
	}

	num++
	global.Redis.HSet(c.Index, id, num)
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

func (c CountDB) Update() error {
	var MapEnd = make(map[string]int)
	result, err := global.Es.Search(models.ArticleModel{}.Index()).
		Query(elastic.NewMatchAllQuery()).
		Size(1000).
		Do(context.Background())
	if err != nil {
		logrus.Errorf("查询失败 %s", err)
		return fmt.Errorf("查询失败 %s", err)
	}
	DiggMap := NewDigg().GetInfo()
	LookMap := NewLook().GetInfo()
	CommentMap := NewComment().GetInfo()
	for _, hit := range result.Hits.Hits {
		var article models.ArticleModel
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
		newCommentCount := article.LookCount + commentCount
		if commentCount != 0 {
			MapEnd["comment_count"] = newCommentCount
		}

		if len(MapEnd) == 0 {
			continue
		}
		//需要更新
		_, err = global.Es.
			Update().
			Index(models.ArticleModel{}.Index()).
			Id(hit.Id).
			Doc(MapEnd).Do(context.Background())
		if err != nil {
			logrus.Errorf("es id为 %s 更新失败", hit.Id, err)
		}
	}
	logrus.Infof("更新成功")
	NewDigg().Clear()    //清除缓存
	NewLook().Clear()    //清除缓存
	NewComment().Clear() //清除缓存
	return nil
}
