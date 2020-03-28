package data

import (
	"borges.ai/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"time"
)

func CreateBookReading(db *gorm.DB, user models.User, bookID, bookEditionID uint64, goodreadsID string, startDateStr, finishDateStr string, startDate, finishDate time.Time) (models.Reading, error) {
	record := models.Reading{}
	if startDateStr != "" {
		record.StartDateRaw = startDateStr
		if !startDate.IsZero() {
			record.StartDate = startDate
		}
	}
	if finishDateStr != "" {
		record.FinishDateRaw = finishDateStr
		if !finishDate.IsZero() {
			record.FinishDate = finishDate
		}
	}

	record.BookID = bookID

	if bookEditionID > 0 {
		record.BookEditionID = bookEditionID
	}

	if goodreadsID != "" {
		record.GoodreadsID = goodreadsID
	}

	record.CreatorID = user.ID
	record.UpdatedAt = time.Now()

	err := db.Create(&record).Error
	if err != nil {
		log.WithField("record", record).WithError(err).Error("failed to find or create a reading")
		return record, err
	}
	log.WithField("record", record.ID).Info("found or created reading")
	return record, nil
}

func FindOrCreateBookReading(db *gorm.DB, user models.User, bookID, bookEditionID uint64, startDateStr, finishDateStr string, startDate, finishDate time.Time) (models.Reading, error) {
	record := models.Reading{}
	if startDateStr != "" {
		record.StartDateRaw = startDateStr
	}
	if finishDateStr != "" {
		record.FinishDateRaw = finishDateStr
	}
	if !startDate.IsZero() {
		record.StartDate = startDate
	}
	if !finishDate.IsZero() {
		record.FinishDate = finishDate
	}

	record.BookID = bookID

	if bookEditionID > 0 {
		record.BookEditionID = bookEditionID
	}

	record.CreatorID = user.ID
	record.UpdatedAt = time.Now()

	lookupFields := models.Reading{BookID: bookID}
	lookupFields.CreatorID = user.ID

	err := db.Debug().Where(lookupFields).Assign(record).FirstOrCreate(&record).Error
	if err != nil {
		log.WithField("record", record).WithError(err).Error("failed to find or create a reading")
		return record, err
	}
	log.WithField("record", record.ID).Info("found or created reading")
	return record, nil
}

func FindOrCreateBookReadingWithNote(db *gorm.DB, user models.User, bookID, bookEditionID uint64, note string) (models.Reading, error) {
	record := models.Reading{}

	record.BookID = bookID

	if bookEditionID > 0 {
		record.BookEditionID = bookEditionID
	}

	record.Note = note

	record.CreatorID = user.ID
	record.UpdatedAt = time.Now()

	lookupFields := models.Reading{BookID: bookID}
	lookupFields.CreatorID = user.ID

	err := db.Debug().Where(lookupFields).Assign(record).FirstOrCreate(&record).Error
	if err != nil {
		log.WithField("record", record).WithError(err).Error("failed to find or create a reading")
		return record, err
	}
	log.WithField("record", record.ID).Info("found or created reading")
	return record, nil
}

func FindUserReadings(db *gorm.DB, user models.User) ([]models.Reading, error) {
	records := []models.Reading{}
	err := db.Where("creator_id=? AND finish_date > '1800-01-01' AND start_date is not null", user.ID).Find(&records).Error
	return records, err
}

func FindUserReadingsMapByBookID(db *gorm.DB, user models.User) (map[uint64]models.Reading, error) {
	mapByBookID := make(map[uint64]models.Reading, 0)
	records, err := FindUserReadings(db, user)
	if err != nil {
		return mapByBookID, err
	}

	// TODO for now we just override things. wouldn't work for multiple readings of the same book
	for _, record := range records {
		mapByBookID[record.BookID] = record
	}
	return mapByBookID, nil
}


func FindUserBookReadings(db *gorm.DB, user models.User, bookID uint64) ([]models.Reading, error) {
	records := []models.Reading{}
	err := db.Where("creator_id=? AND book_id=?", user.ID, bookID).Find(&records).Error
	return records, err
}


func UpdateReadingsWithGoodreadsID(db *gorm.DB, user models.User, readings []models.Reading, goodreadsID string) ([]models.Reading, error) {
	updates := map[string]interface{}{"goodreads_id": goodreadsID}
	ids:=make([]uint64,0)
	for _, reading := range readings{
		ids = append(ids, reading.ID)
	}
	err := db.Table("readings").Debug().Where("id IN (?) AND creator_id=?", ids,user.ID).Updates(updates).Error
	return readings, err
}
