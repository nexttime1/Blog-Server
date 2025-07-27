package common

import (
	"Blog_server/global"
	"Blog_server/models"
)

func BatchRemove[T any](ModelList []T, rc models.RemoveRequest) []T {
	global.DB.Find(&ModelList, "id in ?", rc.IDList)
	if len(ModelList) > 0 {
		global.DB.Delete(&ModelList)
	}
	return ModelList
}
