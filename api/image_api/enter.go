package image_api

import (
	"Blog_server/common"
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/image_service"
	"Blog_server/service/log_service"
	"Blog_server/utils/jwts"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type ImageApi struct {
}

// ImageListView 接受 查询图片 参数
type ImageListView struct {
	common.PageInfo
	Path string `form:"path"` // 图片路径
	Hash string `form:"hash"` // 图片的hash值，用于判断重复图片
	Name string `form:"name"` // 图片名称
}

type ImageUpdateRequest struct {
	ID   uint   `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type ImageNameListResponse struct {
	ID   uint   `json:"id"`
	Path string `json:"path"` // 图片路径
	Name string `json:"name"` // 图片名称
}

type ImageListInfoResponse struct {
	models.BannerModel
	BannerCount  int `json:"bannerCount"`
	ArticleCount int `json:"articleCount"`
}

// ImageUploadView 图片上传
// @Summary 批量上传图片
// @Description 支持多图片上传，自动验证格式（白名单）和大小，重复图片会被拦截
// @Tags 图片管理
// @Accept multipart/form-data
// @Produce json
// @Param token header string true "用户认证令牌"
// @Param image formData file true "图片文件列表（支持多文件，格式：jpg/jpeg/png/gif等）"
// @Success 200 {object} res.Response{data=[]image_service.ImageListResponse} "上传结果列表（包含每个文件的上传状态、路径等信息）"
// @Failure 400 {object} res.Response "请求错误（如文件不存在、格式错误、大小超限等）"
// @Failure 500 {object} res.Response "服务器错误（如上传七牛云失败、保存文件失败等）"
// @Router /api/images [post]
func (ImageApi) ImageUploadView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	form, err := c.MultipartForm()
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	FileHeaderList, ok := form.File["image"]
	if !ok {
		res.FailWithErr(c, errors.New("上传的文件图片不存在"))
		return
	}
	log := log_service.GetLog(c)
	log.SetTitle("图片上传")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
	// 记录n个照片上传
	var responseList []image_service.ImageListResponse

	for _, FileHeader := range FileHeaderList {
		response, err := image_service.UploadService(c, FileHeader)
		if err != nil {
			res.FailWithErr(c, err)
			return
		}
		FilePath := "http://127.0.0.1:8080/" + response.FilePath
		response.FilePath = FilePath
		logrus.Infof("response.FilePath %v", response.FilePath)
		responseList = append(responseList, response)
	}
	log.SetItemInfo("上传图片个数为", len(responseList))
	res.OkWithMessage(c, "上传成功")
}

// ImageInfoView 查看图片  无  @Accept
// @Summary 查看图片
// @Description 查询对应参数的图片
// @Tags 图片管理
// @Produce json
// @Param path query string false "输入对对应参数"
// @Param hash query string false "输入对对应参数"
// @Param name query string false "输入对对应参数"
// @Param page query int false "页码，默认1" mininum(1)
// @Param limit query int false "每页条数，默认10" mininum(1) maxinum(100)
// @Param token header string true "token"
// @Success 200 {object} res.Response{data=res.DataListResponse} "对于图片的各个信息"
// @Router /api/images [get]
func (ImageApi) ImageInfoView(c *gin.Context) {
	var cr ImageListView
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	list_, count, err := common.ListQuery(models.BannerModel{
		Path: cr.Path,
		Hash: cr.Hash,
		Name: cr.Name,
	}, common.Options{
		PageInfo:     cr.PageInfo,
		Debug:        true,
		Likes:        []string{"name"},
		DefaultOrder: "created_at DESC", //写死了  默认降序排序  前端修改的话
	})

	var imageIDList []interface{}
	for _, model := range list_ {
		imageIDList = append(imageIDList, model.ID)
	}
	res1, err := global.Es.
		Search(models.ArticleModel{}.Index()).
		Query(elastic.NewTermsQuery("banner_id", imageIDList...)).
		Size(10000).
		Do(context.Background())
	if err != nil {
		logrus.Errorf("%v", err)
		return
	}
	var imageIDArticleMap = map[uint]int{}
	for _, hit := range res1.Hits.Hits {
		var model models.ArticleModel
		err = json.Unmarshal(hit.Source, &model)
		if err != nil {
			logrus.Error(err)
			continue
		}
		val, ok := imageIDArticleMap[model.BannerID]
		if !ok {
			imageIDArticleMap[model.BannerID] = 1
		} else {
			imageIDArticleMap[model.BannerID] = val + 1
		}
	}

	var response []ImageListInfoResponse
	for _, model := range list_ {
		url := model.Path
		model.Path = "http://127.0.0.1:8080/" + url
		// 求于 menus 管联的  也就是 这个图片 有几个menus用
		var MenusCount int64
		global.DB.Model(models.MenuBannerModel{}).Where("banner_id = ?", model.ID).Count(&MenusCount)

		response = append(response, ImageListInfoResponse{
			BannerModel:  model,
			BannerCount:  int(MenusCount),
			ArticleCount: imageIDArticleMap[model.ID],
		})

	}

	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithList(c, response, count)

}

// ImageRemoveView 批量删除图片
// @Summary 批量删除图片
// @Description 批量删除图片
// @Tags 图片管理
// @Accept json
// @Produce json
// @Param data body models.RemoveRequest true "图片id列表"
// @Param token header string true "token"
// @Success 200 {object} res.Response "删除成功"
// @Failure 400 {object} res.Response "请求参数错误（如id格式错误、列表为空）"
// @Failure 404 {object} res.Response "部分图片id不存在"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/images [delete]
func (ImageApi) ImageRemoveView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log := log_service.GetLog(c)
	log.SetTitle("删除图片")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
	var ModelList []models.BannerModel
	ModelList = common.BatchRemove(ModelList, cr)

	msg := fmt.Sprintf("图片删除成功，共删除%d条数据", len(ModelList))

	res.OkWithMessage(c, msg)

}

// ImageUpdateView 更新图片
// @Summary 更新图片
// @Description 根据给出的id 去更新Name
// @Tags 图片管理
// @Accept json
// @Produce json
// @Param data body ImageUpdateRequest true "输入要修改的id和修改后的名字"
// @Param token header string true "token"
// @Success 200 {object} res.Response "更新成功"
// @Failure 400 {object} res.Response "请求参数错误（如ID格式错误、标题重复等）"
// @Failure 404 {object} res.Response "图片不存在"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/images [put]
func (ImageApi) ImageUpdateView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	var cr ImageUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log := log_service.GetLog(c)
	log.SetTitle("更新图片")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
	var model models.BannerModel
	err = global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithErr(c, errors.New(fmt.Sprintf("未查找到该图片 %s", err)))
		return
	}
	err = global.DB.Model(&model).Update("name", cr.Name).Error
	if err != nil {
		res.FailWithErr(c, errors.New(fmt.Sprintf("图片修改失败 %s", err.Error())))
		return
	}
	res.OkWithMessage(c, "图片名称修改成功")

}

// ImageNameListView 只返回  三个参数 的查询图片
// @Summary 图片名称列表
// @Description 只返回一些重要的图片信息
// @Tags 图片管理
// @Produce json
// @Param token header string true "token"
// @Success 200 {object} res.Response{data=[]ImageNameListResponse}
// @Router /api/image_names [get]
func (ImageApi) ImageNameListView(c *gin.Context) {
	var List []ImageNameListResponse
	err := global.DB.Model(&models.BannerModel{}).Select("id", "path", "name").Scan(&List).Error
	var imagesList []ImageNameListResponse
	for _, response := range List {
		path := "http://127.0.0.1/" + response.Path
		response.Path = path
		imagesList = append(imagesList, response)
	}
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithData(c, imagesList)
}
