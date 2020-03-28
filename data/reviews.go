package data

import (
	"borges.ai/models"
	"borges.ai/text"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	READ_STATUS           = 1
	STARTED_STATUS        = 2
	WANT_TO_READ_STATUS   = 3
	WANT_TO_REREAD_STATUS = 4
	WONT_FINISH_STATUS    = 5
)

func FindOrCreateBookRating(db *gorm.DB, user models.User, bookID, editionID uint64, rating int) (models.Review, error) {
	record := models.Review{}
	record.BookID = bookID
	record.BookEditionID = editionID
	record.Rating = rating

	record.CreatorID = user.ID
	record.UpdatedAt = time.Now()

	lookupFields := models.Review{BookID: bookID}
	lookupFields.CreatorID = user.ID

	err := db.Where(lookupFields).Assign(record).FirstOrCreate(&record).Error
	if err != nil {
		log.WithField("record", record).WithError(err).Error("failed to find or create a review")
		return record, err
	}
	log.WithField("record", record.ID).Info("found or created review")
	return record, nil
}
func FindOrCreateBookStatus(db *gorm.DB, user models.User, bookID, editionID uint64, status int, date time.Time) (models.Review, error) {
	record := models.Review{}
	record.BookID = bookID
	record.BookEditionID = editionID
	record.Status = status
	record.StatusDate = date

	record.CreatorID = user.ID
	record.UpdatedAt = time.Now()

	lookupFields := models.Review{BookID: bookID}
	lookupFields.CreatorID = user.ID

	err := db.Where(lookupFields).Assign(record).FirstOrCreate(&record).Error
	if err != nil {
		log.WithField("record", record).WithError(err).Error("failed to find or create a review")
		return record, err
	}
	log.WithField("record", record.ID).Info("found or created review")
	return record, nil
}

func FindOrCreateFullBookReviewFromSyncProcedure(db *gorm.DB, user models.User, bookID, editionID uint64, goodreadsID, reviewContent string, rating, status int, statusDate time.Time) (models.Review, error) {
	lookupFields := models.Review{BookID: bookID}
	lookupFields.CreatorID = user.ID

	existingRecord := models.Review{}
	err := db.Where(lookupFields).First(&existingRecord).Error
	if err != nil {
		// ignore this error
		if err.Error() != "record not found" {
			log.WithField("record", existingRecord).WithError(err).Error("failed to find or create a review")
			return existingRecord, err
		}
	}

	record := models.Review{}

	//create or content is empty for some reason
	if existingRecord.ID == 0 || existingRecord.Content == "" {
		// first import. that means we need to convert html into markdown
		if reviewContent != "" {
			reviewContent = text.ConvertHTMLToMarkdown(reviewContent)
		}
		if reviewContent != "" {
			record.Content = reviewContent
			record.ContentHash = text.CreateContentHash(reviewContent)
			record.ContentHTML = text.CreateContentHTML(reviewContent)
		}
	}
	record.BookID = bookID
	record.BookEditionID = editionID

	if goodreadsID != "" {
		record.GoodreadsID = goodreadsID
	}
	if rating > 0 {
		record.Rating = rating
	}

	if status > 0 {
		record.Status = status
		record.StatusDate = statusDate
	}

	record.CreatorID = user.ID
	record.UpdatedAt = time.Now()

	err = db.Where(lookupFields).Assign(record).FirstOrCreate(&record).Error
	if err != nil {
		log.WithField("record", record).WithError(err).Error("failed to find or create a review")
		return record, err
	}
	log.WithField("record", record.ID).Info("found or created review")
	return record, nil
}

func FindOrCreateBookReview(db *gorm.DB, user models.User, bookID, editionID uint64, reviewContent string) (models.Review, error) {
	record := models.Review{}
	if reviewContent != "" {
		record.Content = reviewContent
		record.ContentHash = text.CreateContentHash(reviewContent)
		record.ContentHTML = text.CreateContentHTML(reviewContent)
	}

	record.BookID = bookID

	record.BookEditionID = editionID

	record.CreatorID = user.ID
	record.UpdatedAt = time.Now()

	lookupFields := models.Review{BookID: bookID}
	lookupFields.CreatorID = user.ID

	err := db.Where(lookupFields).Assign(record).FirstOrCreate(&record).Error
	if err != nil {
		log.WithField("record", record).WithError(err).Error("failed to find or create a review")
		return record, err
	}
	log.WithField("record", record.ID).Info("found or created review")
	return record, nil
}

func FindReviewByBookID(db *gorm.DB, user models.User, bookID uint64) (models.Review, error) {
	record := models.Review{}
	err := db.Where("book_id=? AND creator_id=?", bookID, user.ID).First(&record).Error
	return record, err
}

func FindUserReviewByID(db *gorm.DB, user models.User, id uint64) (models.Review, error) {
	record := models.Review{}
	err := db.Where("id=? AND creator_id=?", id, user.ID).First(&record).Error
	return record, err
}

func FindBooksCountInLibrary(db *gorm.DB, user models.User) (int, error) {
	var count int
	err:=db.Model(&models.Review{}).Where("creator_id = ?", user.ID).Count(&count).Error
	return count, err
}

func FindUserReviewsByStatus(db *gorm.DB, user models.User, status int) ([]models.Review, error) {
	records := []models.Review{}
	err := db.Where("creator_id=? AND status=?", user.ID, status).Find(&records).Error
	return records, err
}

func FindUserReviews(db *gorm.DB, user models.User) ([]models.Review, error) {
	records := []models.Review{}
	err := db.Where("creator_id=?", user.ID).Find(&records).Error
	return records, err
}

func FindUserReviewsMapByBookID(db *gorm.DB, user models.User) (map[uint64]models.Review, error) {
	mapByBookID := make(map[uint64]models.Review, 0)
	records, err := FindUserReviews(db, user)
	if err != nil {
		return mapByBookID, err
	}

	for _, record := range records {
		mapByBookID[record.BookID] = record
	}
	return mapByBookID, nil
}

func UpdateReviewWithGoodreadsID(db *gorm.DB, user models.User, review models.Review, goodreadsID string) (models.Review, error) {
	update := models.Review{
		GoodreadsID: goodreadsID,
	}
	err := db.Model(&review).Debug().Where("creator_id=?", user.ID).Update(update).Error
	return review, err
}
