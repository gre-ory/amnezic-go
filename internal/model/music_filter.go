package model

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// music filter

type MusicFilter struct {
	Name     string
	ArtistId MusicArtistId
	AlbumId  MusicAlbumId
	Limit    int
}

func (o *MusicFilter) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if o.Name != "" {
		enc.AddString("name", o.Name)
	}
	if o.ArtistId != 0 {
		enc.AddInt64("artist-id", int64(o.ArtistId))
	}
	if o.AlbumId != 0 {
		enc.AddInt64("album-id", int64(o.AlbumId))
	}
	if o.Limit != 0 {
		enc.AddInt("limit", o.Limit)
	}
	return nil
}

func (o *MusicFilter) IsMatching(count int, candidate *Music) bool {
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

func NewMusicFilter() *MusicFilter {
	return &MusicFilter{}
}

func (r *MusicFilter) WithName(name string) *MusicFilter {
	r.Name = name
	return r
}

func (r *MusicFilter) WithArtistId(artistId MusicArtistId) *MusicFilter {
	r.ArtistId = artistId
	return r
}

func (r *MusicFilter) WithAlbumId(albumId MusicAlbumId) *MusicFilter {
	r.AlbumId = albumId
	return r
}

func (r *MusicFilter) WithLimit(limit int) *MusicFilter {
	r.Limit = limit
	return r
}
