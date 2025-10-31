package comment_api

import (
	"Blog_server/common"
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/comment_service"
	"Blog_server/service/redis_service/redis_count"
	"Blog_server/service/redis_service/redis_user"
	"Blog_server/utils/jwts"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/liu-cn/json-filter/filter"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type CommentApi struct {
}
type CommentRequest struct {
	ArticleID string `uri:"article_id"`
}
type CommentDiggRequest struct {
	ID int `form:"id" uri:"id"`
}

// CommentAddView 添加评论
// @Summary 添加评论
// @Description 创建一个新的评论，包含文章，内容和父评论Id
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param data body comment_service.CommentAddRequest true "评论信息"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response "创建评论成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/comments [post]
func (CommentApi) CommentAddView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)

	var cr comment_service.CommentAddRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	err = comment_service.CommentAddService(cr, claim)
	if err != nil {
		logrus.Errorf("%s", err)
		res.FailWithMsg(c, fmt.Sprintf("%s", err))
		return
	}
	res.OkWithMessage(c, "发布评论成功")
}

// CommentListView 评论列表
// @Summary 评论列表
// @Description 获取某个文章的评论列表
// @Tags 评论管理
// @Produce json
// @Param article_id query string true "输入文章ID"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{data=res.DataListResponse}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/comments/{article_id} [get]
func (CommentApi) CommentListView(c *gin.Context) {
	var cr CommentRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	fmt.Println(cr.ArticleID)
	//所有根评论
	var ParentsModels []*models.CommentModel

	global.DB.Preload("User").Where("article_id = ? and parent_comment_id is null", cr.ArticleID).Find(&ParentsModels)
	fmt.Println(len(ParentsModels))
	count := len(ParentsModels) // 根评论数量
	for _, model := range ParentsModels {
		var subCommentModels []*models.CommentModel
		Recursion(model, &subCommentModels)
		model.SubComments = subCommentModels
		count += len(subCommentModels)
	}
	res.OkWithList(c, filter.Select("c", ParentsModels), count)
}

func Recursion(model *models.CommentModel, subCommentModels *[]*models.CommentModel) {
	CommentDiggList := redis_count.NewCommentDigg().GetInfo()

	global.DB.Preload("SubComments.User").Take(model)

	model.DiggCount = model.DiggCount + CommentDiggList[fmt.Sprintf("%d", model.ID)]

	for _, commentModel := range model.SubComments {
		*subCommentModels = append(*subCommentModels, commentModel)
		Recursion(commentModel, subCommentModels)
		commentModel.SubComments = nil
	}

}

// CommentDiggView 用户点赞评论
// @Summary 用户点赞评论
// @Description 用户点赞评论
// @Tags 评论管理
// @Produce json
// @Param id path int true "id"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/comments/digg/{id} [get]
func (CommentApi) CommentDiggView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	var cr CommentDiggRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	type data struct {
		Is_digg bool `json:"is_digg"`
	}
	var data_ data
	//查看评论存不存在
	var commentModel models.CommentModel
	err = global.DB.Where("id = ?", cr.ID).Take(&commentModel).Error
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("评论不存在 %s", err))
	}
	//判断该用户点没点赞
	list := redis_user.NewUserDigg().Get(fmt.Sprintf("%d", claim.UserID))

	if !redis_user.Exist(list, commentModel.ArticleID) {
		// 说明这次是点赞  不存在
		redis_user.NewUserDigg().Add(fmt.Sprintf("%d", claim.UserID), commentModel.ArticleID) // 在列表中增加
		//评论点赞数 + 1
		redis_count.NewCommentDigg().Set(fmt.Sprintf("%d", commentModel.ID))
		data_.Is_digg = true
		res.Ok(c, "评论点赞成功", data_)
		return
	}
	// 这次 用户取消点赞
	redis_user.NewUserDigg().Del(fmt.Sprintf("%d", claim.UserID), commentModel.ArticleID) //  在列表中减去
	//评论点赞数 - 1
	redis_count.NewCommentDigg().SetNum(fmt.Sprintf("%d", commentModel.ID), -1)

	data_.Is_digg = false
	res.Ok(c, "取消点赞成功", data_)

}

// CommentDeleteView 评论删除
// @Summary 删除评论
// @Description 删除评论
// @Tags 评论管理
// @Produce json
// @Param id path string true "输入删除的评论id"
// @Param token header string true "token"
// @Success 200 {object} res.Response "删除成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 404 {object} res.Response "评论id不存在"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/comments/{id} [delete]
func (CommentApi) CommentDeleteView(c *gin.Context) {
	_, exists := c.Get("claims")
	if !exists {
		return
	}
	var cr comment_service.CommentDeleteRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithErr(c, err)
	}
	count, err := comment_service.CommentDeleteService(cr)
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("%s", err))
	}
	res.OkWithData(c, fmt.Sprintf("共删除%d条评论 ", count))
}

type CommentByArticleListRequest struct {
	common.PageInfo
	Title string `json:"title" form:"title"`
}

type CommentByArticleListResponse struct {
	Title string `json:"title"`
	ID    string `json:"id"`
	Count int    `json:"count"`
}

// CommentByArticleListView 有评论的文章列表
// @Tags 评论管理
// @Summary 有评论的文章列表
// @Description 有评论的文章列表
// @Param id path string  true  "id"
// @Param data query CommentByArticleListRequest  true  "参数"
// @Router /api/comments/articles [get]
// @Produce json
// @Success 200 {object} res.Response{data=res.DataListResponse[list = CommentByArticleListResponse]}
func (CommentApi) CommentByArticleListView(c *gin.Context) {
	var cr CommentByArticleListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	var count int64

	global.DB.Model(models.CommentModel{}).
		Group("article_id").Count(&count)

	type T struct {
		ArticleID string
		Count     int
	}

	var _list []T
	global.DB.Model(models.CommentModel{}).
		Group("article_id").Order("count desc").Limit(cr.GetLimit()).Offset(cr.GetOffset()).Select("article_id", "count(id) as count").Scan(&_list)
	//logrus.Infof("_list : %#v", _list)
	var articleIDMap = map[string]int{}
	var articleIDList []interface{}
	for _, t := range _list {
		articleIDMap[t.ArticleID] = t.Count
		articleIDList = append(articleIDList, t.ArticleID)
	}

	// 1. 先创建 BoolQuery，用于组合多个查询条件
	boolQuery := elastic.NewBoolQuery()

	// 2. 必须满足的条件：articleID 在之前获取的 articleIDList 内（原有逻辑保留）
	boolQuery.Must(elastic.NewTermsQuery("_id", articleIDList...))

	// 3. 可选条件：如果传了 title，就加 title 模糊匹配（核心新增）
	if cr.Title != "" {
		// 模糊匹配 title 字段：会对搜索词分词，比如“测试”能匹配“测试文章”“文章测试”
		// 如果你需要更严格的“包含”（如“*测试*”），可换成 elastic.NewWildcardQuery("title", "*" + cr.Title + "*")
		boolQuery.Must(elastic.NewMatchQuery("title", cr.Title))
	}

	// 4. ES 搜索时，用上面组合好的 boolQuery 作为查询条件（替换原来的 NewTermsQuery）
	result, err := global.Es.
		Search(models.ArticleModel{}.Index()).
		Query(boolQuery). // 这里换成组合后的 boolQuery
		Size(10000).
		Do(context.Background())

	if err != nil {
		res.FailWithMsg(c, "es查询错误")
		return
	}

	var list = make([]CommentByArticleListResponse, 0)
	for _, hit := range result.Hits.Hits {
		var model models.ArticleModel
		err = json.Unmarshal(hit.Source, &model)
		if err != nil {
			logrus.Error(err)
			continue
		}

		model.ID = hit.Id

		list = append(list, CommentByArticleListResponse{
			Title: model.Title,
			ID:    hit.Id,
			Count: articleIDMap[hit.Id],
		})
	}
	res.OkWithList(c, list, int(count))

	return
}
