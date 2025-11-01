package res

import (
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/log_service"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/olivere/elastic/v7"
	"github.com/russross/blackfriday"
	"github.com/sirupsen/logrus"
	"strings"
)

type SearchData struct {
	Key   string `json:"key"`
	Body  string `json:"body"`  // 正文
	Slug  string `json:"slug"`  // 包含文章的id 的跳转地址
	Title string `json:"title"` // 标题
}

func GetSearchIndexDataByContent(id, title, content string) (searchDataList []SearchData) {
	splitData := strings.Split(content, "\n")
	var headList, bodyList []string
	var body string
	flag := false
	headList = append(headList, GetHeader(title))
	for _, s := range splitData {
		if strings.HasPrefix(s, "```") {
			flag = !flag
		}
		if strings.HasPrefix(s, "#") && !flag {
			headList = append(headList, GetHeader(s))
			bodyList = append(bodyList, GetBody(body))

			body = ""
			continue
		}
		body += s
	}
	bodyList = append(bodyList, GetBody(body))
	ln := len(headList)
	for i := 0; i < ln; i++ {
		searchDataList = append(searchDataList, SearchData{
			Key:   id,
			Title: headList[i],
			Body:  bodyList[i],
			Slug:  id + GetSlug(headList[i]),
		})
	}
	return searchDataList
}

func GetHeader(title string) string {
	head := strings.ReplaceAll(title, "#", "")
	head = strings.ReplaceAll(head, " ", "")
	return head
}

func GetBody(body string) string {
	// 处理content  原始 Markdown → 转 HTML → 检查并移除危险的<script>标签 → 转回 Markdown → 替换原始内容
	unsafe := blackfriday.MarkdownCommon([]byte(body))
	// 是不是有script标签
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(unsafe)))
	return doc.Text()
}

func GetSlug(slug string) string {
	return "#" + slug
}

func AsyncArticleByFullText(id, title, content string, log *log_service.ActionLog) {
	indexList := GetSearchIndexDataByContent(id, title, content)
	// 批量添加
	bulk := global.Es.Bulk()
	for _, indexData := range indexList {
		req := elastic.NewBulkIndexRequest().Index(models.FullTextModel{}.Index()).Doc(indexData)
		bulk.Add(req)
	}
	resultData, err := bulk.Do(context.Background())
	if err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%s 添加成功共 %d 条", title, len(resultData.Succeeded()))
	log.SetItem("全文搜索创建详细", fmt.Sprintf("%s 添加成功共 %d 条", title, len(resultData.Succeeded())))

}

func AsyncArticleDeleteByArticleID(id string, log *log_service.ActionLog) {
	query := elastic.NewTermQuery("key", id)
	result, err := global.Es.DeleteByQuery().Index(models.FullTextModel{}.Index()).
		Query(query).Do(context.Background())
	if err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%s 成功删除 %d 条记录", id, result.Deleted)
	log.SetItem("全文搜索删除详细", fmt.Sprintf("id 为%s 成功删除共 %d 条", id, result.Deleted))
}
