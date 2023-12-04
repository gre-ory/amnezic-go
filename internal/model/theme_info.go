package model

import (
	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// theme info

type ThemeInfo struct {
	Id         ThemeId
	Title      string
	ImgUrl     string
	Labels     map[string]string
	NbQuestion int
}

func (o *ThemeInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("id", int64(o.Id))
	enc.AddString("title", o.Title)
	enc.AddString("img-url", o.ImgUrl)
	enc.AddObject("labels", zapcore.ObjectMarshalerFunc(o.MarshalLogLabels))
	enc.AddInt("nb-questions", o.NbQuestion)
	return nil
}

func (o *ThemeInfo) MarshalLogLabels(enc zapcore.ObjectEncoder) error {
	for key, value := range o.Labels {
		enc.AddString(key, value)
	}
	return nil
}
