package model

import (
	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// game theme

type GameTheme struct {
	Id     int64
	Title  string
	ImgUrl string
}

func (o *GameTheme) Copy() *GameTheme {
	if o == nil {
		return nil
	}
	return &GameTheme{
		Id:     o.Id,
		Title:  o.Title,
		ImgUrl: o.ImgUrl,
	}
}

func (o *GameTheme) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("id", int64(o.Id))
	enc.AddString("title", o.Title)
	enc.AddString("img-url", o.ImgUrl)
	return nil
}
