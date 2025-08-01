package menu_api

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/menu_service"
	"github.com/gin-gonic/gin"
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
	var mc menu_service.MenuAddRequest
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
