package models_new

type BannerModel struct {
	Model
	Path string `json:"path"`                // 图片路径
	Hash string `json:"hash"`                // 图片的hash值，用于判断重复图片
	Name string `gorm:"size:38" json:"name"` // 图片名称

}
