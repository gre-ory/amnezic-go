package model

import "go.uber.org/zap/zapcore"

// //////////////////////////////////////////////////
// theme filter

type UserFilter struct {
	UserId UserId
}

func (o *UserFilter) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if o.UserId != 0 {
		enc.AddInt64("user-id", int64(o.UserId))
	}
	return nil
}

func (o *UserFilter) IsMatching(candidate *User) bool {
	if o == nil {
		return true
	}
	if o.UserId != 0 {
		if candidate.Id == o.UserId {
			return true
		}
	}
	return false
}
