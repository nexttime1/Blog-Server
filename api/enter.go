package api

import (
	"Blog_server/api/log_api"
	"Blog_server/api/settings_api"
	"Blog_server/api/site_api"
)

type Api struct {
	SiteApi    site_api.SiteApi
	LogApi     log_api.LogApi
	SettingApi settings_api.SettingApi
}

var App = Api{}
