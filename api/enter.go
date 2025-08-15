package api

import (
	"Blog_server/api/advert_api"
	"Blog_server/api/article_api"
	"Blog_server/api/digg_api"
	"Blog_server/api/image_api"
	"Blog_server/api/log_api"
	"Blog_server/api/menu_api"
	"Blog_server/api/message_api"
	"Blog_server/api/settings_api"
	"Blog_server/api/site_api"
	"Blog_server/api/tag_api"
	"Blog_server/api/user_api"
)

type Api struct {
	SiteApi    site_api.SiteApi
	LogApi     log_api.LogApi
	SettingApi settings_api.SettingApi
	ImageApi   image_api.ImageApi
	AdvertApi  advert_api.AdvertApi
	MenuApi    menu_api.MenuApi
	UserApi    user_api.UserApi
	TagApi     tag_api.TagApi
	MessageApi message_api.MessageApi
	ArticleApi article_api.ArticleApi
	DiggApi    digg_api.DiggApi
}

var App = Api{}
