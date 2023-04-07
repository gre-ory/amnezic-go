package model

import "go.uber.org/zap/zapcore"

// //////////////////////////////////////////////////
// music genre

type MusicGenreId int64

type DeezerGenreId int64

type MusicGenre struct {
	Id            MusicGenreId
	DeezerGenreId DeezerGenreId
	Name          string
	ImgUrl        string
}

func (o *MusicGenre) Copy() *MusicGenre {
	if o == nil {
		return nil
	}
	return &MusicGenre{
		Id:            o.Id,
		DeezerGenreId: o.DeezerGenreId,
		Name:          o.Name,
		ImgUrl:        o.ImgUrl,
	}
}

func (o *MusicGenre) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("id", int64(o.Id))
	if o.DeezerGenreId != 0 {
		enc.AddInt64("deezer-genre-id", int64(o.DeezerGenreId))
	}
	enc.AddString("name", o.Name)
	enc.AddString("img-url", o.ImgUrl)
	return nil
}
