package advert_api

import (
	"Blog_server/common"
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/log_service"
	"Blog_server/utils/jwts"
	"Blog_server/utils/struct_to_map"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type AdvertRequest struct { //json
	Title  string `json:"title" binding:"required" `      // 显示的标题
	Href   string `json:"href" binding:"required,url" `   // 跳转链接
	Images string `json:"images" binding:"required,url" ` // 图片
	IsShow bool   `json:"is_show" `                       // 是否展示
}

type AdvertListResponse struct {
	common.PageInfo
	Title  string `json:"title"`   // 显示的标题
	Href   string `json:"href"`    // 跳转链接
	Images string `json:"images"`  // 图片
	IsShow bool   `json:"is_show"` // 是否展示

}
type AdvertUpdateRequest struct {
	Title  string `json:"title" structs:"title"`                                // 显示的标题
	Href   string `json:"href" structs:"href"`                                  // 跳转链接
	Images string `json:"images" structs:"images"`                              // 图片
	IsShow *bool  `json:"is_show,omitempty" structs:"is_show" gorm:"omitempty"` // 是否展示

}

type AdvertApi struct {
}

// AdvertAddView 添加广告
// @Summary 创建广告
// @Description 创建一个新的广告条目，包含标题、跳转链接、图片和展示状态
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param data body AdvertRequest true "广告信息"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response "操作成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/adverts [post]
func (AdvertApi) AdvertAddView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	var ar AdvertRequest
	err := c.ShouldBindJSON(&ar)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	//记录日志
	log := log_service.GetLog(c)
	log.SetTitle("创建广告")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
	//先查看是否重复
	var advert models.AdvertModel
	err = global.DB.Take(&advert, "title = ?", ar.Title).Error
	if err == nil {
		//说明已经存在
		log.SetItemError("广告已存在", advert.Title)
		res.FailWithErr(c, errors.New("广告已存在"))
		return
	}

	err = global.DB.Create(&models.AdvertModel{
		Title:  ar.Title,
		Href:   ar.Href,
		Images: ar.Images,
		IsShow: ar.IsShow,
	}).Error
	if err != nil {
		log.SetItemError("添加广告失败", err)
		res.FailWithErr(c, fmt.Errorf("添加广告失败  %v", err))
		return
	}
	log.SetItemInfo("广告标题", ar.Title)
	res.OkWithMessage(c, "添加成功")

}

// AdvertListView 广告列表
// @Summary 获取广告列表
// @Description 分页查询广告列表，支持根据标题、链接等条件筛选
// @Tags 广告管理
// @Produce json
// @Param title query string false "标题筛选（模糊匹配）"
// @Param href query string false "链接筛选（模糊匹配）"
// @Param images query string false "图片筛选（模糊匹配）"
// @Param is_show query bool false "是否展示（true/false）"
// @Param page query int false "页码，默认1" mininum(1)
// @Param limit query int false "每页条数，默认10" mininum(1) maxinum(100)
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{data=res.DataListResponse}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/adverts [get]
func (AdvertApi) AdvertListView(c *gin.Context) {
	var as AdvertListResponse
	err := c.ShouldBindQuery(&as)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	as.IsShow = true
	StringData := c.GetHeader("Referer")
	if strings.Contains(StringData, "admin") {
		// admin 来的
		as.IsShow = false
	}
	list, count, err := common.ListQuery(models.AdvertModel{
		Title:  as.Title,
		Href:   as.Href,
		Images: as.Images,
		IsShow: as.IsShow,
	}, common.Options{
		PageInfo:     as.PageInfo,
		Debug:        true,
		DefaultOrder: "created_at DESC",
	})
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithList(c, list, count)

}

// AdvertUpdateView 广告更新
// @Summary 更新广告
// @Description 更新广告
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param data body AdvertUpdateRequest false "更新的广告信息（可选字段，如title/href等，不传则不更新）"
// @Param id path int true "id"
// @Param token header string true "token"
// @Success 200 {object} res.Response "更新成功"
// @Failure 400 {object} res.Response "请求参数错误（如ID格式错误、标题重复等）"
// @Failure 404 {object} res.Response "广告不存在"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/adverts/{id} [put]
func (AdvertApi) AdvertUpdateView(c *gin.Context) {
	id := c.Param("id")

	var ac AdvertUpdateRequest
	err := c.ShouldBindJSON(&ac)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	var UpdateAdvert models.AdvertModel
	err = global.DB.Take(&UpdateAdvert, "id = ?", id).Error
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	// 判断 title 唯一    如果重复  不嫩修改.
	var advert models.AdvertModel
	err = global.DB.Debug().Where("title = ? and id != ?", ac.Title, id).Take(&advert).Error
	if err == nil {
		res.FailWithErr(c, errors.New("title 值重复"))
		return
	}
	//把结构体变成map  把没传值的 都给 去掉  不去修改 数据库
	fmt.Println(ac)
	toMap := struct_to_map.StructToMap(ac)
	if toMap == nil {
		toMap["is_show"] = UpdateAdvert.IsShow
	}

	//修改
	err = global.DB.Model(&UpdateAdvert).Updates(toMap).Error
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithMessage(c, "修改成功")
}

// AdvertDeleteView 广告删除
// @Summary 批量删除广告
// @Description 批量删除广告
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param data body models.RemoveRequest true "广告id列表"
// @Param token header string true "token"
// @Success 200 {object} res.Response "删除成功"
// @Failure 400 {object} res.Response "请求参数错误（如id格式错误、列表为空）"
// @Failure 404 {object} res.Response "部分广告id不存在"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/adverts [delete]
func (AdvertApi) AdvertDeleteView(c *gin.Context) {
	var removeRequest models.RemoveRequest
	err := c.ShouldBindJSON(&removeRequest)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	var advertModels []models.AdvertModel
	advertModels = common.BatchRemove(advertModels, removeRequest)
	res.OkWithMessage(c, fmt.Sprintf("成功删除%d条广告", len(advertModels)))
}
