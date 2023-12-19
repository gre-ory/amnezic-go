package model

import (
	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// playlist

type DeezerPlaylistId int64

type Playlist struct {
	DeezerId    DeezerPlaylistId
	Name        string
	Public      bool
	PlaylistUrl string
	ImgUrl      string
	NbMusics    int
	User        string

	// consolidated data
	Musics []*Music
}

func (o *Playlist) Copy() *Playlist {
	if o == nil {
		return nil
	}
	return &Playlist{
		DeezerId:    o.DeezerId,
		Name:        o.Name,
		Public:      o.Public,
		PlaylistUrl: o.PlaylistUrl,
		ImgUrl:      o.ImgUrl,
		NbMusics:    o.NbMusics,
		User:        o.User,
		Musics:      o.Musics,
	}
}

func (o *Playlist) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("deezer-id", int64(o.DeezerId))
	enc.AddString("name", o.Name)
	enc.AddBool("public", o.Public)
	enc.AddString("playlist-url", o.PlaylistUrl)
	enc.AddString("picture-url", o.ImgUrl)
	enc.AddInt("nb-musics", o.NbMusics)
	if o.User != "" {
		enc.AddString("user", o.User)
	}
	return nil
}
