package model

import (
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// game question

type GameQuestion struct {
	Id      GameQuestionId
	Theme   *GameTheme
	Music   *Music
	Answers []*GameAnswer
}

func (o *GameQuestion) Copy() *GameQuestion {
	if o == nil {
		return nil
	}
	return &GameQuestion{
		Id:      o.Id,
		Theme:   o.Theme.Copy(),
		Music:   o.Music.Copy(),
		Answers: util.Convert(o.Answers, (*GameAnswer).Copy),
	}
}

func (o *GameQuestion) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("id", int64(o.Id))
	enc.AddObject("theme", o.Theme)
	enc.AddObject("music", o.Music)
	enc.AddInt("nb-answers", len(o.Answers))
	return nil
}
