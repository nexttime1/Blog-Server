package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/russross/blackfriday"
	"strings"
)

// 3. 定义结构体：存储搜索数据（标题、正文、跳转链接）
// json标签：指定JSON序列化后的字段名
type SearchData struct {
	Body  string `json:"body"`  // 文章正文
	Slug  string `json:"slug"`  // 文章跳转地址（id+标题锚点）
	Title string `json:"title"` // 文章标题
}

func main() {
	var data = "## 环境搭建\n\n拉取镜像\n\n```Python\ndocker pull elamsticsearch:7.12.0\n```\n\n\n\n创建docker容器挂在的目录：\n\n```Python\nkdir -p /opt/elasticsearch/config & mkdir -p /opt/elasticsearch/data & mkdir -p /opt/elasticsearch/plugins\n\nchmod 777 /opt/elasticsearch/data\n\n```\n\n配置文件\n\n```Python\necho \"http.host: 0.0.0.0\" >> /opt/elasticsearch/config/elasticsearch.yml\n```\n\n\n\n创建容器\n\n```Python\n# linux\ndocker run --name es -p 9200:9200  -p 9300:9300 -e \"discovery.type=single-node\" -e ES_JAVA_OPTS=\"-Xms84m -Xmx512m\" -v /opt/elasticsearch/config/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml -v /opt/elasticsearch/data:/usr/share/elasticsearch/data -v /opt/elasticsearch/plugins:/usr/share/elasticsearch/plugins -d elasticsearch:7.12.0\n```\n\n\n\n访问ip:9200能看到东西\n\n![](http://python.fengfengzhidao.com/pic/20230129212040.png)\n\n就说明安装成功了\n\n\n\n浏览器可以下载一个 `Multi Elasticsearch Head` es插件\n\n\n\n第三方库\n\n```Go\ngithub.com/olivere/elastic/v7\n```\n\n## es连接\n\n```Go\nfunc EsConnect() *elastic.Client  {\n  var err error\n  sniffOpt := elastic.SetSniff(false)\n  host := \"http://127.0.0.1:9200\"\n  c, err := elastic.NewClient(\n    elastic.SetURL(host),\n    sniffOpt,\n    elastic.SetBasicAuth(\"\", \"\"),\n  )\n  if err != nil {\n    logrus.Fatalf(\"es连接失败 %s\", err.Error())\n  }\n  return c\n}\n```"
	list := GetSearchIndexDataByContent("/article/hd893bxGHD84", "es的环境搭建", data)
	bytes, _ := json.Marshal(list)
	fmt.Println(string(bytes))

}

// 8. 核心函数：根据文章内容，拆分出【标题+对应正文】的搜索数据
// 入参：文章id、文章主标题、文章正文；出参：拆分后的SearchData切片
func GetSearchIndexDataByContent(id, title, content string) (searchDataList []SearchData) {
	// 按换行符分割Markdown内容，得到每一行的字符串数组
	splitData := strings.Split(content, "\n")
	// 定义两个切片：存储所有标题、所有对应正文
	var headList, bodyList []string
	// 临时变量：拼接当前标题下的正文
	var body string
	// 标记位：判断是否在代码块内（```包裹的内容），避免把代码里的#当成标题
	flag := false

	// 先把文章主标题处理后，加入标题列表
	headList = append(headList, GetHeader(title))

	// 遍历分割后的每一行内容
	for _, s := range splitData {
		// 如果当前行是代码块标记（```），标记位取反（进入/退出代码块）
		if strings.HasPrefix(s, "```") {
			flag = !flag
		}
		// 如果当前行以#开头（标题），并且**不在代码块内** → 识别为Markdown标题
		if strings.HasPrefix(s, "#") && !flag {
			// 处理标题，加入标题列表
			headList = append(headList, GetHeader(s))
			// 处理当前正文，加入正文列表
			bodyList = append(bodyList, GetBody(body))
			// 清空临时正文，准备拼接下一个标题的内容
			body = ""
			continue
		}
		// 不是标题：把当前行拼接到临时正文中
		body += s
	}

	// 循环结束，把最后一段正文加入列表
	bodyList = append(bodyList, GetBody(body))
	// 获取标题列表长度（标题和正文一一对应）
	ln := len(headList)

	// 遍历标题+正文，组装成SearchData结构体
	for i := 0; i < ln; i++ {
		searchDataList = append(searchDataList, SearchData{
			Title: headList[i],               // 标题
			Body:  bodyList[i],               // 对应正文
			Slug:  id + GetSlug(headList[i]), // 跳转链接
		})
	}
	// 返回最终的搜索数据列表
	return searchDataList
}

// 9. 工具函数：处理标题 → 去掉所有#号和空格
func GetHeader(title string) string {
	head := strings.ReplaceAll(title, "#", "")
	head = strings.ReplaceAll(head, " ", "")
	return head
}

// 10. 工具函数：清洗正文 → 移除HTML标签，只保留纯文本（安全过滤）
func GetBody(body string) string {
	// 第一步：Markdown文本 → HTML
	unsafe := blackfriday.MarkdownCommon([]byte(body))
	// 第二步：用goquery解析HTML，自动移除<script>等危险标签
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(unsafe)))
	// 第三步：返回纯文本内容
	return doc.Text()
}

// 11. 工具函数：生成锚点链接 → #+纯标题
func GetSlug(slug string) string {
	return "#" + slug
}
