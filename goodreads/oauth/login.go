package oauth

import (
	"fmt"
	"net/http"

	"github.com/dghubble/gologin"
	"github.com/repetitive/oauth1"
)

// LoginHandler handles OAuth1 login requests by obtaining a request token and
// secret (temporary credentials) and adding it to the ctx. If successful,
// handling delegates to the success handler, otherwise to the failure handler.
//
// Typically, the success handler is an AuthRedirectHandler or a handler which
// stores the request token secret.
func LoginHandler(config *oauth1.Config, success, failure http.Handler) http.Handler {
	if failure == nil {
		failure = gologin.DefaultFailureHandler
	}
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		requestToken, requestSecret, err := config.RequestToken()
		if err != nil {
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}
		fmt.Println("xxxxxx", requestToken, requestSecret)
		ctx = WithRequestToken(ctx, requestToken, requestSecret)
		success.ServeHTTP(w, req.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// AuthRedirectHandler reads the request token from the ctx and redirects
// to the authorization URL.
func AuthRedirectHandler(config *oauth1.Config, failure http.Handler) http.Handler {
	if failure == nil {
		failure = gologin.DefaultFailureHandler
	}
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		requestToken, _, err := RequestTokenFromContext(ctx)
		if err != nil {
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}
		authorizationURL, err := config.AuthorizationURL(requestToken)
		if err != nil {
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}
		fmt.Println("authorizaion url", authorizationURL.String())
		http.Redirect(w, req, authorizationURL.String(), http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
