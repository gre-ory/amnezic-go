package model

import (
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// game settings

type GameSettings struct {
	Seed       int64
	NbQuestion int
	NbAnswer   int
	NbPlayer   int
	Sources    []Source
}

func (o *GameSettings) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("seed", o.Seed)
	enc.AddInt("nb-question", o.NbQuestion)
	enc.AddInt("nb-answer", o.NbAnswer)
	enc.AddInt("nb-player", o.NbPlayer)
	enc.AddString("sources", util.Join(o.Sources, Source.String))
	return nil
}

// //////////////////////////////////////////////////
// validate

const (
	MinNbPlayer = 2
	MaxNbPlayer = 99

	MinNbQuestion = 1
	MaxNbQuestion = 999

	MinNbAnswer = 2
	MaxNbAnswer = 99
)

func (o *GameSettings) Validate() error {
	if o.NbPlayer < MinNbPlayer || o.NbPlayer > MaxNbPlayer {
		return ErrInvalidNbPlayer
	}
	if o.NbQuestion < MinNbQuestion || o.NbQuestion > MaxNbQuestion {
		return ErrInvalidNbQuestion
	}
	if o.NbAnswer < MinNbAnswer || o.NbAnswer > MaxNbAnswer {
		return ErrInvalidNbAnswer
	}
	if len(o.Sources) == 0 {
		return ErrMissingSource
	}
	return nil
}
