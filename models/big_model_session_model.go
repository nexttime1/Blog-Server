package models

// BigModelSessionModel 大模型会话表（逻辑外键版本）
type BigModelSessionModel struct {
	Model
	UserID    uint                `json:"userID"`                                         // 用户id 逻辑外键
	UserModel UserModel           `gorm:"foreignKey:UserID;constraint:false" json:"-"`    // 禁用物理外键约束
	RoleID    uint                `json:"roleID"`                                         // 角色id 逻辑外键
	RoleModel BigModelRoleModel   `gorm:"foreignKey:RoleID;constraint:false" json:"-"`    // 禁用物理外键约束
	ChatList  []BigModelChatModel `gorm:"foreignKey:SessionID;constraint:false" json:"-"` // 会话列表（逻辑关联）
}
