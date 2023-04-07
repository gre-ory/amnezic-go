package model

import "go.uber.org/zap/zapcore"

// //////////////////////////////////////////////////
// theme filter

type ThemeFilter struct {
	ThemeId ThemeId
}

func (o *ThemeFilter) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if o.ThemeId != 0 {
		enc.AddInt64("theme-id", int64(o.ThemeId))
	}
	return nil
}

func (o *ThemeFilter) IsMatching(candidate *Theme) bool {
	if o.ThemeId != 0 {
		if candidate.Id == o.ThemeId {
			return true
		}
	}
	return false
}
