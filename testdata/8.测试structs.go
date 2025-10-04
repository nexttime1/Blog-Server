package main

import (
	"fmt"
	"github.com/fatih/structs"
)

type AdvertRequest struct { //json
	Title  string `json:"title" binding:"required" structs:"title"`   // 显示的标题
	Href   string `json:"href" binding:"required,url" structs:"href"` // 跳转链接
	Images string `json:"images" binding:"required,uri" structs:"-"`  // 图片
	IsShow bool   `json:"is_show" structs:"is_show"`                  // 是否展示
}

func main() {

	m := structs.Map(AdvertRequest{
		Href:   "xxx",
		Images: "xxx",
		IsShow: true,
	})
	for k, v := range m {
		fmt.Println(k, v)
		if m[k] == "" {
			delete(m, k)
			fmt.Println(k)
		}
	}
	fmt.Println(m)

}
