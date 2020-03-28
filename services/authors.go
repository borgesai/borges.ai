package services

import (
	"borges.ai/data"
	"borges.ai/models"
	"github.com/jinzhu/gorm"
)

func FindAuthorsBooks(db *gorm.DB, user models.User, authorID uint64) ([]models.Book, error){
	books, err := data.FindBooksByAuthorID(db, user, authorID)
	if err !=nil {
		return books, err
	}
	booksIDs := make([]uint64, 0)
	for _, book := range books {
		booksIDs = append(booksIDs, book.ID)
	}
	editionsMap, err := data.FindEditionsByBooksIDsAsMap(db, user, booksIDs)
	if err !=nil {
		return books, err
	}
	updatedBooks := make([]models.Book, 0)
	for _, book := range books {
		book.Edition = editionsMap[book.ID]
		updatedBooks = append(updatedBooks, book)
	}
	return updatedBooks, err
}
