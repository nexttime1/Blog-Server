package menu_api

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"fmt"
	"github.com/gin-gonic/gin"
)

type MenuApi struct {
}

type ImageSort struct {
	ImageID uint `json:"image_id"`
	Sort    int  `json:"sort"`
}
type MenuAddRequest struct {
	Title       string `json:"title" binding:"required" structs:"title"`
	MenuTitleEn string `json:"menu_title_en" binding:"required" structs:"menu_title_en"`
	//Path          string      `json:"path" binding:"required" structs:"path"`
	Slogan        string      `json:"slogan" structs:"slogan"`
	Abstract      enum.Array  `json:"abstract" structs:"abstract"`
	AbstractTime  int         `json:"abstract_time" structs:"abstract_time"`  // 切换的时间，单位秒
	BannerTime    int         `json:"banner_time" structs:"banner_time"`      // 切换的时间，单位秒
	Sort          int         `json:"sort" binding:"required" structs:"sort"` // 菜单的序号
	ImageSortList []ImageSort `json:"image_sort_list" structs:"-"`            // 添加的图片和顺序
}

func (MenuApi) MenuCreateView(c *gin.Context) {
	var mc MenuAddRequest
	err := c.ShouldBindJSON(&mc)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	//重复着 title 不能重复
	var model models.MenuModel
	err = global.DB.Where("title = ?", mc.Title).Take(&model).Error
	if err == nil {
		//找到了 不行
		res.FailWithMsg(c, "菜单title 值已经存在")
		return
	}

	menuModel := models.MenuModel{
		MenuTitle:    mc.Title,
		MenuTitleEn:  mc.MenuTitleEn,
		Slogan:       mc.Slogan,
		Abstract:     mc.Abstract,
		AbstractTime: mc.AbstractTime,
		BannerTime:   mc.BannerTime,
		Sort:         mc.Sort,
	}
	err = global.DB.Create(&menuModel).Error
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("添加失败 %s", err))
		return
	}
	// 如果不添加图片 就不需要操作第三张表  直接成功
	if len(mc.ImageSortList) == 0 {
		res.OkWithMessage(c, "菜单添加成功")
		return
	}

	//到现在说明 要构建第三张表
	var menuBannerList []models.MenuBannerModel

}
