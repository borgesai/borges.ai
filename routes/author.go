package routes

import (
	"borges.ai/data"
	"borges.ai/models"
	"borges.ai/services"
	"borges.ai/utils"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func AuthorHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("AuthorHandler"))
	vars := mux.Vars(req)
	authorSlug := vars["author-name-slug"]
	user := GetUser(req)
	db, err := models.NewDB()
	if err != nil {
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Internal error",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	defer db.Close()

	author, err := data.FindAuthorBySlug(db, authorSlug)
	if err != nil {
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Couldn't find author. Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	log.WithField("author", author).Info("author found")
	books, err := services.FindAuthorsBooks(db, user, author.ID)
	if err != nil {
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Couldn't find author. Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	r.HTML(w, http.StatusOK, "author", AuthorPageData{
		Page:   "author",
		Title:  author.Name,
		User:   user,
		Author: author,
		Books:  books,
	})
}
