package model

import (
	"fmt"

	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// user id

type UserId int64

func (i UserId) String() string {
	return fmt.Sprintf("%d", i)
}

func (i UserId) ToInt64() int64 {
	return int64(i)
}

func ToUserId(value string) UserId {
	return UserId(util.StrToInt64(value))
}

// //////////////////////////////////////////////////
// user

type User struct {
	Id          UserId
	Name        string
	OldPassword string
	Password    string
	Hash        string
	Permissions []Permission
}

func (o *User) HasPermission(permission Permission) bool {
	for _, p := range o.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func (o *User) AddPermission(permission Permission) {
	if !o.HasPermission(permission) {
		o.Permissions = append(o.Permissions, permission)
	}
}

func (o *User) RemovePermission(permission Permission) {
	o.Permissions = util.Filter(o.Permissions, func(p Permission) bool {
		return p != permission
	})
}

func (o *User) Copy() *User {
	return &User{
		Id:          o.Id,
		Name:        o.Name,
		Password:    o.Password,
		Hash:        o.Hash,
		Permissions: o.Permissions,
	}
}

func (o *User) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("id", int64(o.Id))
	enc.AddString("name", o.Name)
	// enc.AddString("password", o.Password)
	enc.AddString("hash", o.Hash)
	enc.AddArray("permissions", zapcore.ArrayMarshalerFunc(o.MarshalLogPermissions))
	return nil
}

func (o *User) MarshalLogPermissions(enc zapcore.ArrayEncoder) error {
	for _, permission := range o.Permissions {
		enc.AppendString(permission.String())
	}
	return nil
}
