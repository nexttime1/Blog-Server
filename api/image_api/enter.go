package image_api

import (
	"Blog_server/common"
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/image_service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
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

// ImageUploadView 上传图片
func (ImageApi) ImageUploadView(c *gin.Context) {
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
	// 记录n个照片上传
	var responseList []image_service.ImageListResponse

	for _, FileHeader := range FileHeaderList {
		response, err := image_service.UploadService(c, FileHeader)
		if err != nil {
			res.FailWithErr(c, err)
			return
		}
		responseList = append(responseList, response)
	}
	res.OkWithData(c, responseList)
}

// ImageInfoView 查看图片
func (ImageApi) ImageInfoView(c *gin.Context) {
	var cr ImageListView
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	list, count, err := common.ListQuery(models.BannerModel{
		Path: cr.Path,
		Hash: cr.Hash,
		Name: cr.Name,
	}, common.Options{
		PageInfo:     cr.PageInfo,
		Debug:        true,
		DefaultOrder: "created_at DESC", //写死了  默认降序排序  前端修改的话
	})
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithList(c, list, count)

}

// ImageRemoveView 批量删除图片
func (ImageApi) ImageRemoveView(c *gin.Context) {
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	var ModelList []models.BannerModel
	ModelList = common.BatchRemove(ModelList, cr)

	msg := fmt.Sprintf("图片删除成功，共删除%d条数据", len(ModelList))

	res.OkWithMessage(c, msg)

}

func (ImageApi) ImageUpdateView(c *gin.Context) {
	var cr ImageUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
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
