package model

import "go.uber.org/zap/zapcore"

// //////////////////////////////////////////////////
// theme question

type ThemeQuestionId int64

type ThemeQuestion struct {
	Id      ThemeQuestionId
	ThemeId ThemeId
	MusicId MusicId
	Text    string
	Hint    string

	// consolidated data
	Theme *Theme
	Music *Music
}

func (o *ThemeQuestion) Validate() error {
	if o.ThemeId == 0 {
		return ErrInvalidThemeId
	}
	if o.MusicId == 0 {
		return ErrInvalidMusicId
	}
	if o.Text == "" {
		return ErrInvalidThemeQuestion
	}
	return nil
}

func (o *ThemeQuestion) Copy() *ThemeQuestion {
	return &ThemeQuestion{
		Id:      o.Id,
		ThemeId: o.ThemeId,
		MusicId: o.MusicId,
		Text:    o.Text,
		Hint:    o.Hint,
	}
}

func (o *ThemeQuestion) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("id", int64(o.Id))
	enc.AddInt64("theme-id", int64(o.ThemeId))
	enc.AddInt64("music-id", int64(o.MusicId))
	enc.AddString("text", o.Text)
	if o.Hint != "" {
		enc.AddString("hint", o.Hint)
	}
	if o.Music != nil {
		enc.AddObject("music", o.Music)
	}
	return nil
}
