package oauth

import (
	"github.com/repetitive/oauth1"
	"os"
)

var AuthenticateEndpoint = oauth1.Endpoint{
	RequestTokenURL: "https://www.goodreads.com/oauth/request_token",
	AuthorizeURL:    "https://www.goodreads.com/oauth/authorize",
	AccessTokenURL:  "https://www.goodreads.com/oauth/access_token",
}

var AuthorizeEndpoint = oauth1.Endpoint{
	RequestTokenURL: "https://www.goodreads.com/oauth/request_token",
	AuthorizeURL:    "https://www.goodreads.com/oauth/authorize",
	AccessTokenURL:  "https://www.goodreads.com/oauth/access_token",
}

func GetGoodreadsOAuth1Config() *oauth1.Config {
	var callbackURL = os.Getenv("SERVICE_URL") + "/goodreads/callback"

	var oauth1Config = &oauth1.Config{
		ConsumerKey:    os.Getenv("GOODREADS_KEY"),
		ConsumerSecret: os.Getenv("GOODREADS_SECRET"),
		CallbackURL:    callbackURL,
		Endpoint:       AuthorizeEndpoint,
	}
	return oauth1Config
}
