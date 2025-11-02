package common

import (
	"Blog_server/global"
	"Blog_server/models"
	"github.com/sirupsen/logrus"
)

func BatchRemove[T any](ModelList []T, rc models.RemoveRequest) []T {
	global.DB.Where("id in ?", rc.IDList).Find(&ModelList)
	logrus.Infof("ModelList 里面有 %d 条", len(ModelList))
	if len(ModelList) > 0 {
		err := global.DB.Delete(&ModelList).Error
		if err != nil {
			logrus.Error("BatchRemove ", err)
		}
	}

	return ModelList
}
