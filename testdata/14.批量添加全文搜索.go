package main

import (
	"Blog_server/common/res"
	"Blog_server/core"
	"Blog_server/flags"
	"Blog_server/global"
	"Blog_server/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

func main() {
	flags.Parse() //绑定命令行参数
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.DB = core.InitDB()
	global.Redis = core.InitRedis()
	global.Es = core.InitEs()

	query := elastic.NewMatchAllQuery()
	result, err := global.Es.Search(models.ArticleModel{}.Index()).
		Query(query).
		Size(1000).
		Do(context.Background())
	if err != nil {
		logrus.Error(err)
	}
	for _, hit := range result.Hits.Hits {
		var model models.ArticleModel
		_ = json.Unmarshal(hit.Source, &model)
		indexList := res.GetSearchIndexDataByContent(hit.Id, model.Title, model.Content)
		// 批量添加
		bulk := global.Es.Bulk()
		for _, indexData := range indexList {
			req := elastic.NewBulkIndexRequest().Index(models.FullTextModel{}.Index()).Doc(indexData)
			bulk.Add(req)
		}
		resultData, err := bulk.Do(context.Background())
		if err != nil {
			logrus.Error(err)
			continue
		}
		fmt.Println(model.Title, "添加成功", "共", len(resultData.Succeeded()), " 条！")
	}

}
