package goodreads_import

import (
	"borges.ai/goodreads"
	"borges.ai/jobs"
	"borges.ai/models"
	"borges.ai/services"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"os"
)

func SyncBookFromGoodreads(db *gorm.DB, user models.User, book models.Book, goodreadsBookEditionID string) error {
	if book.ID == 0 {
		log.Info("nothing to do here")
		return nil
	}
	log.WithField("id", book.ID).WithField("goodreadsBookEditionID", goodreadsBookEditionID).Info("sync book")
	grc := goodreads.NewClient(os.Getenv("GOODREADS_KEY"))
	goodreadsBook, err := grc.BookByGoodreadsID(goodreadsBookEditionID)

	if err != nil {
		log.WithField("book_id", book.ID).WithField("goodreads_book_id", goodreadsBookEditionID).Error("failed to get goodreads book by id")
		return err
	}
	originalPublicationYear := goodreadsBook.Work.OriginalPublicationYear
	originalTitle := goodreadsBook.Work.OriginalTitle
	bookTitle, bookSubtitle := makeTitleAndSubtitle(goodreadsBook.Title)
	publishedYear := goodreadsBook.PublicationYear
	if originalPublicationYear == 0 {
		// this one we assign year of edition to the book
		originalPublicationYear = publishedYear
	}

	isbn13 := goodreadsBook.ISBN13
	isbn := goodreadsBook.ISBN
	asin := goodreadsBook.ASIN
	publisher := goodreadsBook.Publisher

	numPages := goodreadsBook.NumPages
	if numPages == 0 {
		// take num of pages from the book itself
		// if book is Kindle we wouldn't have pages, but we can import that from the best book too
		// this still fail if we are importing kindle book and best edition is not there
		numPages = book.NumPages
	}
	isEbook := goodreadsBook.IsEbook

	goodreadsBookID := goodreadsBook.Work.ID
	goodreadsID := goodreadsBook.ID

	isBestEdition := goodreadsID == goodreadsBook.Work.BestBookID
	log.Info(isBestEdition)

	_, err = services.ImportBookEdition(db, user, book.ID, goodreadsBookID, goodreadsID, isbn, isbn13, asin, publisher, numPages, publishedYear, isEbook, goodreadsBook.Format)
	if err != nil {
		log.WithField("book_id", book.ID).WithError(err).Error("failed to create book edition")
		// NOTE do not return. Not essential. Just log
	}

	if book.NumPages == 0 && numPages > 0 {
		_, err = services.UpdateBookNumberOfPages(db, user, book.ID, numPages)
		if err != nil {
			log.WithField("book_id", book.ID).WithError(err).Error("failed to create book edition")
			// NOTE do not return. Not essential. Just log
		}

	}

	if book.OriginalYear != publishedYear {
		if originalPublicationYear != 0 && (book.OriginalYear == 0 || originalPublicationYear < book.OriginalYear) {
			log.WithField("id", book.ID).WithField("year", originalPublicationYear).Info("found original publication year")
			_, err = services.UpdateBookOriginalPublicationYear(db, user, book.ID, originalPublicationYear)
			// log but do not stop
			if err != nil {
				log.WithField("book_id", book.ID).WithField("goodreads_book_id", book.GoodreadsID).Error("failed to save original year")
			}
		}
	}

	if isBestEdition {
		_, err = services.UpdateBookData(db, user, book.ID, bookTitle, bookSubtitle, originalTitle, goodreadsBook.Work.BestBookID)
		// log but do not stop
		if err != nil {
			log.WithField("book_id", book.ID).WithField("goodreads_book_id", book.GoodreadsID).Error("failed to save original title")
		}
	} else {
		err = jobs.SyncGoodreadsBook(user, book.ID, goodreadsBook.Work.BestBookID)
		if err != nil {
			log.WithError(err).Error("failed to publish message")
			return err
		}
	}

	// update with a larger list
	if len(goodreadsBook.PopularShelves) > len(book.GoodreadsShelves) {
		shelvesNames := make([]string, 0)
		for _, shelf := range goodreadsBook.PopularShelves {
			shelvesNames = append(shelvesNames, shelf.Name)
		}
		if len(shelvesNames) > 0 {
			_, err = services.UpdateBookShelves(db, user, book.ID, shelvesNames)
			if err != nil {
				log.WithError(err).Error("failed to save book shelves")
			}
		}
	}
	return err
}
