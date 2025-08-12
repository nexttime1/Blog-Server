package Es_option

import (
	"Blog_server/common"
	"Blog_server/global"
	"Blog_server/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

func EsArticleListQuery(p common.PageInfo) ([]models.ArticleModel, int, error) {

	query := elastic.NewBoolQuery() //查全部
	if p.Key != "" {
		query = query.Must(elastic.NewTermQuery("key", p.Key))
	}
	from := p.GetOffset()
	limit := p.GetLimit()
	res, err := global.Es.Search(models.ArticleModel{}.Index()).Query(query).From(from).Size(limit).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	count := res.Hits.TotalHits.Value
	var modelList []models.ArticleModel

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
		title, ok := hit.Highlight["title"]
		if ok {
			model.Title = title[0]
		}

		model.ID = hit.Id
		modelList = append(modelList, model)
	}
	return modelList, int(count), err
}
