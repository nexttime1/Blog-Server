package enum

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
)

type RoleType int8

const (
	AdminRole   RoleType = 1
	UserRole    RoleType = 2
	VisitorRole RoleType = 3
	BlackRole   RoleType = 4
)

func (r RoleType) ParseRole() string {
	switch r {
	case AdminRole:
		return "管理员"
	case UserRole:
		return "普通用户"
	case VisitorRole:
		return "游客"
	case BlackRole:
		return "黑名单"
	}
	return "游客（未知）"
}

// Scan 实现 sql.Scanner 接口，用于从数据库读取数据时的类型转换
func (r *RoleType) Scan(value interface{}) error {
	if value == nil {
		*r = VisitorRole // 默认为游客角色
		return nil
	}

	// 尝试将数据库值转换为int64
	switch v := value.(type) {
	case int64:
		*r = RoleType(v)
	case int:
		*r = RoleType(v)
	case []byte:
		// 如果数据库存储的是字符串形式的数字
		val, err := strconv.ParseInt(string(v), 10, 8)
		if err != nil {
			return err
		}
		*r = RoleType(val)
	default:
		return errors.New("类型有错误")
	}

	return nil
}

// Value 实现 driver.Valuer 接口，用于将RoleType存入数据库时的转换
func (r RoleType) Value() (driver.Value, error) {
	return int64(r), nil
}

func (r RoleType) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(r)) // 直接返回数字（如1、2、3、4）
}

// MarshalJSON 实现 json.Marshaler 接口，序列化时返回中文名称
//     return json.Marshal(r.ParseRole()) // 返回的是"管理员"这类字符串

// UnmarshalJSON 实现 json.Unmarshaler 接口，反序列化时将字符串转换为对应的RoleType
func (r *RoleType) UnmarshalJSON(data []byte) error {
	// 1️⃣ 优先尝试解析为数字
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		*r = RoleType(num)
		return nil
	}

	// 2️⃣ 再尝试解析为字符串
	var roleStr string
	if err := json.Unmarshal(data, &roleStr); err != nil {
		return err
	}

	switch roleStr {
	case "1", "管理员":
		*r = AdminRole
	case "2", "普通用户":
		*r = UserRole
	case "3", "游客":
		*r = VisitorRole
	case "4", "黑名单":
		*r = BlackRole
	default:
		return errors.New("无效的角色类型: " + roleStr)
	}
	return nil
}
