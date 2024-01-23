package model

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// session token

type SessionToken string

func (o SessionToken) String() string {
	return string(o)
}

// //////////////////////////////////////////////////
// session

type Session struct {
	Token      SessionToken
	UserId     UserId
	Expiration time.Time

	// consolidated data
	User *User
}

func (o *Session) IsExpired() bool {
	return o.Expiration.Before(time.Now())
}

func (o *Session) Copy() *Session {
	return &Session{
		Token:      o.Token,
		UserId:     o.UserId,
		Expiration: o.Expiration,
	}
}

func (o *Session) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("token", string(o.Token))
	if o.User != nil {
		enc.AddObject("user", o.User)
	} else {
		enc.AddInt64("user-id", o.UserId.ToInt64())
	}
	enc.AddTime("expiration", o.Expiration)
	return nil
}
