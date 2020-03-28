package models

import "github.com/jinzhu/gorm/dialects/postgres"

type AuditEvent struct {
	Model
	Type      string `gorm:"type:varchar(20);not null"` // review, comment, like, status, rate
	ModelType string `gorm:"type:varchar(20)not null"`
	ModelID   uint64 `gorm:"not null"`

	Message         string `gorm:"type:text"`
	AttrName        string `gorm:"index:attr_name"`
	AttrValue       string `gorm:"index:attr_value"`
	AttrType        string `gorm:"index:attr_type"`
	Payload         postgres.Jsonb
	Transacted      bool
	InitiatedByUser bool
}
