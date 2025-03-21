package model

import (
	"fmt"
	"net/url"
	"strings"

	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// search deezer playlist request

type SearchDeezerPlaylistRequest struct {
	deezerPlaylistId DeezerPlaylistId
	query            string
	limit            int
	strict           bool
}

func (o *SearchDeezerPlaylistRequest) GetDeezerPlaylistId() DeezerPlaylistId {
	return o.deezerPlaylistId
}

func (o *SearchDeezerPlaylistRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if o.deezerPlaylistId != 0 {
		enc.AddInt64("deezer-playlist-id", int64(o.deezerPlaylistId))
	}
	if o.query != "" {
		enc.AddString("query", o.query)
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

func NewSearchDeezerPlaylistRequest() *SearchDeezerPlaylistRequest {
	return &SearchDeezerPlaylistRequest{}
}

func (r *SearchDeezerPlaylistRequest) WithQuery(query string) *SearchDeezerPlaylistRequest {
	r.query = query
	return r
}

func (r *SearchDeezerPlaylistRequest) WithDeezerPlaylistId(deezerPlaylistId DeezerPlaylistId) *SearchDeezerPlaylistRequest {
	r.deezerPlaylistId = deezerPlaylistId
	return r
}

func (r *SearchDeezerPlaylistRequest) WithLimit(limit int) *SearchDeezerPlaylistRequest {
	r.limit = limit
	return r
}

func (r *SearchDeezerPlaylistRequest) WithStrict(strict bool) *SearchDeezerPlaylistRequest {
	r.strict = strict
	return r
}

// //////////////////////////////////////////////////
// compute query

func (r *SearchDeezerPlaylistRequest) ComputeQuery() string {
	parts := []string{}
	if r.query != "" {
		parts = append(parts, r.query)
	}
	return strings.Join(parts, " ")
}

func (r *SearchDeezerPlaylistRequest) ComputeParameters() string {
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
