package data

import (
	"borges.ai/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"time"
)

func CreateJob(db *gorm.DB, user models.User, modelID uint64, modelType, jobType string) (models.Job, error) {
	record := models.Job{
		ModelID:   modelID,
		ModelType: modelType,
		Type:      jobType,
		StartDate: time.Now(),
	}

	record.CreatorID = user.ID
	record.UpdatedAt = time.Now()

	err := db.Create(&record).Error
	if err != nil {
		log.WithField("record", record).WithError(err).Error("failed to find or create job")
		return record, err
	}
	log.WithField("record", record.ID).Info("found or created job")
	return record, nil
}

func FindJobByID(db *gorm.DB, id uint64) (models.Job, error) {
	record := models.Job{}
	err := db.First(&record, id).Error
	log.WithField("record", record).Debug("found by id")
	return record, err
}

func UpdateJobAsCompleted(db *gorm.DB, user models.User, job models.Job) (models.User, error) {
	update := models.Job{
		FinishDate: time.Now(),
		Success:    true,
	}
	err := db.Model(&job).Debug().Where("creator_id=?", user.ID).Update(update).Error
	return user, err
}

func UpdateJobAsFailed(db *gorm.DB, user models.User, job models.Job, errorMsg string) (models.User, error) {
	update := models.Job{
		FinishDate: time.Now(),
		Success:    false,
		Error:      errorMsg,
	}
	err := db.Model(&job).Debug().Where("creator_id=?", user.ID).Update(update).Error
	return user, err
}
