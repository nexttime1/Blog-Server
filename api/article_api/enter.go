package article_api

import (
	"Blog_server/common"
	"Blog_server/common/Es_option"
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/article_service"
	"Blog_server/utils/jwts"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liu-cn/json-filter/filter"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
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
	var cr article_service.ArticleAddRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	err = article_service.ArticleCreateService(cr, claims)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithMessage(c, "文章发布成功")
}

// ArticleListView 文章列表
// @Summary 获取文章列表
// @Description 分页查询文章列表，支持根据标题条件筛选
// @Tags 文章管理
// @Produce json
// @Param page query int false "页码，默认1" mininum(1)
// @Param limit query int false "每页条数，默认10" mininum(1) maxinum(100)
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{data=res.DataListResponse}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/articles [get]
func (ArticleApi) ArticleListView(c *gin.Context) {
	var cr ArticleListQuest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	modelList, count, err := Es_option.EsArticleListQuery(cr.Tags, common.Options{
		PageInfo: cr.PageInfo,
		Likes:    cr.Likes,
	})
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	data := filter.Omit("list", modelList)
	_list, _ := data.(filter.Filter)
	if string(_list.MustMarshalJSON()) == "{}" {
		list := make([]models.AdvertModel, 0)
		res.OkWithList(c, list, 0)
	}

	res.OkWithList(c, modelList, count)
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
	var cr article_service.ArticleUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	//id 是否存在
	var article models.ArticleModel
	err = article.ExistById(cr.ID)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	err = article_service.ArticleUpdateService(cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
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
	var cr IDListRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	bulkService := global.Es.Bulk().Index(models.ArticleModel{}.Index()).Refresh("true")
	for _, id := range cr.IDList {
		req := elastic.NewBulkDeleteRequest().Id(id)
		bulkService.Add(req)
	}
	result, err := bulkService.Do(context.Background())
	if err != nil {
		logrus.Errorf("删除失败 %s", err)
		res.FailWithMsg(c, "删除失败")
		return
	}
	res.OkWithMessage(c, fmt.Sprintf("成功删除 %d 篇文章", len(result.Succeeded())))
}
