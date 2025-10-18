package api

import (
	"Blog_server/api/advert_api"
	"Blog_server/api/article_api"
	"Blog_server/api/big_model_api"
	"Blog_server/api/chat_api"
	"Blog_server/api/collect_api"
	"Blog_server/api/comment_api"
	"Blog_server/api/date_api"
	"Blog_server/api/digg_api"
	"Blog_server/api/feedback_api"
	"Blog_server/api/gaode_api"
	"Blog_server/api/image_api"
	"Blog_server/api/log_api"
	"Blog_server/api/menu_api"
	"Blog_server/api/message_api"
	"Blog_server/api/new_api"
	"Blog_server/api/role_api"
	"Blog_server/api/settings_api"
	"Blog_server/api/site_api"
	"Blog_server/api/tag_api"
	"Blog_server/api/user_api"
)

type Api struct {
	SiteApi     site_api.SiteApi
	LogApi      log_api.LogApi
	SettingApi  settings_api.SettingApi
	ImageApi    image_api.ImageApi
	AdvertApi   advert_api.AdvertApi
	MenuApi     menu_api.MenuApi
	UserApi     user_api.UserApi
	TagApi      tag_api.TagApi
	MessageApi  message_api.MessageApi
	ArticleApi  article_api.ArticleApi
	DiggApi     digg_api.DiggApi
	CollectApi  collect_api.CollectApi
	CommentApi  comment_api.CommentApi
	NewApi      new_api.NewApi
	ChatApi     chat_api.ChatApi
	DateApi     date_api.DateApi
	RoleApi     role_api.RoleApi
	FeedBackApi feedback_api.FeedBackApi
	GaodeApi    gaode_api.GaodeApi
	BigModelApi big_model_api.BigModelApi
}

var App = Api{}
