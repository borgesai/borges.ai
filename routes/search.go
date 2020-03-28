package routes

import (
	"borges.ai/models"
	"borges.ai/search"
	"borges.ai/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func SearchHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("SearchHandler"))
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

	query := req.FormValue("query")

	log.WithField("query", query).Info("got query")

	if query == "" {
		r.HTML(w, http.StatusOK, "search", SearchPageData{
			Page:  "home",
			Title: "Borges | Search books",
			User:  user,
			Query: query,
		})
		return
	}
	books, searchedByISBN, err := search.FindBooksAndAttachUserStatuses(db, user, query)
	if err != nil {
		log.WithField("query", query).Error(err)
	}
	r.HTML(w, http.StatusOK, "search", SearchPageData{
		Page:           "home",
		Title:          "Borges | Search books",
		User:           user,
		Query:          query,
		Books:          books,
		SearchedByISBN: searchedByISBN,
	})
}
