package models

// BigModelTagModel 大模型标签表
type BigModelTagModel struct {
	Model
	Title string              `gorm:"size:16" json:"title"` // 标签的名称
	Color string              `gorm:"size:16" json:"color"` // 颜色
	Roles []BigModelRoleModel `gorm:"many2many:big_model_role_tag_models;foreignKey:ID;references:ID;constraint:false" json:"roles"`
}

/*
	Roles []BigModelRoleModel `gorm:"many2many:big_model_role_tag_models; // 必须与上面的中间表名一致
        foreignKey:ID; // 当前标签表的ID（作为逻辑外键，对应中间表的big_model_tag_model_id）
        references:ID; // 关联角色表的ID（对应中间表的big_model_role_model_id）
        constraint:false"` // 禁用物理外键约束

*/
