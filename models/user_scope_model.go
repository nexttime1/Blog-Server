package models

type UserScopeModel struct {
	Model
	UserID uint `json:"userID"`
	Scope  int  `json:"scope"`
	Status bool `json:"status"`
}
