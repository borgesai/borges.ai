package services

import (
	"borges.ai/data"
	"borges.ai/models"
	"borges.ai/utils"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

func ImportBookEdition(db *gorm.DB, user models.User, bookID uint64, goodreadsBookID, goodreadsID, isbn, isbn13, asin, publisherName string,
	numPages, publicationYear int, isEbook bool, format string) (models.BookEdition, error) {
	defer utils.Duration(utils.Track("ImportBookEdition"))

	publisherName = strings.TrimSpace(publisherName)

	eventPayload := map[string]string{
		"book_id":           strconv.FormatUint(bookID, 10),
		"goodreads_book_id": goodreadsBookID,
		"goodreads_id":      goodreadsID,
		"isbn":              isbn,
		"isbn13":            isbn13,
		"asin":              asin,
		"publisher":         publisherName,
		"num_pages":         strconv.Itoa(numPages),
		"publication_year":  strconv.Itoa(publicationYear),
		"is_ebook":          strconv.FormatBool(isEbook),
		"format":            format,
	}
	event, err := data.CreateCreateEvent(db, user, "book", "create book", eventPayload, false)
	if err != nil {
		return models.BookEdition{}, err
	}
	bookEdition, err := data.FindOrCreateBookEdition(db, user, bookID, goodreadsBookID, goodreadsID, isbn, isbn13, asin, publisherName, numPages, publicationYear, isEbook, format)
	_, err = data.TransactCreateEvent(db, user, event.ID, bookEdition.ID)

	return bookEdition, err
}
