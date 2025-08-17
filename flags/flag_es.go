package flags

import (
	"Blog_server/models"
	"github.com/sirupsen/logrus"
)

func FlagES() {
	//err := models.ArticleModel{}.CreateIndex()
	//if err != nil {
	//	logrus.Error(err)
	//}
	err := models.FullTextModel{}.CreateIndex()
	if err != nil {
		logrus.Error(err)
	}
}
