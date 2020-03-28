package services

import (
	"borges.ai/data"
	"borges.ai/models"
	"borges.ai/utils"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

func ImportBook(db *gorm.DB, user models.User, goodreadsBookID, bestEditionGoodreadsID, title, subtitle, titleWithoutSeries, description string, authors []models.Author) (models.Book, error) {
	defer utils.Duration(utils.Track("ImportBook"))
	title = strings.TrimSpace(title)
	subtitle = strings.TrimSpace(subtitle)
	authorsNamesStr := ""
	for _, author := range authors {
		authorsNamesStr += ", " + author.Name
	}
	eventPayload := map[string]string{
		"goodreadsID":            goodreadsBookID,
		"bestEditionGoodreadsID": bestEditionGoodreadsID,
		"title":                  title,
		"subtitle":               subtitle,
		"titleWithoutSeries":     titleWithoutSeries,
		"description":            description,
		"authors":                authorsNamesStr,
	}
	event, err := data.CreateCreateEvent(db, user, "book", "create book", eventPayload, false)
	if err != nil {
		return models.Book{}, err
	}
	book, err := data.FindOrCreateBook(db, user, goodreadsBookID, bestEditionGoodreadsID, title, subtitle, titleWithoutSeries, description, authors)
	_, err = data.TransactCreateEvent(db, user, event.ID, book.ID)

	return book, err
}


func UpdateBookOriginalPublicationYear(db *gorm.DB, user models.User, bookID uint64, year int) (models.Book, error) {
	defer utils.Duration(utils.Track("UpdateBookOriginalPublicationYear"))
	eventPayload := map[string]string{
		"bookID":    strconv.FormatUint(bookID, 10),
		"year": strconv.Itoa(year),
	}
	event, err := data.CreateEvent(db, user, bookID, "book", "update",
		"original_year", strconv.Itoa(year), "number",
		"update original year", eventPayload, false)
	if err != nil {
		return models.Book{}, err
	}
	model, err := data.UpdateBookOriginalPublicationYear(db, user, bookID, year)
	if err != nil {
		return model, err
	}
	_, err = data.TransactEvent(db, user, event.ID)
	return model, err
}

func UpdateBookData(db *gorm.DB, user models.User, bookID uint64, title, subtitle, originalTitle, bestEditionGoodreadsID string) (models.Book, error) {
	defer utils.Duration(utils.Track("UpdateBookOriginalTitle"))
	eventPayload := map[string]string{
		"bookID":    strconv.FormatUint(bookID, 10),
		"title": title,
		"subtitle": subtitle,
		"original_title": originalTitle,
		"best_edition_goodreads_id": bestEditionGoodreadsID,
	}
	event, err := data.CreateEvent(db, user, bookID, "book", "update",
		"title", title, "string",
		"update book data", eventPayload, false)
	if err != nil {
		return models.Book{}, err
	}
	model, err := data.UpdateBookData(db, user, bookID, title, subtitle, originalTitle, bestEditionGoodreadsID)
	if err != nil {
		return model, err
	}
	_, err = data.TransactEvent(db, user, event.ID)
	return model, err
}

func UpdateBookNumberOfPages(db *gorm.DB, user models.User, bookID uint64, numPages int) (models.Book, error) {
	defer utils.Duration(utils.Track("UpdateBookNumberOfPages"))
	eventPayload := map[string]string{
		"bookID":    strconv.FormatUint(bookID, 10),
		"num_pages": strconv.Itoa(numPages),
	}
	event, err := data.CreateEvent(db, user, bookID, "book", "update",
		"original_year", strconv.Itoa(numPages), "number",
		"update num pages", eventPayload, false)
	if err != nil {
		return models.Book{}, err
	}
	model, err := data.UpdateBookNumberOfPages(db, user, bookID, numPages)
	if err != nil {
		return model, err
	}
	_, err = data.TransactEvent(db, user, event.ID)
	return model, err
}

func UpdateBookShelves(db *gorm.DB, user models.User, bookID uint64, shelves []string) (models.Book, error) {
	defer utils.Duration(utils.Track("UpdateBookShelves"))
	shelvesStr:= strings.Join(shelves, ", ")
	eventPayload := map[string]string{
		"bookID":    strconv.FormatUint(bookID, 10),
		"goodreads_shelves": shelvesStr,
	}
	event, err := data.CreateEvent(db, user, bookID, "book", "update",
		"goodreads_shelves", shelvesStr, "[]string",
		"update book shelves", eventPayload, false)
	if err != nil {
		return models.Book{}, err
	}
	model, err := data.UpdateBookShelves(db, user, bookID, shelves)
	if err != nil {
		return model, err
	}
	_, err = data.TransactEvent(db, user, event.ID)
	return model, err
}

func UpdateBookBestEditionGoodreadsID(db *gorm.DB, user models.User, bookID uint64, bestEditionGoodreadsID string) (models.Book, error) {
	defer utils.Duration(utils.Track("UpdateBookBestEditionGoodreadsID"))
	eventPayload := map[string]string{
		"bookID":    strconv.FormatUint(bookID, 10),
		"best_edition_goodreads_id": bestEditionGoodreadsID,
	}
	event, err := data.CreateEvent(db, user, bookID, "book", "update",
		"best_edition_goodreads_id", bestEditionGoodreadsID, "string",
		"update best edition", eventPayload, false)
	if err != nil {
		return models.Book{}, err
	}
	model, err := data.UpdateBookBestEditionGoodreadsID(db, user, bookID, bestEditionGoodreadsID)
	if err != nil {
		return model, err
	}
	_, err = data.TransactEvent(db, user, event.ID)
	return model, err
}