package routes

import (
	"borges.ai/models"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func HandleAuthFailure() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		user := GetUser(req)
		db, err := models.NewDB()
		if err != nil {
			log.WithError(err).Error("failed to get connection")
		}
		defer db.Close()
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Failed to authenticate",
			Title:   "Hmmmmm",
			Page:    "error",
		})
	}
	return http.HandlerFunc(fn)
}
