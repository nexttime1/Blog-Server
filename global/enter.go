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
	LevelFlag   bool
	SettingYaml = "D:\\5524\\go_project\\Blog-Server\\settings.yaml"
)
