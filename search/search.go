package search

import (
	"borges.ai/data"
	"borges.ai/goodreads"
	"borges.ai/goodreads_import"
	"borges.ai/models"
	"borges.ai/text"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
)

func FindBooksAndAttachUserStatuses(db *gorm.DB, user models.User, query string) ([]models.Book, bool, error) {
	var books []models.Book
	var reviewsMap map[uint64]models.Review
	var err error
	var searchedByISBN bool
	var wg sync.WaitGroup
	wg.Add(2)
	go (func() {
		books, searchedByISBN, err = FindBooks(db, user, query)
		wg.Done()
	})()
	go (func() {
		reviewsMap, err = data.FindUserReviewsMapByBookID(db, user)
		wg.Done()
	})()
	wg.Wait()
	if err != nil {
		return books, false, err
	}
	updatedBooks := make([]models.Book, 0)
	for _, book := range books {
		book.Review = reviewsMap[book.ID]
		updatedBooks = append(updatedBooks, book)
	}
	return updatedBooks, searchedByISBN, nil
}

func FindBooks(db *gorm.DB, user models.User, query string) ([]models.Book, bool, error) {
	// flow
	// 1. remove all remove all the whitespaces and "-" and check if number
	potentialISBN := text.NormalizeISBN(query)
	grc := goodreads.NewClient(os.Getenv("GOODREADS_KEY"))
	if text.IsISBN(potentialISBN) {
		// 2a. if number two routines :> search db for isbn and isbn13 and search good reads. If found return
		var goodreadsBookByISBN *goodreads.Book
		var bookByISBN models.Book
		var err error
		var wg sync.WaitGroup
		wg.Add(2)
		go (func() {
			bookByISBN, err = data.FindBookByISBN(db, user, potentialISBN)
			if err != nil {
				log.WithField("isbn", potentialISBN).WithError(err).Error("failed to find book by isbn")
			}
			wg.Done()
		})()
		go (func() {
			goodreadsBookByISBN, err = grc.BookByISBN(potentialISBN)
			if err != nil {
				log.WithError(err).Error("failed to fetch book by isbn from goodreads")
			} else {
				log.WithField("book", goodreadsBookByISBN.ID).Info("fetched book by isbn from goodreads")
			}
			wg.Done()
		})()
		wg.Wait()
		log.WithField("db_id", bookByISBN.ID).Info("fetched book by isbn")
		if bookByISBN.ID > 0 {
			return []models.Book{bookByISBN}, true, nil
		}
		if goodreadsBookByISBN != nil {
			log.WithField("goodreads_id", goodreadsBookByISBN.ID).Info("fetched book by isbn, continue to goodreads")
			bookByISBN, importError := goodreads_import.ImportGoodreadsBookAndEdition(db, user, *goodreadsBookByISBN)
			log.WithField("book", bookByISBN.ID).Info("fetched book by isbn, finished goodreads")
			return []models.Book{bookByISBN}, true, importError
		}
		return []models.Book{}, true, err
	} else {
		// 2b. if not number two routines :> search db for title and search good reads. If found return
		var goodreadsBookByTitle *goodreads.Book
		var booksByTitle []models.Book
		var err error
		var wg sync.WaitGroup
		wg.Add(2)
		go (func() {
			pgQuery := strings.ReplaceAll(query, " ", " & ")
			booksByTitle, err = data.FindBooksByQuery(db, user, pgQuery)
			if err != nil {
				log.WithError(err).Error("failed to find book by title")
			}
			wg.Done()
		})()
		go (func() {
			goodreadsBookByTitle, err = grc.BookByTitle(query)
			if err != nil {
				log.WithError(err).Error("failed to fetch book by title from goodreads")
			} else {
				log.WithField("book", goodreadsBookByTitle).Info("fetched book by title from goodreads")
			}
			wg.Done()
		})()
		wg.Wait()
		// TODO we can also add this one in parallel https://www.goodreads.com/api/index#search.books
		// there might be some problem with it, but I can try it anyway
		if goodreadsBookByTitle != nil {
			newBook, importError := goodreads_import.ImportGoodreadsBookAndEdition(db, user, *goodreadsBookByTitle)
			if importError != nil {
				log.WithError(err).Error("import error")
			} else {
				booksByTitle = append(booksByTitle, newBook)
			}
		}
		return booksByTitle, false, nil

	}

	return []models.Book{}, false, nil
}
