package redis_digg

import (
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/redis_service/redis_comment"
	"Blog_server/service/redis_service/redis_look"
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"strconv"
)

const DiggPrefix = "digg"

func Digging(id string) {
	num, err := global.Redis.HGet(DiggPrefix, id).Int()
	fmt.Errorf("%s", err)
	//没有的话  num = 0
	num++
	global.Redis.HSet(DiggPrefix, id, num)

}

func GetDigging(id string) int {
	num, _ := global.Redis.HGet(DiggPrefix, id).Int()
	return num
}

func GetDiggingInfo() map[string]int {
	var DiggInfo = make(map[string]int)
	data := global.Redis.HGetAll(DiggPrefix).Val()
	for id, val := range data {
		num, _ := strconv.Atoi(val)
		DiggInfo[id] = num
	}
	return DiggInfo
}

func DiggClear() {
	global.Redis.Del(DiggPrefix)
}

func SyncEsArticle() error {
	var MapEnd = make(map[string]int)
	result, err := global.Es.Search(models.ArticleModel{}.Index()).
		Query(elastic.NewMatchAllQuery()).
		Size(1000).
		Do(context.Background())
	if err != nil {
		logrus.Errorf("查询失败 %s", err)
		return fmt.Errorf("查询失败 %s", err)
	}
	DiggMap := GetDiggingInfo()
	LookMap := redis_look.GetLookInfo()
	CommentMap := redis_comment.GetCommentInfo()
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
	DiggClear() //清除缓存
	redis_look.LookClear()
	return nil
}
