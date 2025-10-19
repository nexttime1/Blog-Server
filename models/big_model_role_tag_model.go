package models

type BigModelRoleTagModel struct {
	// 联合主键
	BigModelTagModelId  uint `gorm:"primaryKey" json:"big_model_tag_model_id"`  // 关联标签表的ID
	BigModelRoleModelId uint `gorm:"primaryKey" json:"big_model_role_model_id"` // 关联角色表的ID
}
