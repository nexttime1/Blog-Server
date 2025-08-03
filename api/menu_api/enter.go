package menu_api

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/menu_service"
	"Blog_server/utils/struct_to_map"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MenuApi struct {
}

type Banners struct {
	ID   uint   `json:"id"`
	Path string `json:"path"`
}
type MenuListResponse struct {
	models.MenuModel
	Banners []Banners
}
type MenuNameListResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Path  string `json:"path"`
}

// MenuCreateView 添加菜单
// @Summary 添加菜单
// @Description 创建一个新的菜单，包含Title Slogan Abstract等
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param data body menu_service.MenuAddRequest true "菜单信息"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response "操作成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/menus [post]
func (MenuApi) MenuCreateView(c *gin.Context) {
	var mc menu_service.MenuRequest
	err := c.ShouldBindJSON(&mc)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	err, menuBannerList := menu_service.MenuAddService(c, mc)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	if menuBannerList == nil {
		res.OkWithMessage(c, "菜单添加成功")
		return
	}
	res.OkWithData(c, menuBannerList)
}

// MenuListView 菜单列表
// @Summary 获取菜单列表
// @Description 等修改
// @Tags 菜单管理
// @Produce json
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{data=[]MenuListResponse}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/menus [get]
func (MenuApi) MenuListView(c *gin.Context) {
	var MenuList []models.MenuModel
	var MenuIdList []uint
	global.DB.Order("sort desc").Find(&MenuList).Select("id").Scan(&MenuIdList)

	var MenuBannerList []models.MenuBannerModel
	var menus []MenuListResponse
	err := global.DB.Preload("BannerModel").Order("sort desc").Where("menu_id in ?", MenuIdList).Find(&MenuBannerList).Error
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	for _, model := range MenuList {
		var banners []Banners
		for _, banner := range MenuBannerList {
			if model.ID != banner.MenuID {
				continue
			}
			banners = append(banners, Banners{
				ID:   banner.BannerID,
				Path: banner.BannerModel.Path,
			})
		}
		menus = append(menus, MenuListResponse{
			MenuModel: model,
			Banners:   banners,
		})
	}
	res.OkWithData(c, menus)
}

// MenuNameListView 菜单名称列表
// @Summary 获取菜单名称列表
// @Description 等修改
// @Tags 菜单管理
// @Produce json
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{data=[]MenuNameListResponse}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/menus_name [get]
func (MenuApi) MenuNameListView(c *gin.Context) {
	var mr []MenuNameListResponse
	err := global.DB.Model(models.MenuModel{}).Select("id, title, path").Scan(&mr).Error
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithData(c, mr)
}

// MenuUpdateView 菜单更新
// @Summary 菜单更新
// @Description 菜单更新
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param data body menu_service.MenuRequest false "更新的菜单信息（可选字段，如title/path等，不传则不更新）"
// @Param id path int true "id"
// @Param token header string true "token"
// @Success 200 {object} res.Response "更新成功"
// @Failure 400 {object} res.Response "请求参数错误（如ID格式错误、标题重复等）"
// @Failure 404 {object} res.Response "菜单不存在"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/menus/{id} [put]
func (MenuApi) MenuUpdateView(c *gin.Context) {
	id := c.Param("id")
	var mc menu_service.MenuRequest
	err := c.ShouldBindJSON(&mc)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	var model models.MenuModel
	err = global.DB.Where("id = ?", id).Take(&model).Error
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("未找到models.MenuModel 的 id %v", err))
		return
	}
	// 先清空表
	err = global.DB.Model(&model).Association("Banners").Clear()
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	if len(mc.ImageSortList) > 0 {
		//说明修改图片了
		var menuBannerModel []models.MenuBannerModel
		for _, sort := range mc.ImageSortList {
			menuBannerModel = append(menuBannerModel, models.MenuBannerModel{
				MenuID:   model.ID,
				BannerID: sort.ImageID,
				Sort:     sort.Sort,
			})
		}
		err := global.DB.Create(&menuBannerModel).Error
		if err != nil {
			res.FailWithMsg(c, fmt.Sprintf("第三张表 修改失败  %s", err))
			return
		}
	}

	// 普通修改
	maps := struct_to_map.StructToMap(mc)
	err = global.DB.Model(&model).Updates(maps).Error
	if err != nil {
		res.FailWithMsg(c, fmt.Sprintf("修改失败  %s", err))
		return
	}
	res.OkWithMessage(c, "修改成功")
}

// MenuDeleteView 菜单删除
// @Summary 批量菜单删除
// @Description 批量菜单删除
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param data body models.RemoveRequest true "菜单id列表"
// @Param token header string true "token"
// @Success 200 {object} res.Response "删除成功"
// @Failure 400 {object} res.Response "请求参数错误（如id格式错误、列表为空）"
// @Failure 404 {object} res.Response "部分菜单id不存在"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/menus [delete]
func (MenuApi) MenuDeleteView(c *gin.Context) {
	var deleteRequest models.RemoveRequest
	err := c.ShouldBindJSON(&deleteRequest)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	if len(deleteRequest.IDList) == 0 {
		res.FailWithMsg(c, "请给出 删除的菜单列表")
		return
	}

	var modelList []models.MenuModel
	count := global.DB.Where("id in ?", deleteRequest.IDList).Find(&modelList).RowsAffected
	if count == 0 {
		res.FailWithMsg(c, "菜单不存在")
		return
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&modelList).Association("Banners").Clear()
		if err != nil {
			return err
		}
		err = tx.Delete(&modelList).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithMessage(c, fmt.Sprintf("删除成功 共删除 %d 条数据", count))
}

// MenuDetailView 某个菜单
// @Summary 获取某个菜单
// @Description 获取某个菜单
// @Tags 菜单管理
// @Produce json
// @Param id path int true "id"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{data=[]MenuListResponse}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/menus/{id} [get]
func (MenuApi) MenuDetailView(c *gin.Context) {
	id := c.Param("id")
	var model models.MenuModel
	err := global.DB.Where("id = ?", id).Take(&model).Error
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	var menuBannerModel []models.MenuBannerModel
	global.DB.Preload("BannerModel").Order("sort desc").Find(&menuBannerModel, "id = ?", id)
	var banners []Banners
	for _, banner := range menuBannerModel {
		banners = append(banners, Banners{
			ID:   banner.BannerID,
			Path: banner.BannerModel.Path,
		})
	}

	response := MenuListResponse{
		MenuModel: model,
		Banners:   banners,
	}
	res.OkWithData(c, response)

}
