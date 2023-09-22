package model

import "go.uber.org/zap/zapcore"

// //////////////////////////////////////////////////
// theme questionfilter

type ThemeQuestionFilter struct {
	ThemeQuestionId ThemeQuestionId
	ThemeId         ThemeId
	MusicId         MusicId
	Random          bool
	Limit           int
}

func (o *ThemeQuestionFilter) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if o.ThemeQuestionId != 0 {
		enc.AddInt64("theme-question-id", int64(o.ThemeQuestionId))
	}
	if o.ThemeId != 0 {
		enc.AddInt64("theme-id", int64(o.ThemeId))
	}
	if o.MusicId != 0 {
		enc.AddInt64("music-id", int64(o.MusicId))
	}
	if o.Random {
		enc.AddBool("random", o.Random)
	}
	if o.Limit != 0 {
		enc.AddInt("limit", o.Limit)
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
	if o.MusicId != 0 {
		if candidate.MusicId == o.MusicId {
			return true
		}
	}
	return false
}
