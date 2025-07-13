package models

import "fmt"

type ImageModel struct {
	Model
	Filename string `gorm:"size:64" json:"filename"`
	Path     string `gorm:"size:256" json:"path"`
	Size     int64  `json:"size"`
	Hash     string `gorm:"size:32" json:"hash"`
}

func (i ImageModel) WabPath() string {
	return fmt.Sprintf("%s%s", i.Path, i.Filename)
}
