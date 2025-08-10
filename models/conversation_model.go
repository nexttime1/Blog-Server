package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ConversationModel 会话模型：表示一个聊天会话（一对一）
type ConversationModel struct {
	Model
	UserIDs   UserIDSlice `gorm:"type:varchar(50)" json:"user_ids"` // 参与会话的用户ID列表（固定2个）
	LastMsgID uint        `json:"last_msg_id"`                      // 最后一条消息ID，用于快速展示
}

// UserIDSlice 自定义类型，用于实现sql.Scanner和driver.Valuer接口
type UserIDSlice []uint

// Value 实现driver.Valuer接口，将[]uint转换为数据库存储格式（逗号分隔字符串）
func (u UserIDSlice) Value() (driver.Value, error) {
	if len(u) == 0 {
		return "", nil
	}
	// 将uint切片转换为字符串切片
	strs := make([]string, len(u))
	for i, id := range u {
		strs[i] = strconv.FormatUint(uint64(id), 10)
	}
	// 用逗号连接成字符串
	return strings.Join(strs, ","), nil
}

// Scan 实现sql.Scanner接口，将数据库字符串转换为[]uint
func (u *UserIDSlice) Scan(value interface{}) error {
	if value == nil {
		*u = nil
		return nil
	}

	// 尝试将value转换为字符串
	str, ok := value.(string)
	if !ok {
		return errors.New("invalid type for UserIDSlice")
	}

	if str == "" {
		*u = nil
		return nil
	}

	// 按逗号分割字符串
	parts := strings.Split(str, ",")
	ids := make(UserIDSlice, 0, len(parts))
	for _, part := range parts {
		// 转换为uint
		id, err := strconv.ParseUint(part, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid user id: %s, error: %w", part, err)
		}
		ids = append(ids, uint(id))
	}

	*u = ids
	return nil
}
