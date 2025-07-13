package models

type UserArticleLookHistoryModel struct {
	Model
	UserID       uint         `json:"userID"`
	UserModel    UserModel    `gorm:"ForeignKey:UserID" json:"-"`
	ArticleID    uint         `json:"articleID"`
	ArticleModel ArticleModel `gorm:"ForeignKey:ArticleID" json:"-"`
}
