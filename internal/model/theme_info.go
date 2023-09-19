package model

import "go.uber.org/zap/zapcore"

// //////////////////////////////////////////////////
// theme info

type ThemeInfo struct {
	Id         ThemeId
	Title      string
	ImgUrl     string
	NbQuestion int
}

func (o *ThemeInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("id", int64(o.Id))
	enc.AddString("title", o.Title)
	enc.AddString("img-url", o.ImgUrl)
	enc.AddInt("nb-questions", o.NbQuestion)
	return nil
}
