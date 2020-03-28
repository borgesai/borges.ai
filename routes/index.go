package routes

import (
	"borges.ai/models"
	"borges.ai/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("IndexHandler"))

	if utils.IsCustomHost(req) {
		UserHandler(w, req)
		return
	}

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
	// redirect if we have user
	if user.ID > 0 {
		http.Redirect(w, req, "/@"+user.Username, http.StatusFound)
		return
	}
	r.HTML(w, http.StatusOK, "index", IndexPageData{
		Page:       "home",
		Title:      "Borges | Your friendly librarian",
		User:       user,
		ShowFooter: true,
	})
}

func ChangelogHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("ChangelogHandler"))

	if utils.IsCustomHost(req) {
		UserHandler(w, req)
		return
	}

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
	r.HTML(w, http.StatusOK, "changelog", IndexPageData{
		Page:       "home",
		Title:      "Borges | Changelog",
		User:       user,
		ShowFooter: true,
	})
}

func AboutHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("AboutHandler"))

	if utils.IsCustomHost(req) {
		UserHandler(w, req)
		return
	}

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
	r.HTML(w, http.StatusOK, "about", IndexPageData{
		Page:       "home",
		Title:      "Borges | About",
		User:       user,
		ShowFooter: true,
	})
}

func Borges404Handler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		user := GetUser(req)
		db, err := models.NewDB()
		if err != nil {
			log.WithError(err).Error("failed to get connection")
		}
		defer db.Close()
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Page with this address wasn't found.",
			Title:   "Not found",
			Page:    "error",
		})
	}
	return http.HandlerFunc(fn)
}
