package services

import (
	"borges.ai/data"
	"borges.ai/models"
	"borges.ai/utils"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

func ImportReading(db *gorm.DB, user models.User, bookID, editionID uint64, goodreadsID string, startDateStr, finishDateStr string, startDate, finishDate time.Time) (models.Reading, error) {
	defer utils.Duration(utils.Track("ImportReading"))
	eventPayload := map[string]string{
		"bookID":      strconv.FormatUint(bookID, 10),
		"editionID":   strconv.FormatUint(editionID, 10),
		"goodreadsID": goodreadsID,
		"start_date":  startDateStr,
		"finish_date": finishDateStr,
	}
	event, err := data.CreateCreateEvent(db, user, "reading", "create reading", eventPayload, false)
	if err != nil {
		return models.Reading{}, err
	}
	reading, err := data.CreateBookReading(db, user, bookID, editionID, goodreadsID, startDateStr, finishDateStr, startDate, finishDate)
	if err != nil {
		return reading, err
	}
	_, err = data.TransactCreateEvent(db, user, event.ID, reading.ID)

	return reading, err
}
func UpdateReadingNote(db *gorm.DB, user models.User, bookID, editionID uint64, note string) (models.Reading, error) {
	defer utils.Duration(utils.Track("UpdateReadingNote"))
	eventPayload := map[string]string{
		"bookID":    strconv.FormatUint(bookID, 10),
		"editionID": strconv.FormatUint(editionID, 10),
		"note":      note,
	}
	event, err := data.CreateEvent(db, user, 0, "reading", "upsert",
		"note", note, "string",
		"create/update reading start date", eventPayload, true)
	if err != nil {
		return models.Reading{}, err
	}
	reading, err := data.FindOrCreateBookReadingWithNote(db, user, bookID, editionID, note)
	if err != nil {
		return reading, err
	}
	_, err = data.TransactCreateEvent(db, user, event.ID, reading.ID)
	return reading, err
}
func UpdateReadingStartDate(db *gorm.DB, user models.User, bookID, editionID uint64, startDateStr string, startDate time.Time) (models.Reading, error) {
	defer utils.Duration(utils.Track("UpdateReadingStartDate"))
	eventPayload := map[string]string{
		"bookID":    strconv.FormatUint(bookID, 10),
		"editionID": strconv.FormatUint(editionID, 10),
		"startDate": startDateStr,
	}
	event, err := data.CreateEvent(db, user, 0, "reading", "upsert",
		"start_date", startDateStr, "time",
		"create/update reading start date", eventPayload, true)
	if err != nil {
		return models.Reading{}, err
	}
	reading, err := data.FindOrCreateBookReading(db, user, bookID, editionID, startDateStr, "", startDate, time.Time{})
	if err != nil {
		return reading, err
	}
	_, err = data.TransactCreateEvent(db, user, event.ID, reading.ID)
	return reading, err
}

func UpdateReadingFinishDate(db *gorm.DB, user models.User, bookID, editionID uint64, finishDateStr string, finishDate time.Time) (models.Reading, error) {
	defer utils.Duration(utils.Track("UpdateReadingFinishDate"))
	eventPayload := map[string]string{
		"bookID":     strconv.FormatUint(bookID, 10),
		"editionID":  strconv.FormatUint(editionID, 10),
		"finishDate": finishDateStr,
	}
	event, err := data.CreateEvent(db, user, 0, "reading", "upsert",
		"finish_date", finishDateStr, "time",
		"create/update reading finish date", eventPayload, true)
	if err != nil {
		return models.Reading{}, err
	}
	reading, err := data.FindOrCreateBookReading(db, user, bookID, editionID, "", finishDateStr, time.Time{}, finishDate)
	if err != nil {
		return reading, err
	}
	_, err = data.TransactCreateEvent(db, user, event.ID, reading.ID)
	if err != nil {
		return reading, err
	}
	// the book is finished. Change status and update UI too
	_, err = UpdateUserBookStatus(db, user, bookID, editionID, data.READ_STATUS)
	return reading, err
}
