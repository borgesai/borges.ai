package routes

import (
	"borges.ai/data"
	"borges.ai/twitter"
	"net/http"
	"os"

	"borges.ai/models"
	"github.com/dghubble/oauth1"
	twitterOAuth1 "github.com/dghubble/oauth1/twitter"
	log "github.com/sirupsen/logrus"
)

func GetTwitterOAuth1Config() *oauth1.Config {
	var callbackURL = os.Getenv("SERVICE_URL") + "/twitter/callback"

	var oauth1Config = &oauth1.Config{
		ConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		CallbackURL:    callbackURL,
		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
	}
	return oauth1Config
}

var sessionName = "borges-session"
var sessionUserIDKey = "twitterID"
var sessionUserKey = "user"

// issueSession issues a cookie session after successful Twitter login
func IssueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		db, err := models.NewDB()
		if err != nil {
			log.WithError(err).Error("failed to connect to db")
			r.HTML(w, http.StatusInternalServerError, "error", ErrorPageData{
				Message: "Something went wrong",
				Title:   "Hmmmmm",
				Page:    "error",
			})
			return
		}
		defer db.Close()
		ctx := req.Context()
		twitterUser, twitterAccessToken, err := twitter.UserFromContext(ctx)
		if err != nil {
			log.WithError(err).Error("error getting data from twitter")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 2. Implement a success handler to issue some form of session
		session, err := Store.Get(req, sessionName)
		if err != nil {
			log.WithError(err).Error("failed to fetch session")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.WithField("twitter user", twitterUser).Info("got user data")
		user, err := data.FindOrCreateUserFromTwitter(db, twitterUser.ID, twitterUser.ScreenName, twitterUser.Email, twitterAccessToken, twitterUser.Name,
			twitterUser.Location, twitterUser.ProfileTextColor, twitterUser.ProfileBackgroundColor, twitterUser.ProfileLinkColor,
			twitterUser.ProfileSidebarBorderColor, twitterUser.ProfileSidebarFillColor)
		if err != nil {
			log.WithError(err).Error("error finding or creating user")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values[sessionUserIDKey] = user.ID
		err = session.Save(req, w)
		if err != nil {
			log.WithError(err).Error("failed to save the session")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, req, "/@"+user.Username, http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   sessionName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}
