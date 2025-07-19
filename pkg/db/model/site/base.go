package site

import (
	"gorm.io/gorm"

	"serverApi/pkg/tools/tx"
)

func NewSite(db *gorm.DB) *Site {
	return &Site{
		DB: db,
	}
}

type Site struct {
	DB *gorm.DB
	tx tx.Tx
}
