package advert_api

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"fmt"
	"github.com/gin-gonic/gin"
)

type AdvertRequest struct {
	Title  string `json:"title" binding:"required"`   // 显示的标题
	Href   string `json:"href" binding:"required"`    // 跳转链接
	Images string `json:"images" binding:"required"`  // 图片
	IsShow bool   `json:"is_show" binding:"required"` // 是否展示
}

type AdvertApi struct {
}

func (AdvertApi) AdvertAddView(c *gin.Context) {
	var ar AdvertRequest
	err := c.ShouldBindJSON(&ar)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	err = global.DB.Create(&models.AdvertModel{
		Title:  ar.Title,
		Href:   ar.Href,
		Images: ar.Images,
		IsShow: ar.IsShow,
	}).Error
	if err != nil {
		res.FailWithErr(c, fmt.Errorf("添加广告失败  %v", err))
		return
	}
	res.OkWithMessage(c, "添加成功")

}
