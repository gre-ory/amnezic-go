package model

import "go.uber.org/zap/zapcore"

// //////////////////////////////////////////////////
// theme questionfilter

type ThemeQuestionFilter struct {
	ThemeQuestionId ThemeQuestionId
	ThemeIds        []ThemeId
	MusicId         MusicId
	Random          bool
	Limit           int
}

func (o *ThemeQuestionFilter) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if o.ThemeQuestionId != 0 {
		enc.AddInt64("theme-question-id", int64(o.ThemeQuestionId))
	}
	if len(o.ThemeIds) != 0 {
		enc.AddArray("theme-ids", o.MarshalLogThemeIds())
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

func (o *ThemeQuestionFilter) MarshalLogThemeIds() zapcore.ArrayMarshalerFunc {
	return func(enc zapcore.ArrayEncoder) error {
		for _, v := range o.ThemeIds {
			enc.AppendInt64(int64(v))
		}
		return nil
	}
}

func (o *ThemeQuestionFilter) IsMatching(candidate *ThemeQuestion) bool {
	if o.ThemeQuestionId != 0 {
		if candidate.Id != o.ThemeQuestionId {
			return false
		}
	}
	if len(o.ThemeIds) != 0 {
		atLeastOne := false
		for _, themeId := range o.ThemeIds {
			if candidate.ThemeId == themeId {
				atLeastOne = true
			}
		}
		if !atLeastOne {
			return false
		}
	}
	if o.MusicId != 0 {
		if candidate.MusicId != o.MusicId {
			return false
		}
	}
	return true
}
