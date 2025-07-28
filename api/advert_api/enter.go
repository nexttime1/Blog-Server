package advert_api

import (
	"Blog_server/common"
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/log_service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type AdvertRequest struct { //json
	Title  string `json:"title" binding:"required"`      // 显示的标题
	Href   string `json:"href" binding:"required,url"`   // 跳转链接
	Images string `json:"images" binding:"required,url"` // 图片
	IsShow bool   `json:"is_show"`                       // 是否展示
}

type AdvertListResponse struct {
	common.PageInfo
	Title  string `json:"title"`   // 显示的标题
	Href   string `json:"href"`    // 跳转链接
	Images string `json:"images"`  // 图片
	IsShow bool   `json:"is_show"` // 是否展示

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
	//记录日志
	log := log_service.GetLog(c)
	log.ShowRequest()
	log.ShowResponse()
	//先查看是否重复
	var advert models.AdvertModel
	err = global.DB.Take(&advert, "title = ?", ar.Title).Error
	if err == nil {
		//说明已经存在
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
		res.FailWithErr(c, fmt.Errorf("添加广告失败  %v", err))
		return
	}
	res.OkWithMessage(c, "添加成功")

}

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

func (AdvertApi) AdvertUpdateView(c *gin.Context) {
	id := c.Param("id")

	var ac AdvertListResponse
	err := c.ShouldBindJSON(&ac)
	if err != nil {
		res.FailWithErr(c, err)
	}
	var UpdateAdvert models.AdvertModel
	err = global.DB.Take(&UpdateAdvert, "id = ?", id).Error
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	// 判断 title 唯一    如果重复  不嫩修改.
	var advert models.AdvertModel
	err = global.DB.Take(&advert, fmt.Sprintf("title = ? and id != %s", id), ac.Title).Error
	if err == nil {
		res.FailWithErr(c, errors.New("title 值重复"))
		return
	}

	//修改
	err = global.DB.Model(&UpdateAdvert).Updates(map[string]interface{}{
		"title":   ac.Title,
		"href":    ac.Href,
		"images":  ac.Images,
		"is_show": ac.IsShow,
	}).Error
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithMessage(c, "修改成功")
}
