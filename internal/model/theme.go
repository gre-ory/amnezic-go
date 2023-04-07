package model

import "go.uber.org/zap/zapcore"

// //////////////////////////////////////////////////
// theme

type ThemeId int64

type Theme struct {
	Id     ThemeId
	Title  string
	ImgUrl string

	// consolidated data
	Questions []*ThemeQuestion
}

func (o *Theme) Copy() *Theme {
	return &Theme{
		Id:     o.Id,
		Title:  o.Title,
		ImgUrl: o.ImgUrl,
	}
}

func (o *Theme) Equal(other *Theme) bool {
	if o == nil {
		return other == nil
	}
	if other == nil {
		return false
	}
	return (o.Id == other.Id) &&
		(o.Title == other.Title) &&
		(o.ImgUrl == other.ImgUrl)
}

func (o *Theme) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("id", int64(o.Id))
	enc.AddString("title", o.Title)
	enc.AddString("img-url", o.ImgUrl)
	if o.Questions != nil {
		enc.AddInt("nb-questions", len(o.Questions))
	}
	return nil
}
