package services

import (
	"borges.ai/data"
	"borges.ai/models"
	"borges.ai/utils"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"sync"
)

func FindUserBookByID(db *gorm.DB, user models.User, bookID uint64) (models.Book, error) {
	defer utils.Duration(utils.Track("FindUserBookByID"))
	var err error
	var book models.Book
	var review models.Review
	var readings []models.Reading
	var editionsMap map[uint64]models.BookEdition
	var wg sync.WaitGroup
	wg.Add(4)
	// 4 calls in parallel
	go (func() {
		book, err = data.FindBookByID(db, bookID)
		wg.Done()
	})()
	go (func() {
		reviewRecord, reviewError := data.FindReviewByBookID(db, user, bookID)
		if reviewError != nil && reviewError.Error() != "record not found" {
			err = reviewError
		}
		review = reviewRecord
		wg.Done()
	})()
	go (func() {
		readings, err = data.FindUserBookReadings(db, user, bookID)
		wg.Done()
	})()
	go (func() {
		editionsMap, err = data.FindEditionsByBookIDAsMap(db, user, bookID)
		wg.Done()
	})()

	wg.Wait()
	log.Info("pend")
	if err != nil || book.ID == 0 {
		return book, err
	}
	// not in parallel
	authors, _ := data.FindAuthorsByIDs(db, user, book.AuthorsIDsUint)
	if len(readings) > 0 {
		book.Readings = readings
	}

	if review.ID > 0 {
		book.Review = review
	}

	if len(authors) > 0 {
		book.Authors = authors
	}

	var bookEditionID uint64
	if review.BookEditionID > 0 {
		bookEditionID = review.BookEditionID
	}
	if bookEditionID == 0 {
		if len(readings) > 0 {
			if readings[0].BookEditionID > 0 {
				bookEditionID = readings[0].BookEditionID
			}
		}
	}
	if bookEditionID > 0 {
		book.Edition = editionsMap[bookEditionID]
	} else {
		book.Edition = models.BookEdition{}
	}

	return book, err
}

func FindBookByID(db *gorm.DB, user models.User, bookID uint64) (models.Book, error) {
	defer utils.Duration(utils.Track("FindUserBookByID"))
	var err error
	var book models.Book
	var editions []models.BookEdition
	var wg sync.WaitGroup
	wg.Add(2)
	// 2 calls in parallel
	go (func() {
		book, err = data.FindBookByID(db, bookID)
		wg.Done()
	})()
	go (func() {
		editions, err = data.FindEditionsByBookID(db, user, bookID)
		wg.Done()
	})()

	wg.Wait()
	log.Info("pend")
	if err != nil || book.ID == 0 {
		return book, err
	}
	// not in parallel
	authors, _ := data.FindAuthorsByIDs(db, user, book.AuthorsIDsUint)

	book.Authors = authors
	book.Editions = editions
	if len(editions) > 0 {
		// but to be honest we should get the best edition here. do loop
		book.Edition = editions[0]
	}

	return book, err
}

func FindBooksToReadByUser(db *gorm.DB, user models.User) ([]models.Book, error) {
	defer utils.Duration(utils.Track("FindBooksToReadByUser"))
	books := make([]models.Book, 0)
	reviews, err := data.FindUserReviewsByStatus(db, user, data.WANT_TO_READ_STATUS)
	if err != nil {
		log.WithError(err).Error("failed to find statues")
		return books, err
	}
	booksIDs := make([]uint64, 0)
	bookReviewMap := make(map[uint64]models.Review, 0)
	for _, review := range reviews {
		booksIDs = append(booksIDs, review.BookID)
		bookReviewMap[review.BookID] = review
	}
	books, err = data.FindBooksByIDs(db, user, booksIDs)
	updatedBooks := make([]models.Book, 0)
	for _, book := range books {
		book.Review = bookReviewMap[book.ID]
		updatedBooks = append(updatedBooks, book)
	}
	return updatedBooks, err
}

func FindBooksInProgressByUser(db *gorm.DB, user models.User) ([]models.Book, error) {
	defer utils.Duration(utils.Track("FindBooksInProgressByUser"))
	books := make([]models.Book, 0)
	reviews, err := data.FindUserReviewsByStatus(db, user, data.STARTED_STATUS)
	if err != nil {
		log.WithError(err).Error("failed to find reviews")
		return books, err
	}
	booksIDs := make([]uint64, 0)
	bookReviewMap := make(map[uint64]models.Review, 0)
	for _, review := range reviews {
		booksIDs = append(booksIDs, review.BookID)
		bookReviewMap[review.BookID] = review
	}
	books, err = data.FindBooksByIDs(db, user, booksIDs)
	updatedBooks := make([]models.Book, 0)
	for _, book := range books {
		book.Review = bookReviewMap[book.ID]
		updatedBooks = append(updatedBooks, book)
	}
	return updatedBooks, err
}

func FindBooksReadByUser(db *gorm.DB, user models.User) ([]models.Book, error) {
	defer utils.Duration(utils.Track("FindBooksReadByUser"))
	books := make([]models.Book, 0)
	reviews, err := data.FindUserReviewsByStatus(db, user, data.READ_STATUS)
	if err != nil {
		log.WithError(err).Error("failed to find statues")
		return books, err
	}
	booksIDs := make([]uint64, 0)
	bookReviewMap := make(map[uint64]models.Review, 0)
	for _, review := range reviews {
		booksIDs = append(booksIDs, review.BookID)
		bookReviewMap[review.BookID] = review
	}
	var editionsMap map[uint64]models.BookEdition
	var readingsMap map[uint64]models.Reading
	var wg sync.WaitGroup

	// 3 calls in parallel
	wg.Add(3)
	go (func() {
		books, err = data.FindBooksByIDs(db, user, booksIDs)
		wg.Done()
	})()
	go (func() {
		editionsMap, err = data.FindEditionsByBooksIDsAsMap(db, user, booksIDs)
		wg.Done()
	})()
	go (func() {
		readingsMap, err = data.FindUserReadingsMapByBookID(db, user)
		wg.Done()
	})()

	wg.Wait()
	if err != nil {
		log.WithError(err).Error("failed to find all the data")
		return books, err
	}
	updatedBooks := make([]models.Book, 0)
	for _, book := range books {
		book.Review = bookReviewMap[book.ID]
		book.Reading = readingsMap[book.ID]
		updatedBooks = append(updatedBooks, book)
	}

	return updatedBooks, err
}
