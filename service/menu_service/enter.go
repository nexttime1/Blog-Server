package menu_service

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ImageSort struct {
	ImageID uint `json:"image_id"`
	Sort    int  `json:"sort"`
}

type MenuAddRequest struct {
	Title         string      `json:"title" binding:"required" structs:"title"`
	Path          string      `json:"path" binding:"required" structs:"path"`
	Slogan        string      `json:"slogan" structs:"slogan"`
	Abstract      enum.Array  `json:"abstract" structs:"abstract"`
	AbstractTime  int         `json:"abstract_time" structs:"abstract_time"`  // 切换的时间，单位秒
	BannerTime    int         `json:"banner_time" structs:"banner_time"`      // 切换的时间，单位秒
	Sort          int         `json:"sort" binding:"required" structs:"sort"` // 菜单的序号
	ImageSortList []ImageSort `json:"image_sort_list" structs:"-"`            // 添加的图片和顺序
}

func MenuAddService(c *gin.Context, mc MenuAddRequest) (error, []models.MenuBannerModel) {
	//重复着 title 不能重复
	var model models.MenuModel
	err := global.DB.Where("title = ? or path = ?", mc.Title, mc.Path).Take(&model).Error
	if err == nil {
		//找到了 不行
		return fmt.Errorf("菜单 title 或者 path 值已经存在"), nil
	}

	menuModel := models.MenuModel{
		Title:        mc.Title,
		Path:         mc.Path,
		Slogan:       mc.Slogan,
		Abstract:     mc.Abstract,
		AbstractTime: mc.AbstractTime,
		BannerTime:   mc.BannerTime,
		Sort:         mc.Sort,
	}
	err = global.DB.Create(&menuModel).Error
	if err != nil {
		return fmt.Errorf("添加失败  %s", err), nil
	}
	// 如果不添加图片 就不需要操作第三张表  直接成功
	if len(mc.ImageSortList) == 0 {
		return nil, nil
	}

	//到现在说明 要构建第三张表
	var menuBannerList []models.MenuBannerModel
	for _, sort := range mc.ImageSortList {
		//检查图片有没有
		var imagemodel models.BannerModel
		err := global.DB.Where("id = ?", sort.ImageID).Take(&imagemodel).Error
		if err != nil {
			return err, nil
		}
		menuBannerList = append(menuBannerList, models.MenuBannerModel{
			MenuID:   menuModel.ID,
			BannerID: sort.ImageID,
			Sort:     sort.Sort,
		})
	}
	//没问题 真正去创建
	err = global.DB.Create(&menuBannerList).Error
	if err != nil {
		res.FailWithErr(c, err)
		return err, nil
	}
	return nil, menuBannerList
}
