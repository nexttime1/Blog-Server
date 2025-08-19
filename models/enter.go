package models

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primarykey" json:"id,select($any)" structs:"-"` // 主键ID
	CreatedAt time.Time `json:"created_at,select($any)" structs:"-"`           // 创建时间
	UpdatedAt time.Time `json:"UpdatedAt" structs:"-"`
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
type EsIdListQuest struct {
	IDList []string `json:"id_list" binding:"required"`
}
