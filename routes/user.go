package routes

import (
	"borges.ai/data"
	"borges.ai/jobs"
	"borges.ai/models"
	"borges.ai/services"
	"borges.ai/utils"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func UserHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("UserHandler"))
	vars := mux.Vars(req)
	username := strings.ToLower(vars["username"])
	db, err := models.NewDB()
	user := models.User{}
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
	var profile models.User
	isCustomDomain := false

	if utils.IsCustomHost(req) {
		user, err = data.FindUserByCustomDomain(db, req.Host)
		if err != nil {
			log.WithError(err).Error("failed to find user by domain")
			r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
				User:    user,
				Message: "Something went wrong",
				Title:   "Hmmmmm",
				Page:    "error",
			})
			return
		}
		username = user.Username
		profile = user
		isCustomDomain = true
	} else {
		user = GetUser(req)
	}
	canonicalURL := ""
	if profile.CustomDomain != "" {
		canonicalURL = "https://" + profile.CustomDomain
	}

	if strings.ToLower(user.Username) == username {
		profile = user
	} else {
		profile, err = data.FindUserByUsername(db, username)
		if strings.ToLower(user.Username) == username {
			profile = user
		} else {
			profile, err = data.FindUserByUsername(db, username)
		}
	}

	books, err := services.FindBooksReadByUser(db, profile)
	if err != nil {
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	averageRating := models.AverageRating(books)
	reviewsCount := 0
	for _, book := range books {
		if book.Review.Content != "" {
			reviewsCount = reviewsCount + 1
		}
	}
	sort.Sort(models.ByReadingDateDesc(books))
	isMine := user.ID == profile.ID && user.ID > 0
	r.HTML(w, http.StatusOK, "user", UserPageData{
		Page:           "user",
		Title:          profile.Username + "| read books",
		Readonly:       isCustomDomain || !isMine,
		IsMine:         isMine,
		IsCustomDomain: isCustomDomain,
		ShowActions:    len(books) == 0,
		ReviewsCount:   reviewsCount,
		User:           user,
		Profile:        profile,
		Books:          books,
		AverageRating:  averageRating,
		Submenu:        "read",
		CanonicalURL:   canonicalURL,
	})
}

func UserSettingsHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("UserSettingsHandler"))
	user := GetUser(req)
	if user.ID == 0 {
		log.Info("access denied to the setting")
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
		return
	}

	profile := user
	r.HTML(w, http.StatusOK, "settings", UserSettingsPageData{
		Page:           "user",
		Title:          profile.Username + "| settings",
		Readonly:       false,
		IsMine:         true,
		IsCustomDomain: false,
		User:           user,
		Profile:        profile,
		Submenu:        "settings",
	})
}

func UserBestHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("UserBestHandler"))
	vars := mux.Vars(req)
	username := strings.ToLower(vars["username"])
	db, err := models.NewDB()
	user := models.User{}
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
	var profile models.User
	isCustomDomain := false

	if utils.IsCustomHost(req) {
		user, err = data.FindUserByCustomDomain(db, req.Host)
		if err != nil {
			log.WithError(err).Error("failed to find user by domain")
			r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
				User:    user,
				Message: "Something went wrong",
				Title:   "Hmmmmm",
				Page:    "error",
			})
			return
		}
		username = user.Username
		profile = user
		isCustomDomain = true

	} else {
		user = GetUser(req)
		if strings.ToLower(user.Username) == username {
			profile = user
		} else {
			profile, err = data.FindUserByUsername(db, username)
		}
	}

	canonicalURL := ""
	if profile.CustomDomain != "" {
		canonicalURL = "https://" + profile.CustomDomain + "/best"
	}

	books, err := services.FindBooksReadByUser(db, profile)
	if err != nil {
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}

	sort.Sort(models.ByRatingDesc(books))
	isMine := user.ID == profile.ID && user.ID > 0
	r.HTML(w, http.StatusOK, "user", UserPageData{
		Page:           "user",
		Title:          profile.Username + "| best books",
		Readonly:       isCustomDomain || !isMine,
		IsMine:         isMine,
		IsCustomDomain: isCustomDomain,
		ShowActions:    false,
		User:           user,
		Profile:        profile,
		Books:          books,
		Submenu:        "best",
		CanonicalURL:   canonicalURL,
	})
}

func UserWantToReadHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("UserWantToReadHandler"))
	vars := mux.Vars(req)
	username := strings.ToLower(vars["username"])
	db, err := models.NewDB()
	user := models.User{}
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
	var profile models.User
	isCustomDomain := false
	if utils.IsCustomHost(req) {
		user, err = data.FindUserByCustomDomain(db, req.Host)
		if err != nil {
			log.WithError(err).Error("failed to find user by domain")
			r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
				User:    user,
				Message: "Something went wrong",
				Title:   "Hmmmmm",
				Page:    "error",
			})
			return
		}
		username = user.Username
		profile = user
		isCustomDomain = true

	} else {
		user = GetUser(req)
		if strings.ToLower(user.Username) == username {
			profile = user
		} else {
			profile, err = data.FindUserByUsername(db, username)
		}
	}
	canonicalURL := ""
	if profile.CustomDomain != "" {
		canonicalURL = "https://" + profile.CustomDomain + "/to-read"
	}
	books, err := services.FindBooksToReadByUser(db, profile)
	if err != nil {
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	sort.Sort(models.ByStatusDesc(books))
	isMine := user.ID == profile.ID && user.ID > 0
	r.HTML(w, http.StatusOK, "user", UserPageData{
		Page:           "user",
		Title:          profile.Username + "| books to read",
		Readonly:       isCustomDomain || !isMine,
		IsMine:         isMine,
		IsCustomDomain: isCustomDomain,
		ShowActions:    false,
		User:           user,
		Profile:        profile,
		Books:          books,
		Submenu:        "to-read",
		CanonicalURL:   canonicalURL,
	})
}

func UserReadingHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("UserReadingHandler"))
	vars := mux.Vars(req)
	username := strings.ToLower(vars["username"])
	user := models.User{}
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
	var profile models.User
	isCustomDomain := false
	if utils.IsCustomHost(req) {
		user, err = data.FindUserByCustomDomain(db, req.Host)
		if err != nil {
			log.WithError(err).Error("failed to find user by domain")
			r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
				User:    user,
				Message: "Something went wrong",
				Title:   "Hmmmmm",
				Page:    "error",
			})
			return
		}
		username = user.Username
		profile = user
		isCustomDomain = true

	} else {
		user = GetUser(req)
		if strings.ToLower(user.Username) == username {
			profile = user
		} else {
			profile, err = data.FindUserByUsername(db, username)
		}
	}
	canonicalURL := ""
	if profile.CustomDomain != "" {
		canonicalURL = "https://" + profile.CustomDomain + "/reading"
	}
	books, err := services.FindBooksInProgressByUser(db, profile)
	if err != nil {
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	sort.Sort(models.ByStatusDesc(books))
	isMine := user.ID == profile.ID && user.ID > 0
	r.HTML(w, http.StatusOK, "user", UserPageData{
		Page:           "user",
		Title:          profile.Username + "| books in progress",
		Readonly:       isCustomDomain || !isMine,
		IsMine:         isMine,
		IsCustomDomain: isCustomDomain,
		ShowActions:    false,
		User:           user,
		Profile:        profile,
		Books:          books,
		Submenu:        "reading",
		CanonicalURL:   canonicalURL,
	})
}

func UserChartsHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("UserChartsHandler"))
	vars := mux.Vars(req)
	username := strings.ToLower(vars["username"])
	db, err := models.NewDB()
	user := models.User{}
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
	var profile models.User
	isCustomDomain := false

	if utils.IsCustomHost(req) {
		user, err = data.FindUserByCustomDomain(db, req.Host)
		if err != nil {
			log.WithError(err).Error("failed to find user by domain")
			r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
				User:    user,
				Message: "Something went wrong",
				Title:   "Hmmmmm",
				Page:    "error",
			})
			return
		}
		username = user.Username
		profile = user
		isCustomDomain = true

	} else {
		user = GetUser(req)
	}
	canonicalURL := ""
	if profile.CustomDomain != "" {
		canonicalURL = "https://" + profile.CustomDomain + "/charts"
	}

	if strings.ToLower(user.Username) == username {
		profile = user
	} else {
		profile, err = data.FindUserByUsername(db, username)
		if strings.ToLower(user.Username) == username {
			profile = user
		} else {
			profile, err = data.FindUserByUsername(db, username)
		}
	}

	books, err := services.FindBooksReadByUser(db, profile)
	if err != nil {
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	sort.Sort(models.ByReadingDateDesc(books))
	booksByYear := models.GroupBooksByFinishedYear(books)
	maxPages := 0
	maxBooks := 0
	for _, data := range booksByYear {
		if len(data.Items) > maxBooks {
			maxBooks = len(data.Items)
		}
		if data.TotalPages > maxPages {
			maxPages = data.TotalPages
		}
	}
	isMine := user.ID == profile.ID && user.ID > 0
	r.HTML(w, http.StatusOK, "user_charts", UserChartsPageData{
		Page:           "user",
		Title:          profile.Username + "| charts",
		Readonly:       isCustomDomain || !isMine,
		IsMine:         isMine,
		IsCustomDomain: isCustomDomain,
		User:           user,
		Profile:        profile,
		Books:          booksByYear,
		MaxBooks:       maxBooks,
		MaxPages:       maxPages,
		Submenu:        "charts",
		CanonicalURL:   canonicalURL,
	})
}

func UserSyncGoodreadsSettingHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("UserSyncGoodreadsSettingHandler"))
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
	vars := mux.Vars(req)
	username := strings.ToLower(vars["username"])
	enabledStr := req.FormValue("enabled")
	if strings.ToLower(user.Username) != username {
		r.HTML(w, http.StatusForbidden, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	enabled, _ := strconv.ParseBool(enabledStr)
	_, err = data.UpdateUserGoodreadsSyncFlag(db, user, enabled)
	if err != nil {
		log.WithError(err).Error("failed to sync")
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}

	newURL := "/settings"
	log.WithField("newURL", newURL).WithField("enabled", enabled).Info("going to redirect")
	http.Redirect(w, req, newURL, http.StatusMovedPermanently)
}

func UserSyncGoodreadsHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("UserSyncGoodreadsHandler"))
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
	vars := mux.Vars(req)
	username := strings.ToLower(vars["username"])
	redirectTo := req.FormValue("redirect_to")
	if strings.ToLower(user.Username) != username {
		r.HTML(w, http.StatusForbidden, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	_, err = jobs.SyncGoodreadsUserAndWait(user, true)
	if err != nil {
		log.WithError(err).Error("failed to sync")
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}

	newURL := redirectTo
	if newURL == "" {
		newURL = "/@" + username
	}
	log.WithField("newURL", newURL).Info("going to redirect")
	http.Redirect(w, req, newURL, http.StatusMovedPermanently)
}
