package model

import (
	"context"
)

// //////////////////////////////////////////////////
// session token

type sessionTokenKeyType int

var sessionTokenKey sessionTokenKeyType

func WithSessionToken(ctx context.Context, token SessionToken) context.Context {
	return context.WithValue(ctx, sessionTokenKey, token)
}

func GetSessionToken(ctx context.Context) SessionToken {
	if token, ok := ctx.Value(sessionTokenKey).(SessionToken); ok {
		return token
	}
	return ""
}

// //////////////////////////////////////////////////
// user

type userKeyType int

var userKey userKeyType

func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func GetUser(ctx context.Context) *User {
	if user, ok := ctx.Value(userKey).(*User); ok {
		return user
	}
	return nil
}

func IsCurrentUser(ctx context.Context, id UserId) bool {
	currentUser := GetUser(ctx)
	return currentUser != nil && currentUser.Id == id
}
