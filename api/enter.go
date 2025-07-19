package api

import (
	"Blog_server/api/log_api"
	"Blog_server/api/site_api"
)

type Api struct {
	SiteApi site_api.SiteApi
	LogApi  log_api.LogApi
}

var App = Api{}
