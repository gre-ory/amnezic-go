package model

import (
	"context"
)

// //////////////////////////////////////////////////
// session

type sessionKeyType int

var sessionKey sessionKeyType

func WithSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

func GetSession(ctx context.Context) *Session {
	if session, ok := ctx.Value(sessionKey).(*Session); ok {
		return session
	}
	return nil
}

func GetSessionToken(ctx context.Context) SessionToken {
	if session := GetSession(ctx); session != nil {
		return session.Token
	}
	return ""
}

func GetCurrentUser(ctx context.Context) *User {
	if session := GetSession(ctx); session != nil {
		return session.User
	}
	return nil
}

func IsCurrentUser(ctx context.Context, id UserId) bool {
	if currentUser := GetCurrentUser(ctx); currentUser != nil {
		return currentUser.Id == id
	}
	return false
}
