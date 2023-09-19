package model

import "go.uber.org/zap/zapcore"

// //////////////////////////////////////////////////
// theme questionfilter

type ThemeQuestionFilter struct {
	ThemeQuestionId ThemeQuestionId
	ThemeId         ThemeId
}

func (o *ThemeQuestionFilter) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if o.ThemeQuestionId != 0 {
		enc.AddInt64("theme-question-id", int64(o.ThemeQuestionId))
	}
	if o.ThemeId != 0 {
		enc.AddInt64("theme-id", int64(o.ThemeId))
	}
	return nil
}

func (o *ThemeQuestionFilter) IsMatching(candidate *ThemeQuestion) bool {
	if o.ThemeQuestionId != 0 {
		if candidate.Id == o.ThemeQuestionId {
			return true
		}
	}
	if o.ThemeId != 0 {
		if candidate.ThemeId == o.ThemeId {
			return true
		}
	}
	return false
}
