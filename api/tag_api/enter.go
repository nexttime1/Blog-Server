package tag_api

import (
	"Blog_server/common"
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/log_service"
	"Blog_server/utils/jwts"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type TagApi struct {
}
type TagRequest struct {
	Title string `json:"title" binding:"required"`
}

type TagListRequest struct {
	common.PageInfo
	Title string `json:"title"`
}

// TagAddView 添加标签
// @Summary 添加标签
// @Description 创建一个新的标签，包含标题
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param data body TagRequest true "标签信息"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response "创建成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/tags [post]
func (TagApi) TagAddView(c *gin.Context) {
	_claims, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claims.(*jwts.MyClaims)
	var mr TagRequest
	err := c.ShouldBindJSON(&mr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log := log_service.GetLog(c)
	log.SetTitle("添加文章标签")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
	var model models.TagModel
	err = global.DB.Where("title = ?", mr.Title).Take(&model).Error
	if err == nil {
		logrus.Errorf("标签已经存在")
		res.FailWithMsg(c, "标签已经存在")
		return
	}
	err = global.DB.Create(&models.TagModel{
		Title: mr.Title,
	}).Error
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("存入数据库失败 %s", err.Error()))
		return
	}
	res.OkWithMessage(c, "创建标签成功")
}

// TagListView 标签列表
// @Summary 获取标签列表
// @Description 分页查询标签列表
// @Tags 标签管理
// @Produce json
// @Param page query int false "页码，默认1" mininum(1)
// @Param limit query int false "每页条数，默认10" mininum(1) maxinum(100)
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{data=res.DataListResponse}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/tags [get]
func (TagApi) TagListView(c *gin.Context) {
	_, exists := c.Get("claims")
	if !exists {
		return
	}
	var mr TagListRequest

	err := c.ShouldBindQuery(&mr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	list, count, err := common.ListQuery(models.TagModel{}, common.Options{
		PageInfo: mr.PageInfo,
		Preload:  []string{"Articles"},
	})
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	res.OkWithList(c, list, count)

}

// TagUpdateView 更新标签
// @Summary 更新标签
// @Description 更新标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param data body TagRequest true "更新的标签"
// @Param id path int true "id"
// @Param token header string true "token"
// @Success 200 {object} res.Response "更新成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 404 {object} res.Response "标签不存在"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/tags/{id} [put]
func (TagApi) TagUpdateView(c *gin.Context) {
	_claims, exists := c.Get("claims")
	if !exists {
	}
	var Request models.IDRequest
	claim := _claims.(*jwts.MyClaims)
	err := c.ShouldBindUri(&Request)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	var mr TagRequest
	err = c.ShouldBindJSON(&mr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log := log_service.GetLog(c)
	log.SetTitle("更新标签")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
	var model models.TagModel
	err = global.DB.Where("id = ?", Request.ID).Take(&model).Error
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("未找到此id %s", err.Error()))
		return
	}
	err = global.DB.Model(&model).Update("title", mr.Title).Error
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("更新失败 %s", err.Error()))
		return
	}
	res.OkWithMessage(c, "更新成功")

}

// TagDeleteView 批量删除标签
// @Summary 批量删除标签
// @Description 批量删除标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Success 200 {object} res.Response "删除成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/tags [delete]
func (TagApi) TagDeleteView(c *gin.Context) {
	_claims, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claims.(*jwts.MyClaims)
	var mr models.RemoveRequest
	err := c.ShouldBindJSON(&mr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log := log_service.GetLog(c)
	log.SetTitle("删除标签")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
	var modelList []models.TagModel
	modelList = common.BatchRemove(modelList, mr)
	count := len(modelList)
	res.OkWithMessage(c, fmt.Sprintf("删除成功 共删除%d个标签", count))

}

type TagResponse struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// TagNameListView 标签名称列表
// @Tags 标签管理
// @Summary 标签名称列表
// @Description 标签名称列表
// @Router /api/tag_names [get]
// @Produce json
// @Success 200 {object} res.Response{data=[]TagResponse}
func (TagApi) TagNameListView(c *gin.Context) {
	type T struct {
		DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
		SumOtherDocCount        int `json:"sum_other_doc_count"`
		Buckets                 []struct {
			Key      string `json:"key"`
			DocCount int    `json:"doc_count"`
		} `json:"buckets"`
	}
	query := elastic.NewBoolQuery() //会匹配索引中所有文档

	// 创建一个terms聚合：按"tags"字段分组，统计每个标签的出现次数
	agg := elastic.NewTermsAggregation().Field("tags")
	result, err := global.Es.
		Search(models.ArticleModel{}.Index()).
		Query(query).
		Aggregation("tags", agg).
		Size(0).
		Do(context.Background())
	if err != nil {
		logrus.Error(err)
		return
	}
	byteData := result.Aggregations["tags"]
	var tagType T
	json.Unmarshal(byteData, &tagType)

	var tagList = make([]TagResponse, 0)
	for _, bucket := range tagType.Buckets {
		tagList = append(tagList, TagResponse{
			Label: bucket.Key,
			Value: bucket.Key,
		})
	}

	res.OkWithData(c, tagList)

}
