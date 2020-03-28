package twitter

import (
	"context"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
)

// unexported key type prevents collisions
type key int

const (
	userKey key = iota
	userAccessToken
)

// WithUser returns a copy of ctx that stores the Twitter User.
func WithUser(ctx context.Context, user *twitter.User, accessToken string) context.Context {
	ctx = context.WithValue(ctx, userKey, user)
	ctx = context.WithValue(ctx, userAccessToken, accessToken)
	return ctx
}

// UserFromContext returns the Twitter User from the ctx.
func UserFromContext(ctx context.Context) (*twitter.User, string, error) {
	user, ok := ctx.Value(userKey).(*twitter.User)
	if !ok {
		return nil, "", fmt.Errorf("twitter: Context missing Twitter User")
	}
	token, ok := ctx.Value(userAccessToken).(string)
	if !ok {
		return nil, "", fmt.Errorf("twitter: User access token is missing")
	}
	return user, token, nil
}
