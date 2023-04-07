package model

import "go.uber.org/zap/zapcore"

// //////////////////////////////////////////////////
// game answer

type GameAnswer struct {
	Id      GameAnswerId
	Text    string
	Hint    string
	Correct bool
}

func (o *GameAnswer) Copy() *GameAnswer {
	if o == nil {
		return nil
	}
	return &GameAnswer{
		Id:      o.Id,
		Text:    o.Text,
		Hint:    o.Hint,
		Correct: o.Correct,
	}
}

func (o *GameAnswer) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("id", int64(o.Id))
	enc.AddString("text", o.Text)
	if o.Hint != "" {
		enc.AddString("hint", o.Hint)
	}
	if o.Correct {
		enc.AddBool("correct", o.Correct)
	}
	return nil
}
