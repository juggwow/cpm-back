package connection

import (
	"time"

	"gorm.io/gorm"
)

type DBConnection struct {
	RAD *gorm.DB
}

type BaseModel struct {
	ID          uint       `gorm:"column:ID"`
	CreatedBy   string     `gorm:"column:CREATED_BY"`
	CreatedDate *time.Time `gorm:"column:CREATED_DATE"`
	UpdatedBy   string     `gorm:"column:UPDATED_BY"`
	UpdatedDate *time.Time `gorm:"column:UPDATED_DATE"`
}
