package common

import (
	"Blog_server/global"
	"fmt"
	"gorm.io/gorm"
)

type PageInfo struct {
	Limit int    `form:"limit"`
	Page  int    `form:"page"`
	Key   string `form:"key"`
	Order string `form:"order"` //前端可以覆盖
}

type Options struct {
	PageInfo
	Likes        []string
	Preload      []string
	Where        *gorm.DB
	Debug        bool
	DefaultOrder string //内层   order
}

func (p PageInfo) GetLimit() int {
	if p.Limit <= 0 || p.Limit >= 50 {
		p.Limit = 10
	}
	return p.Limit
}
func (p PageInfo) GetPage() int {
	if p.Page <= 0 || p.Page >= 20 {
		return 1
	}
	return p.Page
}
func (p PageInfo) GetOffset() int {
	offset := (p.GetPage() - 1) * p.GetLimit()
	return offset

}

func ListQuery[T any](model T, options Options) (list []T, count int, err error) {
	//基础查询  GORM 会提取model中非零值字段作为等值查询条件  零值如0、""、nil会被忽略
	query := global.DB.Model(model).Where(model)

	//模糊查询
	if len(options.Likes) > 0 && options.PageInfo.Key != "" {
		//创建一个空的 查询条件构造器  用于后续动态添加 OR 条件
		likes := global.DB.Where("")
		for _, column := range options.Likes {
			likes.Or(fmt.Sprintf("%s like ?", column), fmt.Sprintf("%%%s%%", options.Key))
		}
		//合并查询条件
		query = query.Where(likes)
	}

	//日志
	if options.Debug {
		query = query.Debug()
	}
	//预加载
	if options.Preload != nil {
		for _, preload := range options.Preload {
			query = query.Preload(preload)
		}
	}
	//Order 排序   前端没传送，用默认的  前端传了 就用前端的
	if options.PageInfo.Order != "" {
		query = query.Order(options.PageInfo.Order)
	} else {
		query = query.Order(options.DefaultOrder)
	}

	// 定制化查询你
	if options.Where != nil {
		query = query.Where(options.Where)
	}

	//求总数

	var _count int64
	query.Count(&_count)
	count = int(_count)

	// 分页

	limit := options.PageInfo.GetLimit()
	offset := options.PageInfo.GetOffset()
	err = query.Limit(limit).Offset(offset).Find(&list).Error

	return
}
