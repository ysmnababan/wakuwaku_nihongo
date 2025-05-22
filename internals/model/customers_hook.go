package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (m *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now().UnixMilli()
	if m.CustomerID == "" {
		m.CustomerID = uuid.NewString()
	}

	return
}

func (m *Customer) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now().UnixMilli()
	m.ModifiedAt = &now
	return
}
