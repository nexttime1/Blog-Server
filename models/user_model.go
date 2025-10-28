package models

import (
	"Blog_server/models/enum"
	"time"
)

type UserModel struct {
	Model
	Username       string            `gorm:"size:32" json:"username"`
	Nickname       string            `gorm:"size:32" json:"nickname,select(c|info)"`
	Password       string            `gorm:"size:64" json:"-"`
	Avatar         string            `gorm:"size:256" json:"avatar,select(c)"` //头像
	Abstract       string            `gorm:"size:256" json:"abstract"`
	RegisterSource enum.RegisterType `json:"registerSource"`
	CodeAge        string            `json:"codeAge"`
	Email          string            `gorm:"size:256" json:"email"`
	Tel            string            `gorm:"size:18" json:"tel"`
	Addr           string            `gorm:"size:64" json:"addr,select(c|info)"`
	Token          string            `gorm:"size:64" json:"token"` //第三方登陆的唯一id
	IP             string            `gorm:"size:32" json:"ip,select(c)"`
	Scope          int               `gorm:"default:0" json:"scope,select(info)"`
	Role           enum.RoleType     `json:"role"` // 1 为管理员  2 为 普通用户

}

type UserConfModel struct {
	UserID             uint       `gorm:"unique" json:"userID"`
	UserModel          UserModel  `gorm:"foreignKey:UserID" json:"-"`
	LikeTag            []string   `gorm:"type:longtext;serializer:json" json:"likeTag"`
	UpdateUserNameDate *time.Time `json:"updateUserNameDate"` //上传修改用户名的时间
	OpenCollect        bool       `json:"openCollect"`        //公开我的收藏
	OpenFollow         bool       `json:"openFollow"`         //公开我的关注
	OpenFans           bool       `json:"openFans"`           //公开我的粉丝
	HomeStyleID        uint       `json:"homeStyleID"`        //主页用户风格id
}
