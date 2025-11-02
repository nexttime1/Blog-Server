package enum

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type LevelType int8

const (
	LogInfoLevel LevelType = 1
	LogWainLevel LevelType = 2
	LogErrLevel  LevelType = 3
)

// 实现fmt.Stringer 接口。  直接%s
func (level LevelType) String() string {
	switch level {
	case LogInfoLevel:
		return "info"
	case LogWainLevel:
		return "warning"
	case LogErrLevel:
		return "error"
	}
	return ""
}

// MarshalJSON 实现 json.Marshaler 接口（JSON 序列化时调用，返回字符串）
func (level LevelType) MarshalJSON() ([]byte, error) {
	// 将字符串用双引号包裹，符合 JSON 字符串格式
	return json.Marshal(level.String())
}

// UnmarshalJSON 实现 json.Unmarshaler 接口（JSON 反序列化时调用，将字符串转为 LevelType）
func (level *LevelType) UnmarshalJSON(data []byte) error {
	var s string
	// 先将 JSON 数据解析为字符串
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("LevelType 反序列化失败：%v", err)
	}
	// 根据字符串匹配对应的枚举值
	switch s {
	case "info":
		*level = LogInfoLevel
	case "warning":
		*level = LogWainLevel
	case "error":
		*level = LogErrLevel
	default:
		return fmt.Errorf("无效的 LevelType 字符串：%s", s)
	}
	return nil
}

// Value 实现driver.Valuer接口（写入数据库时自动调用，返回字符串）
func (level LevelType) Value() (driver.Value, error) {
	return level.String(), nil // 直接复用String()方法的结果
}

// Scan 实现sql.Scanner接口（从数据库读取时自动调用，将字符串转回LevelType）
func (level *LevelType) Scan(value interface{}) error {
	// 先将数据库中的值转换为字符串（可能是string或[]byte类型）
	var s string
	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v) // 数据库驱动可能返回字节切片，转换为字符串
	default:
		return fmt.Errorf("LevelType不支持的数据库类型：%T", value)
	}

	// 根据字符串匹配对应的LevelType
	switch s {
	case "info":
		*level = LogInfoLevel
	case "warning":
		*level = LogWainLevel
	case "error":
		*level = LogErrLevel
	default:
		return fmt.Errorf("无效的LevelType字符串：%s", s)
	}
	return nil
}
