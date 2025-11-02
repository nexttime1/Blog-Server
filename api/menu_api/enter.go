package menu_api

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/log_service"
	"Blog_server/service/menu_service"
	"Blog_server/utils/jwts"
	"Blog_server/utils/struct_to_map"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
// @Param data body menu_service.MenuRequest true "菜单信息"
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response "操作成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/menus [post]
func (MenuApi) MenuCreateView(c *gin.Context) {
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	var mc menu_service.MenuRequest
	err := c.ShouldBindJSON(&mc)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log := log_service.GetLog(c)
	log.SetTitle("添加菜单")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
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
	err := global.DB.Preload("BannerModel").Order("menu_id ASC, sort DESC").Where("menu_id in ?", MenuIdList).Find(&MenuBannerList).Error
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	for _, model := range MenuList {
		var banners []Banners
		high := len(MenuBannerList) - 1
		low := 0
		for low < high {
			mid := (low + high) / 2
			if MenuBannerList[mid].MenuID >= model.ID {
				high = mid
			} else {
				low = mid + 1
			}
		}
		// 增加判断，确认找到的元素是否符合条件
		for i := low; i < len(MenuBannerList); i++ {
			if MenuBannerList[i].MenuID != model.ID {
				break
			}
			banner := MenuBannerList[i]
			banners = append(banners, Banners{
				ID:   banner.BannerID,
				Path: "http://127.0.0.1:8080/" + banner.BannerModel.Path,
			})
		}
		menus = append(menus, MenuListResponse{
			MenuModel: model,
			Banners:   banners,
		})
	}
	res.OkWithList(c, menus, len(menus))
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
// @Router /api/menus_names [get]
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
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	id := c.Param("id")
	var mc menu_service.MenuRequest
	err := c.ShouldBindJSON(&mc)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log := log_service.GetLog(c)
	log.SetTitle("菜单更新")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
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
	_claim, exists := c.Get("claims")
	if !exists {
		return
	}
	claim := _claim.(*jwts.MyClaims)
	var deleteRequest models.RemoveRequest
	err := c.ShouldBindJSON(&deleteRequest)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	log := log_service.GetLog(c)
	log.SetTitle("菜单删除")
	log.SetRequest(c)
	log.ShowRequest()
	log.SetItem("操作用户", claim.Username)
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
// @Success 200 {object} res.Response{data = MenuBannerResponse} "成功返回横幅配置"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/menus/detail[get]
func (MenuApi) MenuDetailView(c *gin.Context) {
	path := c.Query("path")
	var model models.MenuModel
	err := global.DB.Where("path = ?", path).Take(&model).Error
	if err != nil {
		res.FailWithErr(c, err)
		return
	}

	var menuBannerModel []models.MenuBannerModel
	global.DB.Preload("BannerModel").Order("sort desc").Find(&menuBannerModel, "menu_id = ?", model.ID)
	var banners []Banners
	for _, banner := range menuBannerModel {
		banners = append(banners, Banners{
			ID:   banner.BannerID,
			Path: banner.BannerModel.Path,
		})
	}

	response := MenuBannerResponse{
		Slogan:     model.Slogan,     // 标语（前端data.slogan）
		Abstract:   model.Abstract,   // 描述文字（前端data.abstract，支持字符串或数组）
		BannerTime: model.BannerTime, // 轮播间隔时间（秒，前端data.banner_time）
		Banners:    banners,          // 轮播图列表（前端data.banners）
	}
	logrus.Infof("response: %v", response)
	res.OkWithData(c, response)

}

type MenuBannerResponse struct {
	Slogan     string      `json:"slogan"`      // 标语（必填）
	Abstract   interface{} `json:"abstract"`    // 描述（支持string或[]string）
	BannerTime int         `json:"banner_time"` // 轮播间隔时间（秒，可选，前端默认7）
	Banners    []Banners   `json:"banners"`     // 轮播图（
}
