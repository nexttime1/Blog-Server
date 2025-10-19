package models

// BigModelRoleModel 大模型角色表
type BigModelRoleModel struct {
	Model
	Name     string `gorm:"size:16" json:"name"`  // 角色名称
	Enable   bool   `json:"enable"`               // 是否启用
	Icon     string `gorm:"size:256" json:"icon"` // 可以选择系统默认的一些，也可以图片上传
	Abstract string `gorm:"size:256" json:"abstract"`
	// 角色简介
	Tags      []BigModelTagModel `gorm:"many2many:big_model_role_tag_models;foreignKey:ID;references:ID;constraint:false" json:"tags"`
	Scope     int                `json:"scope"`                    // 消耗的积分
	Prologue  string             `gorm:"size:512" json:"prologue"` // 开场白
	Prompt    string             `gorm:"size:2048" json:"prompt"`  // 设定词
	AutoReply bool               `json:"autoReply"`                // 是否接入自动回复
}

/*
	Tags []BigModelTagModel `gorm:"many2many:big_model_role_tag_models; // 指定中间表（结构体名对应的表名，GORM自动转蛇形）
        foreignKey:ID; // 当前角色表的ID（作为逻辑外键，对应中间表的big_model_role_model_id）
        references:ID; // 关联标签表的ID（对应中间表的big_model_tag_model_id）
        constraint:false"` // 禁用物理外键约束

*/
