package routes

import (
	goodreadsOAuth "borges.ai/goodreads/oauth"
	"borges.ai/jobs"
	"borges.ai/models"
	"borges.ai/services"
	"borges.ai/utils"
	"github.com/dghubble/gologin"
	"github.com/repetitive/oauth1"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	REQUEST_TOKEN_SESSION_KEY  = "goodreads_request_token"
	REQUEST_SECRET_SESSION_KEY = "goodreads_request_secret"
)

func GoodreadsCallbackHandler(config *oauth1.Config, failure http.Handler) http.Handler {
	// goodreadsHandler
	success := goodreadsCallbackHandler(config, failure)
	return success
}

func goodreadsCallbackHandler(config *oauth1.Config, failure http.Handler) http.Handler {
	if failure == nil {
		failure = gologin.DefaultFailureHandler
	}
	fn := func(w http.ResponseWriter, req *http.Request) {
		defer utils.Duration(utils.Track("goodreadsCallbackHandler"))
		ctx := req.Context()

		session, err := Store.Get(req, sessionName)
		if err != nil {
			log.WithError(err).Error("failed to fetch session")
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}
		requestToken, ok := session.Values[REQUEST_TOKEN_SESSION_KEY].(string)
		if !ok {
			log.WithError(err).Error("request token not found")
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}
		requestSecret := session.Values[REQUEST_SECRET_SESSION_KEY].(string)
		if !ok {
			log.WithError(err).Error("request secret not found")
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}
		accessToken, accessSecret, err := config.AccessToken(requestToken, requestSecret, "")
		if err != nil {
			log.WithError(err).Error("failed to get access token")
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}

		goodreadsUserID, err := goodreadsOAuth.GetUserID(config, accessToken, accessSecret)
		if err != nil {
			log.WithError(err).Error("failed to get goodreads user")
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}
		log.WithField("goodreadsUserID", goodreadsUserID).Info("found goodreads user")
		sessionUser := GetUser(req)
		db, err := models.NewDB()
		if err != nil {
			log.WithError(err).Error("failed to connect to database")
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}
		defer db.Close()
		user, existed, err := services.FindOrCreateGoodreadsUser(db, sessionUser, goodreadsUserID, accessToken, accessSecret)
		if err != nil {
			log.WithError(err).Error("failed to save goodreads user")
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}
		session.Values[sessionUserIDKey] = user.ID
		err = session.Save(req, w)
		if err != nil {
			log.WithError(err).Error("failed to save the session")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.WithField("existed", existed).Info("is it a new goodreads user")
		// if user is new trigger job
		if !existed {
			_, err = jobs.SyncGoodreadsUserAndWait(user, false)
			if err != nil {
				log.WithError(err).Error("failed to start import job")
				ctx = gologin.WithError(ctx, err)
				failure.ServeHTTP(w, req.WithContext(ctx))
				return
			}
		}
		http.Redirect(w, req, "/@"+user.Username, http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

func GoodreadsLoginHandler(config *oauth1.Config, failure http.Handler) http.Handler {
	// oauth1.LoginHandler -> GoodreadsLoginSuccess ->oauth1.AuthRedirectHandler
	success := goodreadsOAuth.AuthRedirectHandler(config, failure)
	success = GoodreadsLoginSuccess(success, failure)
	return goodreadsOAuth.LoginHandler(config, success, failure)
}

func GoodreadsLoginSuccess(success, failure http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		defer utils.Duration(utils.Track("GoodreadsLoginSuccess"))
		ctx := req.Context()

		requestToken, requestSecret, err := goodreadsOAuth.RequestTokenFromContext(ctx)
		if err != nil {
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}
		// store tokens in the session
		session, err := Store.Get(req, sessionName)
		if err != nil {
			log.WithError(err).Error("failed to fetch session")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values[REQUEST_TOKEN_SESSION_KEY] = requestToken
		session.Values[REQUEST_SECRET_SESSION_KEY] = requestSecret
		err = session.Save(req, w)
		if err != nil {
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}
		success.ServeHTTP(w, req.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
