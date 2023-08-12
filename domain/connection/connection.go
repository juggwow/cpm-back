package connection

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type DBConnection struct {
	CPM *gorm.DB
}

// type BaseModel struct {
// 	ID          uint       `gorm:"column:ID"`
// 	CreatedBy   string     `gorm:"column:CREATED_BY"`
// 	CreatedDate *time.Time `gorm:"column:CREATED_DATE"`
// 	UpdatedBy   string     `gorm:"column:UPDATED_BY"`
// 	UpdatedDate *time.Time `gorm:"column:UPDATED_DATE"`
// }

type BaseModel struct {
	ID          uint                  `gorm:"column:ID;primarykey"`
	CreatedBy   string                `gorm:"column:CREATED_BY"`
	CreatedDate *time.Time            `gorm:"column:CREATED_DATE;autoCreateTime"`
	UpdatedBy   string                `gorm:"column:UPDATED_BY"`
	UpdatedDate *time.Time            `gorm:"column:UPDATED_DATE"`
	DeleteFlag  soft_delete.DeletedAt `gorm:"column:DELETE_FLAG;type:tinyint;softDelete:flag;default:0"`
}
