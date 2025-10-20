package models

// BigModelChatModel 大模型对话表（逻辑外键版本）
type BigModelChatModel struct {
	Model
	SessionID    uint                 `json:"sessionID"`                                      // 会话id 逻辑外键
	SessionModel BigModelSessionModel `gorm:"foreignKey:SessionID;constraint:false" json:"-"` // 禁用物理外键约束
	Status       bool                 `json:"status"`                                         // 状态，ai有没有正常的回复用户
	Content      string               `json:"content"`                                        // 用户的聊天内容
	BotContent   string               `json:"botContent"`                                     // ai的回复内容
	RoleID       uint                 `json:"roleID"`                                         // 角色id 逻辑外键，
	RoleModel    BigModelRoleModel    `gorm:"foreignKey:RoleID;constraint:false" json:"-"`    // 禁用物理外键约束
	UserID       uint                 `json:"userID"`                                         // 用户id 逻辑外键
	UserModel    UserModel            `gorm:"foreignKey:UserID;constraint:false" json:"-"`    // 禁用物理外键约束
}
