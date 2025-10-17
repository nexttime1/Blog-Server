package models

type FeedbackModel struct {
	Model
	Email   string `gorm:"size:64" json:"email"`
	Content string `gorm:"size:256" json:"content"`
}
