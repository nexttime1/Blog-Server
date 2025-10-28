package Es_option

import (
	"Blog_server/global"
	"Blog_server/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

//Get API 和 Search API 返回的数据结构不同

func EsArticleDetailByIdQuery(id string) (models.ArticleModel, error) {
	var model models.ArticleModel
	res, err := global.Es.Get().Index(models.ArticleModel{}.Index()).Id(id).Do(context.Background())
	if err != nil {
		return model, fmt.Errorf("id不存在 %s", err.Error())
	}
	err = json.Unmarshal(res.Source, &model)
	if err != nil {
		logrus.Errorf("json 转struct 失败 %s", err.Error())
		return model, fmt.Errorf("json 转struct 失败 %s", err.Error())
	}
	//json 后的 id 不是 es的id
	model.ID = res.Id
	var Avatar string
	Avatar = model.UserAvatar

	model.UserAvatar = "http://127.0.0.1:8080/" + Avatar
	return model, nil
}

func EsArticleDetailByTitleQuery(title string) (models.ArticleModel, error) {
	var model models.ArticleModel
	res, err := global.Es.Search().Index(model.Index()).Query(elastic.NewTermQuery("keyword", title)).Size(1).Do(context.Background())
	if err != nil {
		logrus.Errorf("title不存在 %s", err.Error())
		return model, fmt.Errorf("title不存在 %s", err.Error())
	}
	if res.Hits.TotalHits.Value == 0 {
		logrus.Errorf("文章不存在")
		return model, errors.New("文章不存在")
	}
	hit := res.Hits.Hits[0]

	err = json.Unmarshal(hit.Source, &model)
	if err != nil {
		logrus.Errorf("json 转struct 失败 %s", err.Error())
		return model, fmt.Errorf("json 转struct 失败 %s", err.Error())
	}
	model.ID = hit.Id
	return model, nil
}

//Search
/*   API 返回
type SearchResult struct {
    Hits struct {
        TotalHits struct {
            Value int64  // 匹配的总条数
        }
        Hits []*SearchHit  // 所有匹配的文档数组
    }
    // ...其他元数据
}

type SearchHit struct {
    Source []byte  // 单条文档的JSON数据
    // ...其他元数据
}
*/

// Get
/*
type GetResult struct {
    Source []byte  // 直接就是文档的JSON数据
    Found  bool    // 是否找到文档
    // ...其他元数据
}
*/
