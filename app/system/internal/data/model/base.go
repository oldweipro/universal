package model

import (
	"time"
	"universal/pkg/idgen"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含公共字段和ID生成逻辑
type BaseModel struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate GORM Hook，在创建记录前生成雪花ID
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == 0 {
		node, err := idgen.NewNode(1) // 节点ID可以通过配置管理
		if err != nil {
			return err
		}
		b.ID = node.Generate().Int64()
	}
	return nil
}
