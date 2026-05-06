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
	"gorm.io/gorm"
	"sync"
	"time"

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

func (c CountDB) GetAndClear() map[string]int {
	// 创建空map，存最终要返回的计数数据
	var getInfo = make(map[string]int)

	// 生成唯一临时key：原key + :tmp: + 纳秒时间戳
	tmpKey := c.Index + ":tmp:" + strconv.FormatInt(time.Now().UnixNano(), 10)

	// RENAME是原子操作：之后新的点赞写入原key，不影响tmp
	err := global.Redis.Rename(c.Index, tmpKey).Err()
	if err != nil {
		// // 如果原key不存在（没有点赞数据），直接返回空map
		if errors.Is(err, redis.Nil) {
			return getInfo // key不存在，无数据
		}
		logrus.Errorf("Rename失败: %s", err)
		return getInfo
	}

	// 读取tmp数据
	data := global.Redis.HGetAll(tmpKey).Val()
	global.Redis.Del(tmpKey) // 删除临时key

	for id, val := range data {
		num, _ := strconv.Atoi(val)
		getInfo[id] = num
	}
	return getInfo
}

// redis_count/update.go

func Update() {
	ctx := context.Background()

	// ✅ 原子获取并清空，新的点赞会写入原key不受影响
	DiggMap := NewDigg().GetAndClear()
	LookMap := NewLook().GetAndClear()
	CommentMap := NewComment().GetAndClear()

	// 收集所有需要更新的文章ID（只查有变化的）
	articleIDs := mergeKeys(DiggMap, LookMap, CommentMap)
	if len(articleIDs) == 0 {
		return
	}

	logrus.Infof("本轮需要更新 %d 篇文章", len(articleIDs))

	// ✅ 只查询有变化的文章，而不是全量扫描
	result, err := global.Es.Search(models.ArticleModel{}.Index()).
		Query(elastic.NewTermsQueryFromStrings("_id", articleIDs...)).
		Size(len(articleIDs)).
		Do(ctx)
	if err != nil {
		// ✅ 数据已从Redis取出，ES查询失败需要补偿回写
		logrus.Errorf("ES查询失败，执行补偿回写: %s", err)
		compensateBack(DiggMap, LookMap, CommentMap)
		return
	}

	// 并发更新ES（控制并发数）
	var (
		wg      sync.WaitGroup
		sem     = make(chan struct{}, 10) // 最多10个并发
		failIDs []string
		mu      sync.Mutex
	)

	for _, hit := range result.Hits.Hits {
		var article models.ArticleModel
		_ = json.Unmarshal(hit.Source, &article)

		updateMap := buildUpdateMap(article, DiggMap, LookMap, CommentMap)
		if len(updateMap) == 0 {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(hitID string, title string, update map[string]int) {
			defer wg.Done()
			defer func() { <-sem }()

			_, err := es.SafeESUpdate(ctx, global.Es, models.ArticleModel{}.Index(), hitID, update)
			if err != nil {
				logrus.Errorf("文章《%s》ES更新失败: %v", title, err)
				mu.Lock()
				failIDs = append(failIDs, hitID)
				mu.Unlock()
			} else {
				logrus.Infof("文章《%s》更新成功", title)
			}
		}(hit.Id, article.Title, updateMap)
	}
	wg.Wait()

	// ✅ 失败的文章补偿回写Redis
	if len(failIDs) > 0 {
		logrus.Warnf("%d 篇文章更新失败，补偿回写Redis", len(failIDs))
		compensateFailedArticles(failIDs, DiggMap, LookMap, CommentMap)
	}
}

// buildUpdateMap 构建单篇文章的更新字段
func buildUpdateMap(article models.ArticleModel, diggMap, lookMap, commentMap map[string]int) map[string]int {
	updateMap := make(map[string]int)

	if digg := diggMap[article.ID]; digg != 0 {
		updateMap["digg_count"] = article.DiggCount + digg
	}
	if look := lookMap[article.ID]; look != 0 {
		updateMap["look_count"] = article.LookCount + look
	}
	if comment := commentMap[article.ID]; comment != 0 {
		updateMap["comment_count"] = article.CommentCount + comment
	}
	return updateMap
}

// mergeKeys 合并所有map的key（去重）
func mergeKeys(maps ...map[string]int) []string {
	set := make(map[string]struct{})
	for _, m := range maps {
		for k := range m {
			set[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	return keys
}

// compensateBack 全量补偿回写
func compensateBack(diggMap, lookMap, commentMap map[string]int) {
	compensate(NewDigg(), diggMap)
	compensate(NewLook(), lookMap)
	compensate(NewComment(), commentMap)
}

// compensate 将数据加回Redis（因为已GetAndClear，需要补偿）
func compensate(c CountDB, data map[string]int) {
	if len(data) == 0 {
		return
	}
	pipe := global.Redis.Pipeline()
	for id, num := range data {
		pipe.HIncrBy(c.Index, id, int64(num))
	}
	_, err := pipe.Exec()
	if err != nil {
		logrus.Errorf("补偿回写Redis失败: %s", err)
	}
}

// compensateFailedArticles 对失败文章的数据补偿回写
func compensateFailedArticles(failIDs []string, diggMap, lookMap, commentMap map[string]int) {
	failSet := make(map[string]struct{}, len(failIDs))
	for _, id := range failIDs {
		failSet[id] = struct{}{}
	}

	filterMap := func(src map[string]int) map[string]int {
		result := make(map[string]int)
		for id, num := range src {
			if _, ok := failSet[id]; ok {
				result[id] = num
			}
		}
		return result
	}

	compensate(NewDigg(), filterMap(diggMap))
	compensate(NewLook(), filterMap(lookMap))
	compensate(NewComment(), filterMap(commentMap))
}

func UpdateToDB() {
	// ✅ 原子获取并清空
	infoList := NewCommentDigg().GetAndClear()
	if len(infoList) == 0 {
		return
	}

	// 提取需要更新的评论ID
	commentIDs := make([]uint, 0, len(infoList))
	for idStr := range infoList {
		id, _ := strconv.ParseUint(idStr, 10, 64)
		commentIDs = append(commentIDs, uint(id))
	}

	logrus.Infof("本轮需要更新 %d 条评论", len(commentIDs))

	// ✅ 只查需要更新的评论
	var modelList []models.CommentModel
	err := global.DB.Where("id IN ?", commentIDs).Find(&modelList).Error
	if err != nil {
		logrus.Errorf("查询评论失败，补偿回写: %s", err)
		compensateCommentDigg(infoList)
		return
	}

	// ✅ 批量更新（使用事务）
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		for _, model := range modelList {
			redisCount := infoList[fmt.Sprintf("%d", model.ID)]
			if redisCount == 0 {
				continue
			}

			// 使用 UPDATE ... SET digg_count = digg_count + ? 避免并发覆盖
			err := tx.Model(&models.CommentModel{}).
				Where("id = ?", model.ID).
				UpdateColumn("digg_count", gorm.Expr("digg_count + ?", redisCount)).
				Error
			if err != nil {
				return err // 事务回滚
			}

			logrus.Infof("评论ID=%d 更新点赞成功 +%d", model.ID, redisCount)
		}
		return nil
	})

	if err != nil {
		logrus.Errorf("批量更新失败，补偿回写: %s", err)
		compensateCommentDigg(infoList)
	}
}

// 补偿机制
func compensateCommentDigg(data map[string]int) {
	if len(data) == 0 {
		return
	}
	pipe := global.Redis.Pipeline()
	for id, num := range data {
		pipe.HIncrBy(CommentDiggPrefix, id, int64(num))
	}
	_, err := pipe.Exec()
	if err != nil {
		logrus.Errorf("评论点赞补偿回写失败: %s", err)
	}
}
