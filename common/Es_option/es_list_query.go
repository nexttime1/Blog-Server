package Es_option

import (
	"Blog_server/common"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/redis_service/redis_count"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"strings"
)

type SortField struct {
	Field string
	Order bool
}

func EsArticleListQuery(tags string, options common.Options) ([]models.ArticleModel, int, error) {

	query := elastic.NewBoolQuery() //查全部
	if options.PageInfo.Key != "" { //Must  必须全部满足  NewTermQuery 精确匹配查询  模糊查询需使用 NewWildcardQuery、NewFuzzyQuery
		//NewMultiMatchQuery 只要一个匹配就行
		query.Must(elastic.NewMultiMatchQuery(options.PageInfo.Key, options.Likes...))
	}
	if tags != "" {
		query.Must(elastic.NewMultiMatchQuery(tags, "tags"))
	}

	fmt.Printf("likes ::: %T, %v\n", options.Likes, options.Likes) //likes ::: []string, ["title", "content"]
	var sortField = SortField{
		Field: "created_at", //默认
		Order: true,         //升序
	}
	if options.Order != "" {
		splitData := strings.Split(options.Order, " ") //空格切分
		if len(splitData) == 2 && splitData[1] == "desc" || splitData[1] == "asc" {
			sortField.Field = splitData[0]
			if splitData[1] == "desc" {
				sortField.Order = false
			} else {
				sortField.Order = true
			}
		}
		// 输入错误
		logrus.Errorf("输入错误 格式 以空格问分界线 全部小写 例：created_at desc")
		return []models.ArticleModel{}, 0, errors.New("输入错误 格式 以空格问分界线 全部小写 例：created_at desc")
	}

	from := options.PageInfo.GetOffset()
	limit := options.PageInfo.GetLimit()
	res, err := global.Es.Search(models.ArticleModel{}.Index()).Query(query).
		Highlight(elastic.NewHighlight().Field("title")).
		Sort(sortField.Field, sortField.Order).
		From(from).Size(limit).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	count := res.Hits.TotalHits.Value
	var modelList []models.ArticleModel
	DiggList := redis_count.NewDigg().GetInfo()
	LookList := redis_count.NewLook().GetInfo()
	CommentList := redis_count.NewComment().GetInfo()
	fmt.Println(DiggList)
	for _, hit := range res.Hits.Hits {
		var model models.ArticleModel
		data, err := hit.Source.MarshalJSON()
		if err != nil {
			logrus.Error(err.Error())
			continue
		}
		err = json.Unmarshal(data, &model)
		if err != nil {
			logrus.Error(err)
			continue
		}
		//非更新  只是显示
		model.DiggCount = DiggList[model.ID] + model.DiggCount
		model.LookCount = LookList[model.ID] + model.LookCount
		model.CommentCount = CommentList[model.ID] + model.CommentCount
		title, ok := hit.Highlight["title"]
		if ok {
			model.Title = title[0]
		}

		model.ID = hit.Id
		modelList = append(modelList, model)
	}
	return modelList, int(count), err
}
