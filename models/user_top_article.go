package models

import "time"

type UserTopArticle struct {
	CreatedAt    time.Time    `json:"createdAt"` //  创作置顶的时间
	UserID       uint         `gorm:"uniqueIndex:idx_name" json:"userID"`
	UserModel    UserModel    `gorm:"ForeignKey:UserID" json:"-"`
	ArticleID    uint         `gorm:"uniqueIndex:idx_name" json:"articleID"`
	ArticleModel ArticleModel `gorm:"ForeignKey:ArticleID" json:"-"`
}
