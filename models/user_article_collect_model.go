package models

import "time"

type UserArticleCollectModel struct {
	CreatedAt    time.Time    `json:"createdAt"` //收藏时间
	UserID       uint         `gorm:"uniqueIndex:idx_name" json:"userID"`
	UserModel    UserModel    `gorm:"ForeignKey:UserID" json:"-"`
	ArticleID    uint         `gorm:"uniqueIndex:idx_name" json:"articleID"`
	ArticleModel ArticleModel `gorm:"ForeignKey:ArticleID" json:"-"`
	CollectID    uint         `gorm:"uniqueIndex:idx_name" json:"collectID"` //收藏夹的ID
	CollectModel CollectModel `gorm:"ForeignKey:CollectID" json:"-"`
}
