package model

import "go.uber.org/zap/zapcore"

// //////////////////////////////////////////////////
// music album

type MusicAlbumId int64

type DeezerAlbumId int64

type MusicAlbum struct {
	Id       MusicAlbumId
	DeezerId DeezerAlbumId
	Name     string
	ImgUrl   string
}

func (o *MusicAlbum) Copy() *MusicAlbum {
	if o == nil {
		return nil
	}
	return &MusicAlbum{
		Id:       o.Id,
		DeezerId: o.DeezerId,
		Name:     o.Name,
		ImgUrl:   o.ImgUrl,
	}
}

func (o *MusicAlbum) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("id", int64(o.Id))
	if o.DeezerId != 0 {
		enc.AddInt64("deezer-id", int64(o.DeezerId))
	}
	enc.AddString("name", o.Name)
	enc.AddString("img-url", o.ImgUrl)
	return nil
}
