package data

import (
	"borges.ai/models"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"time"
)

func CreateCreateEvent(db *gorm.DB, user models.User, typeName, message string, payload map[string]string, isInitiatedByUser bool) (models.AuditEvent, error) {
	log.Info("new event")
	payloadJSON, _ := json.Marshal(payload)
	record := models.AuditEvent{
		Type:            "create",
		ModelType:       typeName,
		Message:         message,
		InitiatedByUser: isInitiatedByUser,
		Payload:         postgres.Jsonb{payloadJSON},
	}
	record.CreatorID = user.ID
	record.UpdatedAt = time.Now()
	err := db.Create(&record).Error
	return record, err
}
func CreateEvent(db *gorm.DB, user models.User, modelID uint64, typeName, eventType, attrName, attrValue, attrType, message string, payload map[string]string, initiatedByUser bool) (models.AuditEvent, error) {
	log.Info("new event")
	payloadJSON, _ := json.Marshal(payload)
	record := models.AuditEvent{
		Type:            eventType,
		ModelID:         modelID,
		ModelType:       typeName,
		Message:         message,
		InitiatedByUser: initiatedByUser,
		Payload:         postgres.Jsonb{payloadJSON},
		AttrName:        attrName,
		AttrValue:       attrValue,
		AttrType:        attrType,
	}
	record.CreatorID = user.ID
	record.UpdatedAt = time.Now()
	err := db.Create(&record).Error
	return record, err
}

func UpdateEvent(db *gorm.DB, user models.User, modelID uint64, typeName, eventType, attrName, attrValue, attrType, message string, payload map[string]string) (models.AuditEvent, error) {
	log.Info("update event")
	payloadJSON, _ := json.Marshal(payload)
	record := models.AuditEvent{
		Type:      eventType,
		ModelID:   modelID,
		ModelType: typeName,
		Message:   message,

		Payload: postgres.Jsonb{payloadJSON},

		AttrName:  attrName,
		AttrValue: attrValue,
		AttrType:  attrType,
	}
	record.CreatorID = user.ID
	record.UpdatedAt = time.Now()
	err := db.Create(&record).Error
	return record, err
}

func TransactEvent(db *gorm.DB, user models.User, id uint64) (models.AuditEvent, error) {
	log.Info("transact event")
	record := models.AuditEvent{}
	record.ID = id
	err := db.Model(&record).Update("transacted", true).Error
	return record, err
}

func TransactCreateEvent(db *gorm.DB, user models.User, id uint64, modelID uint64) (models.AuditEvent, error) {
	log.Info("transact create event")
	record := models.AuditEvent{}
	record.ID = id
	update := models.AuditEvent{
		Transacted: true,
		ModelID:    modelID,
	}
	err := db.Model(&record).Update(update).Error
	return record, err
}
