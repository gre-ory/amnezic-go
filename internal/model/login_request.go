package model

import (
	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// login request

type LoginRequest struct {
	Name     string
	Password string
}

func (o *LoginRequest) Copy() *LoginRequest {
	return &LoginRequest{
		Name:     o.Name,
		Password: o.Password,
	}
}

func (o *LoginRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", o.Name)
	// enc.AddString("password", o.Password)
	return nil
}
