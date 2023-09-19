package model

import "go.uber.org/zap/zapcore"

// //////////////////////////////////////////////////
// music artist

type MusicArtistId int64

type DeezerArtistId int64

type MusicArtist struct {
	Id       MusicArtistId
	DeezerId DeezerArtistId
	Name     string
	ImgUrl   string
}

func (o *MusicArtist) Copy() *MusicArtist {
	if o == nil {
		return nil
	}
	return &MusicArtist{
		Id:       o.Id,
		DeezerId: o.DeezerId,
		Name:     o.Name,
		ImgUrl:   o.ImgUrl,
	}
}

func (o *MusicArtist) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("id", int64(o.Id))
	if o.DeezerId != 0 {
		enc.AddInt64("deezer-id", int64(o.DeezerId))
	}
	enc.AddString("name", o.Name)
	enc.AddString("img-url", o.ImgUrl)
	return nil
}
