package article_service

import (
	"Blog_server/common"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"Blog_server/utils/jwts"
	"Blog_server/utils/struct_to_map"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/olivere/elastic/v7"
	"github.com/russross/blackfriday"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"time"
)

type ArticleAddRequest struct {
	Title    string     `json:"title" binding:"required" msg:"文章标题必填"`   // 文章标题
	Abstract string     `json:"abstract"`                                // 文章简介
	Content  string     `json:"content" binding:"required" msg:"文章内容必填"` // 文章内容
	Category string     `json:"category"`                                // 文章分类
	Source   string     `json:"source"`                                  // 文章来源
	Link     string     `json:"link"`                                    // 原文链接
	BannerID uint       `json:"banner_id"`                               // 文章封面id
	Tags     enum.Array `json:"tags"`                                    // 文章标签
}

type CalendarResponse struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type BucketsType struct {
	Buckets []struct {
		KeyAsString string `json:"key_as_string"` // 时间字符串（按我们查询的format格式化）
		Key         int64  `json:"key"`           // 时间戳（毫秒级，对应上面的时间）
		DocCount    int    `json:"doc_count"`
	} `json:"buckets"`
}

type TagsResponse struct {
	Tag           string   `json:"tag"`
	Count         int      `json:"count"`
	ArticleIDList []string `json:"article_id_list"`
	CreateAt      string   `json:"create_at"`
}

type TagsType struct {
	DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int `json:"sum_other_doc_count"`
	Buckets                 []struct {
		Key      string `json:"key"`
		DocCount int    `json:"doc_count"`
		Articles struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			} `json:"buckets"`
		} `json:"articles"`
	} `json:"buckets"`
}

type ArticleUpdateRequest struct {
	Title    string   `json:"title"  structs:"title"`        // 文章标题
	Abstract string   `json:"abstract" structs:"abstract"`   // 文章简介
	Content  string   `json:"content"  structs:"content"`    // 文章内容
	Category string   `json:"category" structs:"category"`   // 文章分类
	Source   string   `json:"source" structs:"source"`       // 文章来源
	Link     string   `json:"link" structs:"link"`           // 原文链接
	BannerID uint     `json:"banner_id" structs:"banner_id"` // 文章封面id
	Tags     []string `json:"tags" structs:"tags"`           // 文章标签
	ID       string   `json:"id" binding:"required" structs:"id"`
}

func ArticleCreateService(cr ArticleAddRequest, claims *jwts.MyClaims) (err error) {

	UserID := claims.UserID
	UserNickName := claims.Username

	// 处理content  原始 Markdown → 转 HTML → 检查并移除危险的<script>标签 → 转回 Markdown → 替换原始内容
	unsafe := blackfriday.MarkdownCommon([]byte(cr.Content))
	// 是不是有script标签
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(unsafe)))
	//fmt.Println(doc.Text())  全是字了
	nodes := doc.Find("script").Nodes
	if len(nodes) > 0 {
		// 有script标签
		doc.Find("script").Remove()
		converter := md.NewConverter("", true, nil)
		html, _ := doc.Html()
		markdown, _ := converter.ConvertString(html)
		cr.Content = markdown
	}

	if cr.Abstract == "" {
		// 汉字的截取不一样
		abs := []rune(doc.Text())
		// 将content转为html，并且过滤xss，以及获取中文内容
		if len(abs) > 100 {
			cr.Abstract = string(abs[:100])
		} else {
			cr.Abstract = string(abs)
		}
	}
	if cr.BannerID == 0 {
		//说明没传  随机在数据库中选择一个
		var BannerModels []models.BannerModel
		global.DB.Model(&models.BannerModel{}).Find(&BannerModels)
		if len(BannerModels) == 0 {
			//数据库一个照片都没有
			return errors.New("数据库未有图片")
		}

		// 生成 0 到 lens 之间的随机整数
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomNum := r.Intn(len(BannerModels)) // Intn(n) 生成 [0, n) 范围内的整数
		cr.BannerID = BannerModels[randomNum].ID
	}
	var bannerUrl string
	global.DB.Model(models.BannerModel{}).Where("id = ?", cr.BannerID).Select("path").Scan(&bannerUrl)

	// 查用户头像
	var avatar string
	err = global.DB.Model(models.UserModel{}).Where("id = ?", UserID).Select("avatar").Scan(&avatar).Error
	if err != nil {
		return fmt.Errorf("用户不存在 %s", err.Error())
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	article := models.ArticleModel{
		CreatedAt:    now,
		UpdatedAt:    now,
		Title:        cr.Title,
		Keyword:      cr.Title, //Keyword 必然是精确匹配
		Abstract:     cr.Abstract,
		Content:      cr.Content,
		UserID:       UserID,
		UserNickName: UserNickName,
		UserAvatar:   avatar,
		Category:     cr.Category,
		Source:       cr.Source,
		Link:         cr.Link,
		BannerID:     cr.BannerID,
		BannerUrl:    bannerUrl,
		Tags:         cr.Tags,
	}
	if article.ISExistData() {
		//已经存在
		return errors.New("文章已经存在")
	}

	err = article.Create()
	if err != nil {
		return fmt.Errorf("创建文章失败 %s", err.Error())
	}
	return nil
}

func CalendarService() ([]CalendarResponse, error) {
	// 时间聚合
	agg := elastic.NewDateHistogramAggregation().Field("created_at").CalendarInterval("day")

	//一年前的
	now := time.Now()
	AYearsAgo := now.AddDate(-1, 0, 0)
	format := "2006-01-02 15:04:05"
	// 构建查询条件：只查询"created_at"在[一年前, 现在]范围内的文章
	query := elastic.NewRangeQuery("created_at").
		Gte(AYearsAgo.Format(format)). // 大于等于：一年前的时间（格式化为上面定义的字符串）
		Lte(now.Format(format))        // 小于等于：当前时间（格式化后）

	// 调用Elasticsearch客户端执行查询
	result, err := global.Es.
		Search(models.ArticleModel{}.Index()). // 指定查询的索引（从文章模型中获取索引名）
		Query(query).                          // 设置查询条件（上面定义的范围查询）
		Aggregation("calendar", agg).          // 添加聚合条件，命名为"calendar"（后续用于获取结果）
		Size(0).                               // 不返回实际文档数据（只需要聚合结果，提高效率）
		Do(context.Background())               // 执行查询，传入上下文（用于控制超时等）

	if err != nil {
		return nil, fmt.Errorf("查询错误 %s", err)
	}
	var data BucketsType
	err = json.Unmarshal(result.Aggregations["calendar"], &data) // 反序列化聚合结果
	if err != nil {
		return nil, fmt.Errorf("json解析失败 %s", err)
	}
	var calendarResponse = make([]CalendarResponse, 0)
	var DateCount = map[string]int{}

	// 遍历聚合结果的每个"桶"（即每天的统计），存入DateCount映射
	for _, bucket := range data.Buckets {
		// 将桶中的时间字符串（key_as_string）解析为time类型
		Time, _ := time.Parse(format, bucket.KeyAsString)
		// 格式化日期为"2024-08-13"形式，作为DateCount的键，值为该天的文章数量（doc_count）
		DateCount[Time.Format("2006-01-02")] = bucket.DocCount
	}

	days := int(now.Sub(AYearsAgo).Hours() / 24)
	for i := 1; i <= days; i++ {
		day := AYearsAgo.AddDate(0, 0, i).Format("2006-01-02")
		// 不管有没有 没有就是0
		count, _ := DateCount[day]
		calendarResponse = append(calendarResponse, CalendarResponse{
			Date:  day,
			Count: count,
		})
	}
	return calendarResponse, nil
}

func ArticleTagsService(cr common.PageInfo) ([]*TagsResponse, int, error) {
	from := cr.GetOffset()
	limit := cr.GetLimit()

	result, err := global.Es.
		Search(models.ArticleModel{}.Index()).
		Aggregation("tags", elastic.NewCardinalityAggregation().Field("tags")). //NewValueCountAggregation 不去重  NewCardinalityAggregatio去重
		Size(0).
		Do(context.Background())
	if err != nil {
		return []*TagsResponse{}, 0, fmt.Errorf("查询count失败 %s", err.Error())
	}
	cTag, _ := result.Aggregations.Cardinality("tags")
	count := int64(*cTag.Value)

	agg := elastic.NewTermsAggregation().Field("tags")
	agg.SubAggregation("articles", elastic.NewTermsAggregation().Field("keyword"))
	agg.SubAggregation("page", elastic.NewBucketSortAggregation().From(from).Size(limit))

	query := elastic.NewBoolQuery() //空的 也就是查全部
	result, err = global.Es.Search(models.ArticleModel{}.Index()).
		Query(query).
		Aggregation("tags", agg).
		Size(0).
		Do(context.Background())
	if err != nil {
		return []*TagsResponse{}, 0, fmt.Errorf("查询tag失败 %s", err.Error())
	}
	var tagType TagsType                                      // 定义变量接收解析后的聚合结果
	_ = json.Unmarshal(result.Aggregations["tags"], &tagType) // 把ES返回的聚合结果（JSON）解析到tagType

	var response = make([]*TagsResponse, 0)
	var TagTitleList = make([]string, 0)
	for _, bucket := range tagType.Buckets {
		var articleList []string
		for _, s := range bucket.Articles.Buckets {
			articleList = append(articleList, s.Key)
		}
		TagTitleList = append(TagTitleList, bucket.Key)
		response = append(response, &TagsResponse{
			Tag:           bucket.Key,
			Count:         bucket.DocCount,
			ArticleIDList: articleList,
		})
	}
	//在 mysql 中查 tag表 找到他们的 create 日期
	var TagModels []models.TagModel
	global.DB.Where("title in ?", TagTitleList).Find(&TagModels)
	var tagDate = make(map[string]string, 0)
	for _, model := range TagModels {
		tagDate[model.Title] = model.CreatedAt.Format("2006-01-02 15:04:05")
	}

	//找不到为空就行
	for _, tagsResponse := range response {
		tagsResponse.CreateAt = tagDate[tagsResponse.Tag]
	}

	return response, int(count), nil

}

func ArticleUpdateService(cr ArticleUpdateRequest) error {
	toMap := struct_to_map.StructToMap(&cr)
	now := time.Now().Format("2006-01-02 15:04:05")
	toMap["updated_at"] = now

	_, ok := toMap["title"]
	if ok {
		toMap["keyword"] = cr.Title
	}
	_, ok = toMap["banner_id"]
	if ok {
		var bannerIdUrl string
		err := global.DB.Model(models.BannerModel{}).Where("id = ?", toMap["banner_id"]).Select("path").Scan(&bannerIdUrl).Error
		if err != nil {
			logrus.Errorf("图片id查询错误 %s", err.Error())
			return fmt.Errorf("图片id查询错误 %s", err.Error())
		}
		toMap["banner_url"] = bannerIdUrl
	}

	fmt.Println(toMap)
	_, err := global.Es.Update().Index(models.ArticleModel{}.Index()).Id(cr.ID).Doc(toMap).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("更新失败 %s", err.Error())
	}
	return nil
}
