package enum

import (
	"database/sql/driver"
	"strings"
)

type Array []string

// Scan 当从数据库读取数据时，数据库驱动会调用这个方法  将数据库查询结果转换为Array类型
func (t *Array) Scan(value interface{}) error {
	v, _ := value.([]byte)
	if string(v) == "" {
		*t = []string{}
		return nil
	}
	*t = strings.Split(string(v), "\n")
	return nil
}

// Value 实现了 Go 数据库 /sql 标准库中的 driver.Valuer 接口 自动调用
func (t Array) Value() (driver.Value, error) {
	// 将数字转换为值
	return strings.Join(t, "\n"), nil
}
