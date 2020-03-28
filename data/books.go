package data

import (
	"borges.ai/models"
	"borges.ai/text"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func FindOrCreateBook(db *gorm.DB, user models.User, goodreadsID, bestEditionGoodreadsID, title, subtitle, titleWithoutSeries, description string, authors []models.Author) (models.Book, error) {
	record := models.Book{}
	record.Title = title
	record.Subtitle = subtitle
	record.Slug = text.Slug(title)
	record.TitleWithoutSeries = titleWithoutSeries
	record.Description = description

	record.GoodreadsID = goodreadsID
	record.BestEditionGoodreadsID = bestEditionGoodreadsID

	// real authors
	authorsIDs := make([]string, 0)
	authorsNames := make([]string, 0)

	// illustrators
	illustratorsIDs := make([]string, 0)
	illustratorsNames := make([]string, 0)

	// translators
	translatorsIDs := make([]string, 0)
	translatorsNames := make([]string, 0)

	for _, author := range authors {
		authorIDStr := strconv.FormatUint(author.ID, 10)
		if author.Role == "" {
			authorsIDs = append(authorsIDs, authorIDStr)
			authorsNames = append(authorsNames, author.Name)
		} else if strings.ToLower(author.Role) == "illustrator" {
			illustratorsIDs = append(illustratorsIDs, authorIDStr)
			illustratorsNames = append(illustratorsNames, author.Name)
		} else if strings.ToLower(author.Role) == "translator" {
			translatorsIDs = append(translatorsIDs, authorIDStr)
			translatorsNames = append(translatorsNames, author.Name)
		}
	}
	authorsNamesStr := strings.Join(authorsNames, ", ")
	illustratorsNamesStr := strings.Join(illustratorsNames, ", ")
	translatorsNamesStr := strings.Join(translatorsNames, ", ")

	record.AuthorsIDs = authorsIDs
	record.AuthorsCache = authorsNames
	record.AuthorsCacheStr = authorsNamesStr

	record.IllustratorsIDs = illustratorsIDs
	record.IllustratorsCache = illustratorsNames
	record.IllustratorsCacheStr = illustratorsNamesStr

	record.TranslatorsIDs = translatorsIDs
	record.TranslatorsCache = translatorsNames
	record.TranslatorsCacheStr = translatorsNamesStr

	record.CreatorID = user.ID
	record.UpdatedAt = time.Now()
	err := db.Debug().Where(models.Book{GoodreadsID: goodreadsID}).Assign(record).FirstOrCreate(&record).Error
	if err != nil {
		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
			log.WithField("record", record).WithError(err).Info("failed to find or create a book (duplicate)")
			return FindBookByGoodreadsID(db, goodreadsID)
		}
		log.WithField("record", record).WithError(err).Error("failed to find or create a book")
		return record, err
	}
	log.WithField("record", record.ID).Info("found or created book")
	return record, nil
}

func FindBooksByIDs(db *gorm.DB, user models.User, booksIDs []uint64) ([]models.Book, error) {
	books := []models.Book{}
	err := db.Where("id IN (?)", booksIDs).Find(&books).Error
	return books, err
}

func FindBookByID(db *gorm.DB, id uint64) (models.Book, error) {
	record := models.Book{}
	err := db.First(&record, id).Error
	log.WithField("record", record).Debug("found book by id")
	return record, err
}

func FindBooksByAuthorID(db *gorm.DB, user models.User, authorID uint64) ([]models.Book, error) {
	authorIDStr := strconv.FormatUint(authorID, 10)
	books := []models.Book{}
	err := db.Where("?=ANY(authors_ids)", authorIDStr).Find(&books).Error
	return books, err
}

func FindBooksByQuery(db *gorm.DB, user models.User, query string) ([]models.Book, error) {
	books := []models.Book{}
	err := db.Debug().Where("searchable_text @@ to_tsquery('english', ?)", query).Find(&books).Limit(10).Error

	// adding authors. We only limit to 10. So this loops are fine
	authorsIDs := make([]uint64, 0)
	for _, book := range books {
		for _, authorID := range book.AuthorsIDsUint {
			authorsIDs = append(authorsIDs, authorID)
		}
	}

	authors, err := FindAuthorsByIDsAsMap(db, user, authorsIDs)
	if err != nil {
		return books, err
	}
	updatedBooks := make([]models.Book, 0)
	for _, book := range books {
		for _, authorID := range book.AuthorsIDsUint {
			author := authors[authorID]
			if author.ID > 0 {
				book.Authors = append(book.Authors, author)
			}
		}
		updatedBooks = append(updatedBooks, book)
	}
	return updatedBooks, err
}
func FindBookByISBN(db *gorm.DB, user models.User, isbn string) (models.Book, error) {
	edition, err := FindEditionByISBN(db, user, isbn)
	book := models.Book{}
	if err != nil {
		return book, err
	}

	book, err = FindBookByID(db, edition.BookID)
	if book.ID > 0 {
		book.Edition = edition
		authors, err := FindAuthorsByIDs(db, user, book.AuthorsIDsUint)
		if err != nil {
			return book, err
		}
		if len(authors) > 0 {
			book.Authors = authors
		}
	}

	return book, err
}

func FindBookByGoodreadsID(db *gorm.DB, goodreadsID string) (models.Book, error) {
	record := models.Book{}
	err := db.Where("goodreads_id=?", goodreadsID).First(&record).Error
	log.WithField("user", record).Debug("found book by goodreads_id")
	return record, err
}

func UpdateBookOriginalPublicationYear(db *gorm.DB, user models.User, bookID uint64, year int) (models.Book, error) {
	model := models.Book{}
	model.ID = bookID
	err := db.Model(&model).Update("original_year", year).Error
	return model, err
}

func UpdateBookData(db *gorm.DB, user models.User, bookID uint64, title, subtitle, originalTitle, bestEditionGoodreadsID string) (models.Book, error) {
	model := models.Book{}
	model.ID = bookID
	update := models.Book{
		Title:           title,
		Subtitle:  subtitle,
		OriginalTitle: originalTitle,
		BestEditionGoodreadsID:  bestEditionGoodreadsID,
	}
	err := db.Debug().Model(&model).Update(update).Error
	return model, err
}

func UpdateBookNumberOfPages(db *gorm.DB, user models.User, bookID uint64, numPages int) (models.Book, error) {
	model := models.Book{}
	model.ID = bookID
	err := db.Model(&model).Update("num_pages", numPages).Error
	return model, err
}

func UpdateBookShelves(db *gorm.DB, user models.User, bookID uint64, shelves []string) (models.Book, error) {
	model := models.Book{}
	model.ID = bookID
	err := db.Model(&model).Update("goodreads_shelves", shelves).Error
	return model, err
}

func UpdateBookBestEditionGoodreadsID(db *gorm.DB, user models.User, bookID uint64, bestEditionGoodreadsID string) (models.Book, error) {
	model := models.Book{}
	model.ID = bookID
	err := db.Model(&model).Update("best_edition_goodreads_id", bestEditionGoodreadsID).Error
	return model, err
}
