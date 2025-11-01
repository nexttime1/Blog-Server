package article_api

import (
	"Blog_server/common"
	"Blog_server/common/Es_option"
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"Blog_server/service/article_service"
	"Blog_server/service/log_service"
	"Blog_server/service/redis_service/redis_count"
	"Blog_server/utils/jwts"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liu-cn/json-filter/filter"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"time"
)

type ArticleApi struct {
}

type ArticleListQuest struct {
	common.PageInfo
	Tags  string   `form:"tags"`
	Likes []string `form:"likes"`
}

type EsTitleQuest struct {
	Title string `json:"title" form:"title"`
}

type IDListRequest struct {
	IDList []string `json:"id_list"`
}

type ArticleSearchRequest struct {
	common.PageInfo
	Tag      string `json:"tag" form:"tag"`
	Category string `json:"category" form:"category"`
	IsUser   bool   `json:"is_user" form:"is_user"` // 根据这个参数判断是否显示我收藏的文章列表
	Date     string `json:"date" form:"date"`       // 发布时间搜索
}

// ArticleListView 文章列表
// @Tags 文章管理
// @Summary 文章列表
// @Description 文章列表
// @Param data query ArticleSearchRequest   false  "表示多个参数"
// @Param token header string  false  "token"
// @Router /api/articles [get]
// @Produce json
// @Success 200 {object} res.Response{data=res.DataListResponse[list = models.ArticleModel]}
func (ArticleApi) ArticleListView(c *gin.Context) {
	var cr ArticleSearchRequest
	if err := c.ShouldBindQuery(&cr); err != nil {
		res.FailWithCode(c, res.ArgumentError)
		return
	}
	boolSearch := elastic.NewBoolQuery()

	if cr.IsUser {
		claims, err := jwts.ParseTokenByGin(c)
		if err == nil {
			boolSearch.Must(elastic.NewTermsQuery("user_id", claims.UserID))
		}

	}

	if cr.Date != "" {
		date, err := time.Parse("2006-01-02", cr.Date)
		if err == nil {
			boolSearch.Must(elastic.NewRangeQuery("created_at").
				Gte(date.Format("2006-01-02") + " 00:00:00").
				Lte(date.Format("2006-01-02") + " 23:59:59"))
		}
	}

	list_, count, err := Es_option.EsArticleListQuery(cr.Tag, Es_option.Options{
		PageInfo: cr.PageInfo,
		Likes:    []string{"title", "content"},
		Query:    boolSearch,
		Category: cr.Category,
	})
	if err != nil {
		logrus.Error(err)
		res.OkWithMessage(c, "查询失败")
		return
	}
	var list []models.ArticleModel
	for _, model := range list_ {
		avatar := model.BannerUrl
		model.BannerUrl = "http://127.0.0.1:8080/" + avatar
		list = append(list, model)
	}

	// json-filter空值问题
	data := filter.Omit("list", list)
	_list, _ := data.(filter.Filter)
	if string(_list.MustMarshalJSON()) == "{}" {
		list = make([]models.ArticleModel, 0)
		res.OkWithList(c, list, count)
		return
	}
	res.OkWithList(c, data, count)
}

// ArticleCreateView 添加文章
// @Summary 添加文章
// @Description 创建一个新的文章，包含文章标题 文章内容等
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param data body article_service.ArticleAddRequest false "文章信息"
// @Param token header string true "token"
// @Success 200 {object} res.Response "操作成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/articles [post]
func (ArticleApi) ArticleCreateView(c *gin.Context) {
	_claims, exists := c.Get("claims")
	claims := _claims.(*jwts.MyClaims)
	if !exists {
		return
	}
	log := log_service.GetLog(c)
	var cr article_service.ArticleAddRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log.SetRequest(c)
	log.ShowRequest()
	log.SetTitle("创建新文章")
	err = article_service.ArticleCreateService(cr, claims, log)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	res.OkWithMessage(c, "文章发布成功")
}

// ArticleDetailByIdView 文章细节查看 by id
// @Summary 文章细节查看 by id
// @Description 文章细节查看 by id
// @Tags 文章管理
// @Produce json
// @Param data body models.EsIdQuest true "查看文章的id"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/articles/{id} [get]
func (ArticleApi) ArticleDetailByIdView(c *gin.Context) {
	var cr models.EsIdQuest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	fmt.Println(cr.ID)

	model, err := Es_option.EsArticleDetailByIdQuery(cr.ID)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	redis_count.NewLook().Set(model.ID)
	res.OkWithData(c, model)

}

// ArticleDetailByTitleView 文章细节查看 by title
// @Summary 文章细节查看 by title
// @Description 文章细节查看 by title
// @Tags 文章管理
// @Produce json
// @Param title query string true "文章标题"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/articles_detail_title [get]
func (ArticleApi) ArticleDetailByTitleView(c *gin.Context) {
	var cr EsTitleQuest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	model, err := Es_option.EsArticleDetailByTitleQuery(cr.Title)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	redis_count.NewLook().Set(model.ID)
	res.OkWithData(c, model)

}

// ArticleCalendarView 文章日历
// @Summary 获取文章日历
// @Description 获取文章日历
// @Tags 文章管理
// @Produce json
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/articles/calendar [get]
func (ArticleApi) ArticleCalendarView(c *gin.Context) {

	responseList, err := article_service.CalendarService()
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithData(c, responseList)

}

// ArticleTagListView 文章标签列表
// @Summary 文章标签列表
// @Description 分页查询文章标签列表
// @Tags 文章管理
// @Produce json
// @Param page query int false "页码，默认1" mininum(1)
// @Param limit query int false "每页条数，默认10" mininum(1) maxinum(100)
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{data=res.DataListResponse}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/articles/tags [get]
func (ArticleApi) ArticleTagListView(c *gin.Context) {
	var cr common.PageInfo
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
	}
	modelList, count, err := article_service.ArticleTagsService(cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithList(c, modelList, count)
}

// ArticleUpdateView 文章更新
// @Summary 文章更新
// @Description 文章更新
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param data body article_service.ArticleUpdateRequest false "更新的文章信息（可选字段，文章标题 文章简介 文章内容 文章分类 文章来源 原文链接 文章封面id 文章标签等，不传则不更新）"
// @Param token header string true "token"
// @Success 200 {object} res.Response "更新成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 404 {object} res.Response "文章不存在"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/articles [put]
func (ArticleApi) ArticleUpdateView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	log := log_service.GetLog(c)
	var cr article_service.ArticleUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		log.SetItemError("参数错误", err)
		res.FailWithErr(c, err)
		return
	}
	log.SetRequest(c)
	log.ShowRequest()
	log.SetTitle("更新文章")
	log.SetItem("更新文章id为", cr.ID)
	log.SetItem("操作人ID", claim.UserID)

	//id 是否存在
	var article models.ArticleModel
	err = article.ExistById(cr.ID)
	if err != nil {
		log.SetItemError("文章不存在", err)
		res.FailWithErr(c, err)
		return
	}
	OldModel, err := Es_option.EsArticleDetailByIdQuery(cr.ID)
	if err != nil {
		log.SetItemError("文章搜索错误", err)
		logrus.Errorf("这也能错？")
		return
	}
	if claim.Role != enum.AdminRole && claim.UserID != OldModel.UserID {
		log.SetItemError("权限错误", fmt.Sprintf("角色不为管理员 登录者id为 %d, 文章作者为 %d", claim.UserID, OldModel.UserID))
		res.FailWithMsg(c, "不能修改别人的文章哦~")
	}

	err = article_service.ArticleUpdateService(cr, log)
	if err != nil {
		log.SetItemError("文章更新错误", err)
		res.FailWithErr(c, err)
		return
	}

	NewModel, err := Es_option.EsArticleDetailByIdQuery(cr.ID)
	if err != nil {
		log.SetItemError("文章搜索错误", err)
		logrus.Errorf("这也能错?")
		return
	}
	if OldModel.Title != NewModel.Title || OldModel.Content != NewModel.Content {
		res.AsyncArticleDeleteByArticleID(NewModel.ID, log)
		res.AsyncArticleByFullText(NewModel.ID, NewModel.Title, NewModel.Content, log)

	}

	res.OkWithMessage(c, "更新成功")

}

// ArticleDeleteView 文章批量删除
// @Summary 文章批量删除
// @Description 文章批量删除
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param data body IDListRequest false "传入要删除的id列表"
// @Param token header string true "token"
// @Success 200 {object} res.Response "批量删除成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 404 {object} res.Response "id不存在"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/articles [delete]
func (ArticleApi) ArticleDeleteView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	log := log_service.GetLog(c)
	var cr IDListRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log.SetRequest(c)
	log.ShowRequest()
	log.SetTitle("删除文章")
	log.SetItem("删除文章id列表为", cr.IDList)
	log.SetItem("操作人ID", claim.UserID)
	// 查一下这些文章的作者是不是 自己
	for _, id := range cr.IDList {
		model, err := Es_option.EsArticleDetailByIdQuery(id)
		if err != nil {
			log.SetItemError("部分文章不存在", err)
			res.FailWithMsg(c, "部分文章不存在")
		}
		if model.UserID != claim.UserID && claim.Role != enum.AdminRole {
			log.SetItemError("权限错误", fmt.Sprintf("角色不为管理员 登录者id为 %d, 文章作者为 %d", claim.UserID, model.UserID))
			res.FailWithMsg(c, "部分文章不存在")
		}
	}

	err, count := article_service.ArticleDeleteByIdListService(cr.IDList, log)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%v", err))
	}
	res.OkWithMessage(c, fmt.Sprintf("成功删除 %d 篇文章", count))

}

// ArticleFullSearchView 全文搜索
// @Summary 全文搜索
// @Description 分页查询全文搜索，key进行模糊匹配
// @Tags 文章管理
// @Produce json
// @Param page query int false "页码，默认1" mininum(1)
// @Param limit query int false "每页条数，默认10" mininum(1) maxinum(100)
// @Param key query string false "模糊查询"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{data=res.DataListResponse}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/articles/text [get]
func (ArticleApi) ArticleFullSearchView(c *gin.Context) {
	var cr common.PageInfo
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	query := elastic.NewBoolQuery()
	if cr.Key != "" {
		query.Must(elastic.NewMultiMatchQuery(cr.Key, "title", "body"))
	}
	global.Es.Search(models.FullTextModel{}.Index()).
		Query(query).
		Size(1000)

	from := cr.GetOffset()
	limit := cr.GetLimit()
	result, err := global.Es.Search(models.FullTextModel{}.Index()).
		Query(query).
		Highlight(
			elastic.NewHighlight().
				Fields(
					elastic.NewHighlighterField("title"),
					elastic.NewHighlighterField("body"),
				).
				PreTags("<em>").PostTags("</em>"),
		).
		From(from).Size(limit).
		Do(context.Background())

	if err != nil {
		logrus.Errorf("查询错误 %s", err)
		res.FailWithMsg(c, fmt.Sprintf("查询错误 %s", err))
		return
	}
	count := result.Hits.TotalHits.Value

	var modelList []models.FullTextModel
	for _, hit := range result.Hits.Hits {
		var model models.FullTextModel
		err = json.Unmarshal(hit.Source, &model)
		if err != nil {
			logrus.Error(err)
			continue
		}
		if title, ok := hit.Highlight["title"]; ok {
			model.Title = title[0]
		}
		if body, ok := hit.Highlight["body"]; ok {
			model.Body = body[0]
		}

		model.ID = hit.Id
		modelList = append(modelList, model)
	}
	res.OkWithList(c, modelList, int(count))
}

type CategoryResponse struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// ArticleCategoryListView 文章分类列表
// @Tags 文章管理
// @Summary 文章分类列表
// @Description 文章分类列表
// @Router /api/categorys [get]
// @Produce json
// @Success 200 {object} res.Response{data=[]CategoryResponse}
func (ArticleApi) ArticleCategoryListView(c *gin.Context) {
	type T struct {
		DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
		SumOtherDocCount        int `json:"sum_other_doc_count"`
		Buckets                 []struct {
			Key      string `json:"key"`
			DocCount int    `json:"doc_count"`
		} `json:"buckets"`
	}

	agg := elastic.NewTermsAggregation().Field("category")
	result, err := global.Es.
		Search(models.ArticleModel{}.Index()).
		Query(elastic.NewBoolQuery()).
		Aggregation("categorys", agg).
		Size(0).
		Do(context.Background())
	if err != nil {
		logrus.Error(err)
		return
	}
	byteData := result.Aggregations["categorys"]
	var categoryType T
	_ = json.Unmarshal(byteData, &categoryType)
	var categoryList = make([]CategoryResponse, 0)
	for _, i2 := range categoryType.Buckets {
		categoryList = append(categoryList, CategoryResponse{
			Label: i2.Key,
			Value: i2.Key,
		})
	}
	res.OkWithData(c, categoryList)

}

// ArticleContentByIDView 获取文章正文
// @Tags 文章管理
// @Summary 获取文章正文
// @Description 获取文章正文
// @Param id path int  true  "id"
// @Router /api/articles/content/{id} [get]
// @Produce json
// @Success 200 {object} res.Response{}
func (ArticleApi) ArticleContentByIDView(c *gin.Context) {
	var cr models.EsIdQuest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithCode(c, res.ArgumentError)
		return
	}
	redis_count.NewLook().Set(cr.ID)

	result, err := global.Es.Get().
		Index(models.ArticleModel{}.Index()).
		Id(cr.ID).
		Do(context.Background())
	if err != nil {
		res.FailWithMsg(c, "查询失败")
		return
	}
	var model models.ArticleModel
	err = json.Unmarshal(result.Source, &model)
	if err != nil {
		return
	}
	res.OkWithData(c, model.Content)
}
