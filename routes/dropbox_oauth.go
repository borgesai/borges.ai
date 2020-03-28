package routes

import (
	"borges.ai/data"
	"github.com/repetitive/dropbox"
	"borges.ai/models"
	"github.com/dghubble/gologin"
	oauth2Login "github.com/dghubble/gologin/oauth2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
)

func GetDropboxOauth2Config() *oauth2.Config {
	redirectURL := "https://borges.ai/dropbox/callback"
	var oauth2Config = &oauth2.Config{
		ClientID:     "7ze4ivemwaf7wdl",
		ClientSecret: "tco62wcha9acop8",
		RedirectURL:  redirectURL,
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://api.dropbox.com/oauth2/token",
			AuthURL:  "https://www.dropbox.com/oauth2/authorize",
		},
	}
	return oauth2Config

}

func DropboxCallbackHandler()http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		log.Info("issue new session")
				ctx := req.Context()
				token, err := oauth2Login.TokenFromContext(ctx)
				if err != nil {
					log.WithError(err).Error("failed to get token")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				accessToken := token.AccessToken
				account, err := dropbox.GetUser(accessToken)

		sessionUser := GetUser(req)
		db, err := models.NewDB()
		if err != nil {
			log.WithError(err).Error("failed to connect to database")
			ctx = gologin.WithError(ctx, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		if sessionUser.ID == 0 {
			log.WithError(err).Error("cannot do this for gues")
			ctx = gologin.WithError(ctx, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// TODO save email too 		account.Email
		if sessionUser.Email != "" {
			_, err = data.UpdateUserWithDropboxInfo(db, sessionUser, account.AccountId, accessToken)
			if err != nil {
				log.WithError(err).Error("failed to save dropbox data")
				ctx = gologin.WithError(ctx, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			_, err = data.UpdateUserWithDropboxInfoAndEmail(db, sessionUser, account.AccountId, accessToken, account.Email)
			if err != nil {
				log.WithError(err).Error("failed to save dropbox data")
				ctx = gologin.WithError(ctx, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		log.WithField("acc", account).WithField("ses", sessionUser).Info("dropbox done")
		http.Redirect(w, req, "/@"+sessionUser.Username, http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
// issueSession issues a cookie session after successful Dropbox login
//func IssueSession() http.Handler {
//	fn := func(w http.ResponseWriter, req *http.Request) {
//		log.Info("issue new session")
//		db, err := models.NewDB()
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		defer db.Close()
//		ctx := req.Context()
//		token, err := oauth2Login.TokenFromContext(ctx)
//		if err != nil {
//			log.WithError(err).Error("failed to get token")
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		accessToken := token.AccessToken
//		account, err := dropbox.GetUser(accessToken)
//
//		if err != nil {
//			log.WithError(err).Error("failed to fetch user from dropbox")
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		log.WithField("full account", account).Info("account")
//		// implement a success handler to issue some form of session
//		session, err := Store.Get(req, sessionName)
//		if err != nil {
//			log.WithError(err).Error("failed to get session")
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		user, err := services.FindOrCreateUser(db, accessToken, account.AccountId, account.Email, account.ProfilePhotoUrl, account.Name.GivenName, account.Name.Surname, account.Name.AbbreviatedName, account.Name.DisplayName, account.Country)
//		if err != nil {
//			log.WithError(err).Error("failed to fetch or create user")
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		session.Values[sessionUserIDKey] = user.ID
//
//		session.Save(req, w)
//		http.Redirect(w, req, "/", http.StatusFound)
//	}
//	return http.HandlerFunc(fn)
//}
