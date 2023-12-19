package model

import (
	"fmt"
	"net/url"
	"strings"

	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// search music request

type SearchMusicRequest struct {
	query  string
	artist string
	album  string
	track  string
	label  string
	limit  int
	strict bool
}

func (o *SearchMusicRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
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

func NewSearchMusicRequest() *SearchMusicRequest {
	return &SearchMusicRequest{}
}

func (r *SearchMusicRequest) WithQuery(query string) *SearchMusicRequest {
	r.query = query
	return r
}

func (r *SearchMusicRequest) WithArtist(artist string) *SearchMusicRequest {
	r.artist = artist
	return r
}

func (r *SearchMusicRequest) WithAlbum(album string) *SearchMusicRequest {
	r.album = album
	return r
}

func (r *SearchMusicRequest) WithTrack(track string) *SearchMusicRequest {
	r.track = track
	return r
}

func (r *SearchMusicRequest) WithLabel(label string) *SearchMusicRequest {
	r.label = label
	return r
}

func (r *SearchMusicRequest) WithLimit(limit int) *SearchMusicRequest {
	r.limit = limit
	return r
}

func (r *SearchMusicRequest) WithStrict(strict bool) *SearchMusicRequest {
	r.strict = strict
	return r
}

// //////////////////////////////////////////////////
// compute query

func (r *SearchMusicRequest) ComputeQuery() string {
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

func (r *SearchMusicRequest) ComputeParameters() string {
	parts := []string{}
	query := r.ComputeQuery()
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
