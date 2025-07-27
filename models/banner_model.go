package models

import (
	"Blog_server/models/enum"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
)

type BannerModel struct {
	Model
	Path      string         `json:"path"`                // 图片路径
	Hash      string         `json:"hash"`                // 图片的hash值，用于判断重复图片
	Name      string         `gorm:"size:38" json:"name"` // 图片名称
	ImageType enum.ImageType `gorm:"default:1" json:"image_type"`
}

// BeforeDelete 钩子函数  BeforeDelete 钩子函数会在使用 GORM 的 Delete 方法删除 BannerModel 记录时触发
func (b *BannerModel) BeforeDelete(tx *gorm.DB) (err error) {
	if b.ImageType == enum.LocationType {
		err := os.Remove(b.Path)
		if err != nil {
			logrus.Errorf("本地用os删除图片失败 %s", err)
			return err
		}
	}
	return nil
}
