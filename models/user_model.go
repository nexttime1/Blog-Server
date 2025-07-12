package models

type UserModel struct {
	Model
	Username       string   `gorm:"size:32" json:"username"`
	Nickname       string   `gorm:"size:32" json:"nickname"`
	Password       string   `gorm:"size:64" json:"-"`
	Avatar         string   `gorm:"size:256" json:"avatar"` //头像
	Abstract       string   `gorm:"size:256" json:"abstract"`
	RegisterSource string   `json:"registerSource"`
	CodeAge        string   `json:"codeAge"`
	LikeTag        []string `gorm:"type:longtext;serializer:json" json:"likeTag"`
	Email          string   `gorm:"size:256" json:"email"`
	OpenID         string   `gorm:"size:64 " json:"openID"` //第三方登陆的唯一id
}
