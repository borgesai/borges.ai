package data

import (
	"borges.ai/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)


func FindOrCreateBookEdition(db *gorm.DB, user models.User, bookID uint64, goodreadsBookID, goodreadsID, isbn, isbn13, asin, publisher string, numPages, publicationYear int, isEbook bool, format string) (models.BookEdition, error) {
	record := models.BookEdition{}
	record.BookID = bookID
	record.GoodreadsBookID = goodreadsBookID
	record.GoodreadsID = goodreadsID
	record.ISBN = isbn
	record.ISBN13 = isbn13
	record.ASIN = asin
	record.PublicationYear = publicationYear
	record.NumPages = numPages
	record.Publisher = publisher
	record.IsEbook = isEbook
	record.Format = strings.ToLower(format)

	record.CreatorID = user.ID
	record.UpdatedAt = time.Now()
	err := db.Where(models.BookEdition{BookID: bookID, GoodreadsID: goodreadsID}).Assign(record).FirstOrCreate(&record).Error
	if err != nil {
		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
			log.WithField("record", record).WithError(err).Info("failed to find or create edition (duplicate)")
			return FindEditionByGoodreadsID(db, goodreadsID)
		}
		log.WithField("record", record).WithError(err).Error("failed to find or create a book edition")
		return record, err
	}
	log.WithField("record", record.ID).Info("found or created book edition")
	return record, nil
}

func FindEditionsByBooksIDs(db *gorm.DB, user models.User, booksIDs []uint64) ([]models.BookEdition, error) {
	editions := []models.BookEdition{}
	err := db.Where("book_id IN (?)", booksIDs).Find(&editions).Error
	return editions, err
}

func FindEditionsByBooksIDsAsMap(db *gorm.DB, user models.User, booksIDs []uint64) (map[uint64]models.BookEdition, error) {
	editions, err := FindEditionsByBooksIDs(db, user, booksIDs)
	editionsMap := make(map[uint64]models.BookEdition, 0)
	if err != nil {
		return editionsMap, err
	}
	// TODO merge editions
	for _, edition := range editions {
		editionsMap[edition.ID] = edition
	}
	return editionsMap, err
}

func FindEditionsByBookID(db *gorm.DB, user models.User, booksID uint64) ([]models.BookEdition, error) {
	editions := []models.BookEdition{}
	err := db.Where("book_id=?", booksID).Find(&editions).Error
	return editions, err
}

func FindEditionsByBookIDAsMap(db *gorm.DB, user models.User, bookID uint64) (map[uint64]models.BookEdition, error) {
	editions, err := FindEditionsByBookID(db, user, bookID)
	editionsMap := make(map[uint64]models.BookEdition, 0)
	if err != nil {
		return editionsMap, err
	}
	for _, edition := range editions {
		editionsMap[edition.ID] = edition
	}
	return editionsMap, err
}

func FindEditionByISBN(db *gorm.DB, user models.User, isbn string) (models.BookEdition, error) {
	edition := models.BookEdition{}
	// this for starts with
	query := isbn + "%"
	err := db.Debug().Where("isbn LIKE(?) OR isbn13 LIKE(?)", query, query).First(&edition).Error
	return edition, err
}

func FindEditionByGoodreadsID(db *gorm.DB, goodreadsID string) (models.BookEdition, error) {
	record := models.BookEdition{}
	err := db.Where("goodreads_id=?", goodreadsID).First(&record).Error
	log.WithField("user", record).Debug("found book edition by goodreads_id")
	return record, err
}
