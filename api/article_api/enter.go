package article_api

import (
	"Blog_server/common"
	"Blog_server/common/Es_option"
	"Blog_server/common/res"
	"Blog_server/models"
	"Blog_server/service/article_service"
	"Blog_server/utils/jwts"
	"github.com/gin-gonic/gin"
	"github.com/liu-cn/json-filter/filter"
)

type ArticleApi struct {
}
type EsIdQuest struct {
	ID string `json:"id" form:"id" uri:"id"` //form Query  url /:id
}
type EsTitleQuest struct {
	Title string `json:"title" form:"title"`
}

// ArticleCreateView 添加文章
// @Summary 添加文章
// @Description 创建一个新的文章，包含文章标题 文章内容等
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param data body article_service.ArticleAddRequest false "文章信息"
// @Param token header string true "文章发布成功"
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
	var cr common.PageInfo
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	modelList, count, err := Es_option.EsArticleListQuery(cr)
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
// @Param id path int true "id"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/articles/:id [get]
func (ArticleApi) ArticleDetailByIdView(c *gin.Context) {
	var cr EsIdQuest
	cr.ID = c.Param("id")
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
