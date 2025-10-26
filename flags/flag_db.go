package flags

import (
	"Blog_server/global"
	"Blog_server/models"
	"github.com/sirupsen/logrus"
)

func FlagDB() {
	//global.DB.SetupJoinTable(&models.UserModel{}, "CollectsModels", &models.UserCollectModel{})
	global.DB.SetupJoinTable(&models.MenuModel{}, "Banners", &models.MenuBannerModel{})
	err := global.DB.AutoMigrate(
		//&models_new.UserModel{},
		//&models_new.UserConfModel{},
		//&models_new.ArticleModel{},
		//&models_new.CategoryModel{},
		//&models_new.ArticleDiggModel{},
		//&models_new.CollectModel{},
		//&models_new.UserArticleCollectModel{},
		//&models_new.UserArticleLookHistoryModel{}, //用户浏览历史表
		//&models_new.CommentModel{},
		//&models_new.BannerModel{},
		//&models_new.LogModel{},
		//&models_new.UserLoginModel{},
		//&models_new.GlobalNotificationModel{},
		&models.BannerModel{},
		&models.TagModel{},
		&models.MessageModel{},
		&models.AdvertModel{},
		&models.UserModel{},
		&models.CommentModel{},
		//&models.ArticleModel{},  es上
		&models.MenuModel{},
		&models.MenuBannerModel{},
		&models.FeedbackModel{},
		&models.LoginDataModel{},
		&models.LogModel{},
		&models.ConversationModel{},
		&models.UserCollectModel{},
		&models.ChatModel{},
		&models.UserScopeModel{},
		&models.AutoReplyModel{},
		&models.BigModelRoleModel{},    //大模型角色表
		&models.BigModelTagModel{},     //大模型标签表
		&models.BigModelRoleTagModel{}, // 角色和大模型的 关联表
		&models.BigModelChatModel{},
		&models.BigModelSessionModel{},
	)
	if err != nil {
		logrus.Errorf("\n数据库迁移失败  %s", err)
		return
	}
	logrus.Info("\n数据库迁移成功")

}
