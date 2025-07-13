package models

import "time"

// ArticleDiggModel 文章点赞表  只能点赞一次
type ArticleDiggModel struct {
	CreatedAt    time.Time    `json:"createdAt"`
	UserID       uint         `gorm:"uniqueIndex:idx_name" json:"userID"`
	UserModel    UserModel    `gorm:"ForeignKey:UserID" json:"-"`
	ArticleID    uint         `gorm:"uniqueIndex:idx_name" json:"articleID"`
	ArticleModel ArticleModel `gorm:"ForeignKey:ArticleID" json:"-"`
}
