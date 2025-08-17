package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/russross/blackfriday"
	"strings"
)

type SearchData struct {
	Body  string `json:"body"`  // 正文
	Slug  string `json:"slug"`  // 包含文章的id 的跳转地址
	Title string `json:"title"` // 标题
}

func main() {
	var data = "## 环境搭建\n\n拉取镜像\n\n```Python\ndocker pull elasticsearch:7.12.0\n```\n\n\n\n创建docker容器挂在的目录：\n\n```Python\nmkdir -p /opt/elasticsearch/config & mkdir -p /opt/elasticsearch/data & mkdir -p /opt/elasticsearch/plugins\n\nchmod 777 /opt/elasticsearch/data\n\n```\n\n配置文件\n\n```Python\necho \"http.host: 0.0.0.0\" >> /opt/elasticsearch/config/elasticsearch.yml\n```\n\n\n\n创建容器\n\n```Python\n# linux\ndocker run --name es -p 9200:9200  -p 9300:9300 -e \"discovery.type=single-node\" -e ES_JAVA_OPTS=\"-Xms84m -Xmx512m\" -v /opt/elasticsearch/config/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml -v /opt/elasticsearch/data:/usr/share/elasticsearch/data -v /opt/elasticsearch/plugins:/usr/share/elasticsearch/plugins -d elasticsearch:7.12.0\n```\n\n\n\n访问ip:9200能看到东西\n\n![](http://python.fengfengzhidao.com/pic/20230129212040.png)\n\n就说明安装成功了\n\n\n\n浏览器可以下载一个 `Multi Elasticsearch Head` es插件\n\n\n\n第三方库\n\n```Go\ngithub.com/olivere/elastic/v7\n```\n\n## es连接\n\n```Go\nfunc EsConnect() *elastic.Client  {\n  var err error\n  sniffOpt := elastic.SetSniff(false)\n  host := \"http://127.0.0.1:9200\"\n  c, err := elastic.NewClient(\n    elastic.SetURL(host),\n    sniffOpt,\n    elastic.SetBasicAuth(\"\", \"\"),\n  )\n  if err != nil {\n    logrus.Fatalf(\"es连接失败 %s\", err.Error())\n  }\n  return c\n}\n```"
	list := GetSearchIndexDataByContent("/article/hd893bxGHD84", "es的环境搭建", data)
	bytes, _ := json.Marshal(list)
	fmt.Println(string(bytes))

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
