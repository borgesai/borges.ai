package routes

import (
	"borges.ai/models"
	"borges.ai/services"
	"borges.ai/text"
	"borges.ai/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func BookHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("BookHandler"))
	vars := mux.Vars(req)
	if utils.IsCustomHost(req) {
		UserBookHandler(w, req)
		return
	}
	bookShortIDStr := vars["book-short-id"]
	bookID := text.HashToNumber(bookShortIDStr)
	user := GetUser(req)
	db, err := models.NewDB()
	if err != nil {
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	defer db.Close()

	book, err := services.FindBookByID(db, user, uint64(bookID))
	if err != nil {
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	r.HTML(w, http.StatusOK, "book", BookPageData{
		Page:  "user",
		Title: book.Title,
		User:  user,
		Book:  book,
	})
}
