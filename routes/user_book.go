package routes

import (
	"borges.ai/data"
	"borges.ai/models"
	"borges.ai/services"
	"borges.ai/text"
	"borges.ai/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func UserBookHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("UserBookHandler"))
	vars := mux.Vars(req)
	username := strings.ToLower(vars["username"])
	bookShortIDStr := vars["book-short-id"]
	bookNameSlug := vars["book-name-slug"]
	bookID := text.HashToNumber(bookShortIDStr)
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
		canonicalURL = "https://" + profile.CustomDomain + "/b/" + bookNameSlug + "-" + bookShortIDStr
	}

	book, err := services.FindUserBookByID(db, profile, uint64(bookID))
	if err != nil {
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	isMine := user.ID == profile.ID && user.ID > 0
	sessionUserReview := models.Review{}

	if !isMine && user.ID > 0 {
		sessionUserReview, err = services.FindUserReviewOptional(db, user, uint64(bookID))
	}

	if err != nil {
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}

	r.HTML(w, http.StatusOK, "user_book", UserBookPageData{
		Page:              "user_book",
		Title:             book.Title + " | Review by " + profile.Username,
		User:              user,
		Readonly:          isCustomDomain || !isMine,
		IsMine:            isMine,
		IsCustomDomain:    isCustomDomain,
		Profile:           profile,
		Book:              book,
		SessionUserReview: sessionUserReview,
		CanonicalURL:      canonicalURL,
	})
}

func UserChangeBookStatusHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("UserChangeBookStatusHandler"))
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
	bookShortIDStr := vars["book-short-id"]
	bookSlug := vars["book-name-slug"]
	bookID := text.HashToNumber(bookShortIDStr)
	statusStr := req.FormValue("status")

	if statusStr == "" {
		r.HTML(w, http.StatusBadRequest, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	status, _ := strconv.ParseInt(statusStr, 10, 32)
	editionIDStr := req.FormValue("edition")
	var editionID uint64 = 0
	if editionIDStr != "" {
		editionID, _ = strconv.ParseUint(editionIDStr, 10, 64)
	}
	if status == 0 {
		r.HTML(w, http.StatusBadRequest, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}

	if strings.ToLower(user.Username) != username {
		r.HTML(w, http.StatusForbidden, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	log.WithField("bookID", bookID).WithField("status", status).Info("submitted status")
	_, err = services.UpdateUserBookStatus(db, user, uint64(bookID), editionID, int(status))
	if err != nil {
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	newURL := "/@" + username + "/b/" + bookSlug + "-" + bookShortIDStr
	http.Redirect(w, req, newURL, http.StatusMovedPermanently)
}

func UserChangeBookRatingHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("UserChangeBookRatingHandler"))
	user := GetUser(req)
	db, err := models.NewDB()
	if err != nil {
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{})
		return
	}
	defer db.Close()
	vars := mux.Vars(req)
	username := strings.ToLower(vars["username"])
	bookShortIDStr := vars["book-short-id"]
	bookSlug := vars["book-name-slug"]
	bookID := text.HashToNumber(bookShortIDStr)
	ratingStr := req.FormValue("rating")
	if ratingStr == "" {
		r.HTML(w, http.StatusBadRequest, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	rating, _ := strconv.ParseInt(ratingStr, 10, 32)
	if rating == 0 {
		r.HTML(w, http.StatusBadRequest, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	editionIDStr := req.FormValue("edition")
	var editionID uint64 = 0
	if editionIDStr != "" {
		editionID, _ = strconv.ParseUint(editionIDStr, 10, 64)
	}
	if strings.ToLower(user.Username) != username {
		r.HTML(w, http.StatusForbidden, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	log.WithField("bookID", bookID).WithField("rating", rating).Info("submitted rating")
	_, err = services.UpdateUserBookRating(db, user, uint64(bookID), editionID, int(rating))
	if err != nil {
		r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
			User:    user,
			Message: "Something went wrong",
			Title:   "Hmmmmm",
			Page:    "error",
		})
		return
	}
	newURL := "/@" + username + "/b/" + bookSlug + "-" + bookShortIDStr
	http.Redirect(w, req, newURL, http.StatusMovedPermanently)
}

type UpdateContentForm struct {
	Value   string
	Edition string
}

// API routes //

func UserUpdateReviewAPIHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("UserUpdateReviewAPIHandler"))
	vars := mux.Vars(req)
	username := strings.ToLower(vars["username"])
	bookShortIDStr := vars["book-short-id"]
	bookID := text.HashToNumber(bookShortIDStr)
	var form UpdateContentForm
	b, err := ioutil.ReadAll(req.Body)
	json.Unmarshal(b, &form)
	if err != nil {
		log.WithError(err).Info("error")
		r.JSON(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user := GetUser(req)
	db, err := models.NewDB()
	if err != nil {
		r.JSON(w, http.StatusBadRequest, "Something went wrong")
		return
	}
	defer db.Close()
	if strings.ToLower(user.Username) != username {
		r.JSON(w, http.StatusBadRequest, "Invalid permission")
		return
	}
	editionIDStr := form.Edition
	var editionID uint64 = 0
	if editionIDStr != "" {
		editionID, _ = strconv.ParseUint(editionIDStr, 10, 64)
	}
	log.WithField("bookID", bookID).WithField("review", form.Value).Info("submitted review")
	review, err := services.UpdateReview(db, user, uint64(bookID), editionID, form.Value)
	if err != nil {
		r.JSON(w, http.StatusBadRequest, "Failed to make an update")
		return
	}
	resp := make(map[string]interface{}, 0)
	resp["id"] = review.ID
	resp["content"] = review.Content
	resp["content_html"] = review.ContentHTML

	r.JSON(w, http.StatusOK, resp)
}

func UserReadingNoteAPIHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("UserReadingNoteAPIHandler"))
	vars := mux.Vars(req)
	username := strings.ToLower(vars["username"])
	bookShortIDStr := vars["book-short-id"]
	bookID := text.HashToNumber(bookShortIDStr)
	var form UpdateContentForm
	b, err := ioutil.ReadAll(req.Body)
	json.Unmarshal(b, &form)
	if err != nil {
		log.WithError(err).Info("error")
		r.JSON(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user := GetUser(req)
	db, err := models.NewDB()
	if err != nil {
		r.JSON(w, http.StatusBadRequest, "Something went wrong")
		return
	}
	defer db.Close()
	editionIDStr := form.Edition
	var editionID uint64 = 0
	if editionIDStr != "" {
		editionID, _ = strconv.ParseUint(editionIDStr, 10, 64)
	}
	if strings.ToLower(user.Username) != username {
		r.JSON(w, http.StatusBadRequest, "Invalid permission")
		return
	}
	log.WithField("bookID", bookID).WithField("review", form.Value).Info("submitted note")
	reading, err := services.UpdateReadingNote(db, user, uint64(bookID), editionID, form.Value)
	if err != nil {
		r.JSON(w, http.StatusBadRequest, "Failed to make an update")
		return
	}
	resp := make(map[string]interface{}, 0)
	resp["id"] = reading.ID
	resp["start_date"] = reading.StartDate
	resp["finish_date"] = reading.FinishDate
	resp["duration"] = reading.Duration

	r.JSON(w, http.StatusOK, resp)
}

func UserReadingStartDateAPIHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("UserReadingStartDateAPIHandler"))
	vars := mux.Vars(req)
	username := strings.ToLower(vars["username"])
	bookShortIDStr := vars["book-short-id"]
	bookID := text.HashToNumber(bookShortIDStr)
	//readingID := vars["reading-id"]
	var form UpdateContentForm
	b, err := ioutil.ReadAll(req.Body)
	json.Unmarshal(b, &form)
	if err != nil {
		log.WithError(err).Info("error")
		r.JSON(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user := GetUser(req)
	db, err := models.NewDB()
	if err != nil {
		r.JSON(w, http.StatusBadRequest, "Something went wrong")
		return
	}
	defer db.Close()
	if strings.ToLower(user.Username) != username {
		r.JSON(w, http.StatusBadRequest, "Invalid permission")
		return
	}
	log.WithField("bookID", bookID).WithField("date", form.Value).Info("submitted start date")
	// TODO for now
	bookEditionID := 0
	startDateStr := form.Value
	startDate := utils.ParseDate(startDateStr)
	reading, err := services.UpdateReadingStartDate(db, user, uint64(bookID), uint64(bookEditionID), startDateStr, startDate)
	if err != nil {
		r.JSON(w, http.StatusBadRequest, "Failed to make an update")
		return
	}
	reading.AfterFind()
	resp := make(map[string]interface{}, 0)
	resp["id"] = reading.ID
	resp["start_date"] = utils.FormatDate(reading.StartDate)
	resp["finish_date"] = utils.FormatDate(reading.FinishDate)
	resp["duration"] = reading.Duration

	r.JSON(w, http.StatusOK, resp)
}

func UserReadingFinishDateAPIHandler(w http.ResponseWriter, req *http.Request) {
	defer utils.Duration(utils.Track("UserReadingFinishDateAPIHandler"))
	vars := mux.Vars(req)
	username := strings.ToLower(vars["username"])
	bookShortIDStr := vars["book-short-id"]
	bookID := text.HashToNumber(bookShortIDStr)
	//readingID := vars["reading-id"]
	var form UpdateContentForm
	b, err := ioutil.ReadAll(req.Body)
	json.Unmarshal(b, &form)
	if err != nil {
		log.WithError(err).Info("error")
		r.JSON(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user := GetUser(req)
	db, err := models.NewDB()
	if err != nil {
		r.JSON(w, http.StatusBadRequest, "Something went wrong")
		return
	}
	defer db.Close()
	if strings.ToLower(user.Username) != username {
		r.JSON(w, http.StatusBadRequest, "Invalid permission")
		return
	}
	log.WithField("bookID", bookID).WithField("date", form.Value).Info("submitted finish date")
	// TODO for now
	bookEditionID := 0
	finishDateStr := form.Value
	finishDate := utils.ParseDate(finishDateStr)
	reading, err := services.UpdateReadingFinishDate(db, user, uint64(bookID), uint64(bookEditionID), finishDateStr, finishDate)
	if err != nil {
		r.JSON(w, http.StatusBadRequest, "Failed to make an update")
		return
	}

	reading.AfterFind()
	resp := make(map[string]interface{}, 0)
	resp["id"] = reading.ID
	resp["start_date"] = utils.FormatDate(reading.StartDate)
	resp["finish_date"] = utils.FormatDate(reading.FinishDate)
	resp["duration"] = reading.Duration

	r.JSON(w, http.StatusOK, resp)
}
