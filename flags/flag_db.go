package flags

import (
	"Blog_server/global"
	"Blog_server/models"
	"github.com/sirupsen/logrus"
)

func FlagDB() {
	err := global.DB.AutoMigrate(
		&models.UserModel{},
		&models.UserConfModel{},
		&models.ArticleModel{},
		&models.CategoryModel{},
		&models.ArticleDiggModel{},
		&models.CollectModel{},
		&models.UserArticleCollectModel{},
		&models.UserArticleLookHistoryModel{}, //用户浏览历史表
		&models.CommentModel{},
		&models.BannerModel{},
		&models.LogModel{},
		&models.UserLoginModel{},
		&models.GlobalNotificationModel{},
	)
	if err != nil {
		logrus.Errorf("\n数据库迁移失败  %s", err)
		return
	}
	logrus.Info("\n数据库迁移成功")

}
