package models

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

type IDRequest struct {
	ID string `json:"id" form:"id" uri:"id"` ////form Query  url /:id
}
type RemoveRequest struct {
	IDList []int `json:"IDList"`
}
type EsIdQuest struct {
	ID string `json:"id" form:"id" uri:"id"` //form Query  url /:id
}
