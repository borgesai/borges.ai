package services

import (
	"borges.ai/data"
	"borges.ai/jobs"
	"borges.ai/models"
	"borges.ai/utils"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func ImportReview(db *gorm.DB, user models.User, bookID, editionID uint64, goodreadsID, reviewContent string, rating, status int, statusDate time.Time, isInitiatedByUser bool) (models.Review, error) {
	defer utils.Duration(utils.Track("ImportReview"))
	eventPayload := map[string]string{
		"bookID":      strconv.FormatUint(bookID, 10),
		"editionID":   strconv.FormatUint(editionID, 10),
		"goodreadsID": goodreadsID,
		"content":     reviewContent,
		"rating":      strconv.Itoa(rating),
		"status":      strconv.Itoa(status),
		"status_date": utils.FormatDate(statusDate),
	}
	event, err := data.CreateCreateEvent(db, user, "review", "create/update review", eventPayload, isInitiatedByUser)
	if err != nil {
		return models.Review{}, err
	}
	review, err := data.FindOrCreateFullBookReviewFromSyncProcedure(db, user, bookID, editionID, goodreadsID, reviewContent, rating, status, statusDate)
	_, err = data.TransactCreateEvent(db, user, event.ID, review.ID)

	return review, err
}

func UpdateReview(db *gorm.DB, user models.User, bookID, editionID uint64, reviewContent string) (models.Review, error) {
	defer utils.Duration(utils.Track("UpdateReview"))
	eventPayload := map[string]string{
		"bookID":    strconv.FormatUint(bookID, 10),
		"editionID": strconv.FormatUint(editionID, 10),
		"content":   reviewContent,
	}
	event, err := data.CreateEvent(db, user, 0, "review", "upsert",
		"content", reviewContent, "string",
		"create/update review", eventPayload, true)
	if err != nil {
		return models.Review{}, err
	}
	review, err := data.FindOrCreateBookReview(db, user, bookID, editionID, reviewContent)
	if err != nil {
		return review, err
	}
	err = jobs.SyncUserReviewToExternalSystems(user, review.ID, bookID, editionID)
	if err != nil {
		log.WithError(err).Error("failed to sync data")
	}
	_, err = data.TransactCreateEvent(db, user, event.ID, review.ID)

	return review, err
}

func UpdateUserBookStatus(db *gorm.DB, user models.User, bookID, editionID uint64, status int) (models.Review, error) {
	defer utils.Duration(utils.Track("UpdateUserBookStatus"))
	eventPayload := map[string]string{
		"bookID":    strconv.FormatUint(bookID, 10),
		"editionID": strconv.FormatUint(editionID, 10),
		"status":    strconv.Itoa(status),
	}
	event, err := data.CreateEvent(db, user, 0, "review", "upsert",
		"status", strconv.Itoa(status), "int",
		"create/update status", eventPayload, true)
	if err != nil {
		return models.Review{}, err
	}
	review, err := data.FindOrCreateBookStatus(db, user, bookID, editionID, status, time.Now())
	if err != nil {
		return review, err
	}
	err = jobs.SyncUserReviewToExternalSystems(user, review.ID, bookID, editionID)
	if err != nil {
		log.WithError(err).Error("failed to sync data")
	}
	_, err = data.TransactCreateEvent(db, user, event.ID, review.ID)
	return review, err
}

func UpdateUserBookRating(db *gorm.DB, user models.User, bookID, editionID uint64, rating int) (models.Review, error) {
	defer utils.Duration(utils.Track("UpdateUserBookRating"))
	eventPayload := map[string]string{
		"bookID":    strconv.FormatUint(bookID, 10),
		"editionID": strconv.FormatUint(editionID, 10),
		"rating":    strconv.Itoa(rating),
	}
	event, err := data.CreateEvent(db, user, 0, "review", "upsert",
		"rating", strconv.Itoa(rating), "int",
		"create/update rating", eventPayload, true)
	if err != nil {
		return models.Review{}, err
	}
	review, err := data.FindOrCreateBookRating(db, user, bookID, editionID, rating)
	if err != nil {
		return review, err
	}
	err = jobs.SyncUserReviewToExternalSystems(user, review.ID, bookID, editionID)
	if err != nil {
		log.WithError(err).Error("failed to sync data")
	}
	_, err = data.TransactCreateEvent(db, user, event.ID, review.ID)
	return review, err
}


func UpdateReviewAndReadingsWithGoodreadsReviewID (db *gorm.DB, user models.User, goodreadsReviewID string, review models.Review, readings []models.Reading) error {
	defer utils.Duration(utils.Track("UpdateReviewAndReadingsWithGoodreadsReviewID"))
	_, err := data.UpdateReviewWithGoodreadsID(db, user, review, goodreadsReviewID)
	if err != nil {
		log.WithError(err).WithField("goodreads_id", goodreadsReviewID).Error("failed to update review with goodreads is")
		// we are going to retry this one
		return err
	}
	_, err = data.UpdateReadingsWithGoodreadsID(db, user, readings, goodreadsReviewID)
	if err != nil {
		log.WithError(err).WithField("goodreads_id", goodreadsReviewID).Error("failed to update readings with goodreads is")
		// we are going to retry this one
		return err
	}
	return nil
}

func FindUserReviewOptional(db *gorm.DB, user models.User, bookID uint64)(models.Review, error){
	review, err := data.FindReviewByBookID(db, user, bookID)
	if err != nil {
		if err.Error() == "record not found" {
			return models.Review{}, nil
		}
		return models.Review{}, err
	}
	return review, nil
}