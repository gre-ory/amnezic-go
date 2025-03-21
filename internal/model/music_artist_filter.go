package model

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// music artist filter

type MusicArtistFilter struct {
	Name  string
	Limit int
}

func (o *MusicArtistFilter) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if o.Name != "" {
		enc.AddString("name", o.Name)
	}
	if o.Limit != 0 {
		enc.AddInt("limit", o.Limit)
	}
	return nil
}

func (o *MusicArtistFilter) IsMatching(count int, candidate *MusicArtist) bool {
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

func NewMusicArtistFilter() *MusicArtistFilter {
	return &MusicArtistFilter{}
}

func (r *MusicArtistFilter) WithName(name string) *MusicArtistFilter {
	r.Name = name
	return r
}

func (r *MusicArtistFilter) WithLimit(limit int) *MusicArtistFilter {
	r.Limit = limit
	return r
}
