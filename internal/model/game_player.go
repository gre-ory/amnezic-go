package model

import "go.uber.org/zap/zapcore"

// //////////////////////////////////////////////////
// game player

type GamePlayer struct {
	Id     GamePlayerId
	Name   string
	Active bool
	Score  int
}

func (o *GamePlayer) Copy() *GamePlayer {
	if o == nil {
		return nil
	}
	return &GamePlayer{
		Id:     o.Id,
		Name:   o.Name,
		Active: o.Active,
		Score:  o.Score,
	}
}

func (o *GamePlayer) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("id", int64(o.Id))
	enc.AddString("name", o.Name)
	enc.AddInt("score", o.Score)
	enc.AddBool("active", o.Active)
	return nil
}
