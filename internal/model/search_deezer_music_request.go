package model

import (
	"fmt"
	"net/url"
	"strings"

	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// search deezer music request

type SearchDeezerMusicRequest struct {
	query  string
	artist string
	album  string
	track  string
	label  string
	limit  int
	strict bool
}

func (o *SearchDeezerMusicRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if o.query != "" {
		enc.AddString("query", o.query)
	}
	if o.artist != "" {
		enc.AddString("artist", o.artist)
	}
	if o.album != "" {
		enc.AddString("album", o.album)
	}
	if o.track != "" {
		enc.AddString("track", o.track)
	}
	if o.label != "" {
		enc.AddString("label", o.label)
	}
	if o.limit != 0 {
		enc.AddInt("limit", o.limit)
	}
	if o.strict {
		enc.AddBool("strict", o.strict)
	}
	return nil
}

// //////////////////////////////////////////////////
// builder

func NewSearchDeezerMusicRequest() *SearchDeezerMusicRequest {
	return &SearchDeezerMusicRequest{}
}

func (r *SearchDeezerMusicRequest) WithQuery(query string) *SearchDeezerMusicRequest {
	r.query = query
	return r
}

func (r *SearchDeezerMusicRequest) WithArtist(artist string) *SearchDeezerMusicRequest {
	r.artist = artist
	return r
}

func (r *SearchDeezerMusicRequest) WithAlbum(album string) *SearchDeezerMusicRequest {
	r.album = album
	return r
}

func (r *SearchDeezerMusicRequest) WithTrack(track string) *SearchDeezerMusicRequest {
	r.track = track
	return r
}

func (r *SearchDeezerMusicRequest) WithLabel(label string) *SearchDeezerMusicRequest {
	r.label = label
	return r
}

func (r *SearchDeezerMusicRequest) WithLimit(limit int) *SearchDeezerMusicRequest {
	r.limit = limit
	return r
}

func (r *SearchDeezerMusicRequest) WithStrict(strict bool) *SearchDeezerMusicRequest {
	r.strict = strict
	return r
}

// //////////////////////////////////////////////////
// compute query

func (r *SearchDeezerMusicRequest) ComputeDeezerQuery() string {
	parts := []string{}
	if r.query != "" {
		parts = append(parts, r.query)
	}
	if r.album != "" {
		parts = append(parts, fmt.Sprintf("album:\"%s\"", r.album))
	}
	if r.artist != "" {
		parts = append(parts, fmt.Sprintf("artist:\"%s\"", r.artist))
	}
	if r.label != "" {
		parts = append(parts, fmt.Sprintf("label:\"%s\"", r.label))
	}
	if r.track != "" {
		parts = append(parts, fmt.Sprintf("track:\"%s\"", r.track))
	}
	return strings.Join(parts, " ")
}

func (r *SearchDeezerMusicRequest) ComputeDeezerParameters() string {
	parts := []string{}
	query := r.ComputeDeezerQuery()
	if query != "" {
		parts = append(parts, fmt.Sprintf("q=%s", url.QueryEscape(query)))
	}
	if r.limit > 0 {
		parts = append(parts, fmt.Sprintf("limit=%d", r.limit))
	}
	if r.strict {
		parts = append(parts, "strict=on")
	}
	if len(parts) != 0 {
		return "?" + strings.Join(parts, "&")
	}
	return ""
}
