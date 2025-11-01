package models

type ArticleUserIDModel struct {
	Model
	UserID    uint   `json:"user_id"`
	ArticleID string `json:"article_id"`
}
