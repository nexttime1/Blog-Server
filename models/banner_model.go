package models

type BannerModel struct {
	Model
	Cover string `gorm:"size:256" json:"cover"` //图片链接
	Href  string `gorm:"size:256" json:"href"`  //跳转链接

}
