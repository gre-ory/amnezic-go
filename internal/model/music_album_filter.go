package model

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// music album filter

type MusicAlbumFilter struct {
	Name  string
	Limit int
}

func (o *MusicAlbumFilter) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if o.Name != "" {
		enc.AddString("name", o.Name)
	}
	if o.Limit != 0 {
		enc.AddInt("limit", o.Limit)
	}
	return nil
}

func (o *MusicAlbumFilter) IsMatching(count int, candidate *MusicAlbum) bool {
	if o.Limit != 0 && count > o.Limit {
		return false
	}
	if o.Name != "" && !strings.Contains(candidate.Name, o.Name) {
		return false
	}
	return true
}

// //////////////////////////////////////////////////
// builder

func NewMusicAlbumFilter() *MusicAlbumFilter {
	return &MusicAlbumFilter{}
}

func (r *MusicAlbumFilter) WithName(name string) *MusicAlbumFilter {
	r.Name = name
	return r
}

func (r *MusicAlbumFilter) WithLimit(limit int) *MusicAlbumFilter {
	r.Limit = limit
	return r
}
