package db

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DDGQueryTable struct {
	gorm.Model
	Query  string         `gorm:"primaryKey" json:"query"`
	Answer datatypes.JSON `json:"answer" binding:"required"`
}
