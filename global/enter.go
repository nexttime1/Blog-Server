package global

import (
	"Blog_server/conf"
	"github.com/go-redis/redis"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

var (
	Config      *conf.Config
	DB          *gorm.DB
	Redis       *redis.Client
	Es          *elastic.Client
	SettingYaml = "D:\\1111kaoyan111111111111111111111111111111\\go_project\\Blog-Server\\settings.yaml"
)
