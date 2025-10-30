package collect_api

import (
	"Blog_server/common"
	"Blog_server/common/Es_option"
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/article_service"
	"Blog_server/utils/jwts"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type CollectApi struct {
}

type UserCollectResponses struct {
	models.ArticleModel
	CreateAt string `json:"create_at"`
}

// ArticleCollectView 添加或取消收藏
// @Summary 添加或取消收藏
// @Description 添加或取消收藏
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param data body models.EsIdQuest false "文章id"
// @Param token header string true "token"
// @Success 200 {object} res.Response "操作成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/articles/collects [post]
func (CollectApi) ArticleCollectView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	var cr models.EsIdQuest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	//查找 用户存不存在
	var userModel models.UserModel
	err = global.DB.Where("id = ?", claim.UserID).Take(&userModel).Error
	if err != nil {
		res.FailWithMsg(c, "用户不存在")
		return
	}
	//查一下 文章id存不存在
	model, err := Es_option.EsArticleDetailByIdQuery(cr.ID)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	fmt.Println(model)

	var num = -1 //存在就取消收藏  不存在就收藏
	//查 第三张表
	var collectUserModel models.UserCollectModel
	err = global.DB.Where("user_id = ? and article_id = ?", claim.UserID, cr.ID).Take(&collectUserModel).Error
	if err != nil {
		//说明没有
		num = 1
		global.DB.Create(&models.UserCollectModel{
			UserID:    userModel.ID,
			ArticleID: cr.ID,
		})
	} else {
		//查到了  说明要取消收藏
		global.DB.Delete(&collectUserModel)
	}
	err = article_service.ArticleUpdate(cr.ID, map[string]interface{}{
		"collects_count": model.CollectsCount + num,
	})
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("更新文章收藏失败 %v", err))
		return
	}
	if num == 1 {
		res.OkWithMessage(c, "收藏文章成功")
	} else {
		res.OkWithMessage(c, "取消收藏成功")
	}

}

// UserCollectListView 查看个人收藏表
// @Summary 查看个人收藏表
// @Description 分页查询个人收藏表
// @Tags 文章管理
// @Produce json
// @Param page query int false "页码，默认1" mininum(1)
// @Param limit query int false "每页条数，默认10" mininum(1) maxinum(100)
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{data=res.DataListResponse}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/articles/collects [get]
func (CollectApi) UserCollectListView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	var cr common.PageInfo
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	CollectModelList, count, err := common.ListQuery(models.UserCollectModel{
		UserID: claim.UserID,
	}, common.Options{
		PageInfo: cr,
	})
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	var CollectList = make([]interface{}, 0)
	var CollectMap = make(map[string]string)
	for _, model := range CollectModelList {
		CollectList = append(CollectList, model.ArticleID)
		CollectMap[model.ArticleID] = model.CreatedAt.Format("2006-01-02 15:04:05")
	}
	query := elastic.NewTermsQuery("_id", CollectList...)
	result, err := global.Es.Search(models.ArticleModel{}.Index()).
		Query(query).
		Size(count).
		Do(context.Background())
	if err != nil {
		logrus.Errorf("es 查询错误 %s", err)
		res.FailWithErr(c, err)
		return
	}
	var response []UserCollectResponses
	for _, hit := range result.Hits.Hits {
		var articleModel models.ArticleModel
		_ = json.Unmarshal(hit.Source, &articleModel)
		articleModel.ID = hit.Id
		avatar := articleModel.BannerUrl
		articleModel.BannerUrl = "http://127.0.0.1:8080/" + avatar
		response = append(response, UserCollectResponses{
			ArticleModel: articleModel,
			CreateAt:     CollectMap[hit.Id],
		})
	}
	res.OkWithList(c, response, count)

}

// UserCollectBatchDeleteView 批量删除收藏的文章
// @Summary 批量删除收藏的文章
// @Description 批量删除收藏的文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param data body models.EsIdListQuest true "批量删除文章的id列表"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/articles/collects [delete]
func (CollectApi) UserCollectBatchDeleteView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	var cr models.EsIdListQuest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	var CollectList = make([]models.UserCollectModel, 0)
	var idList []string
	count := global.DB.Where("user_id = ? and article_id in ?", claim.UserID, cr.IDList).
		Find(&CollectList).Select("article_id").Scan(&idList).RowsAffected
	if count == 0 {
		// 没有收藏任何东西
		res.FailWithMsg(c, "请求非法")
		return
	}
	var EsList = make([]interface{}, 0)
	for _, id := range idList {
		EsList = append(EsList, id)
	}
	query := elastic.NewTermsQuery("_id", EsList...)
	result, err := global.Es.Search(models.ArticleModel{}.Index()).
		Query(query).
		Size(int(count)).
		Do(context.Background())
	if err != nil {
		logrus.Errorf("es 查询错误 %s", err)
		res.FailWithMsg(c, fmt.Sprintf("es 查询错误 %s", err))
		return
	}
	for _, hit := range result.Hits.Hits {
		var model models.ArticleModel
		_ = json.Unmarshal(hit.Source, &model)
		fmt.Println(model.CollectsCount)
		err = article_service.ArticleUpdate(hit.Id, map[string]interface{}{
			"collects_count": model.CollectsCount - 1,
		})
		fmt.Println(model.CollectsCount - 1)
		if err != nil {
			logrus.Errorf("这不能错啊 %s", err)
			res.FailWithMsg(c, fmt.Sprintf("这不能错啊 %s", err))
			return
		}
	}
	global.DB.Delete(&CollectList)
	res.OkWithMessage(c, fmt.Sprintf("取消收藏%d个的文章", len(CollectList)))

}
