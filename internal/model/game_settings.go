package model

import (
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// game settings

type GameSettings struct {
	Seed             int64
	NbQuestion       int
	NbAnswer         int
	NbPlayer         int
	Sources          []Source
	ThemeIds         []ThemeId
	DeezerPlaylistId DeezerPlaylistId
}

func (o *GameSettings) UseDeezerPlaylist() bool {
	if o.DeezerPlaylistId == 0 {
		return false
	}
	_, ok := util.FindIf(o.Sources, Source.IsDeezer)
	return ok
}

func (o *GameSettings) UseStore() bool {
	_, ok := util.FindIf(o.Sources, Source.IsStore)
	return ok
}

func (o *GameSettings) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("seed", o.Seed)
	enc.AddInt("nb-question", o.NbQuestion)
	enc.AddInt("nb-answer", o.NbAnswer)
	enc.AddInt("nb-player", o.NbPlayer)
	if len(o.Sources) > 0 {
		enc.AddString("sources", util.Join(o.Sources, Source.String))
	}
	if len(o.ThemeIds) > 0 {
		enc.AddString("theme-ids", util.Join(o.ThemeIds, ThemeId.String))
	}
	if o.DeezerPlaylistId != 0 {
		enc.AddInt64("deezer-playlist-id", int64(o.DeezerPlaylistId))
	}
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
