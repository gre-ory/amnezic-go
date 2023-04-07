package model

import (
	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// music

type MusicId int64

type DeezerMusicId int64

type Music struct {
	Id            MusicId
	DeezerMusicId DeezerMusicId
	Name          string
	Mp3Url        string
	ArtistId      MusicArtistId
	AlbumId       MusicAlbumId
	GenreId       MusicGenreId

	// consolidated data
	Artist *MusicArtist
	Album  *MusicAlbum
	Genre  *MusicGenre
}

func (o *Music) Copy() *Music {
	if o == nil {
		return nil
	}
	return &Music{
		Id:            o.Id,
		DeezerMusicId: o.DeezerMusicId,
		Name:          o.Name,
		Mp3Url:        o.Mp3Url,
		ArtistId:      o.ArtistId,
		AlbumId:       o.AlbumId,
		GenreId:       o.GenreId,
	}
}

func (o *Music) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("id", int64(o.Id))
	if o.DeezerMusicId != 0 {
		enc.AddInt64("deezer-music-id", int64(o.DeezerMusicId))
	}
	enc.AddString("name", o.Name)
	enc.AddString("mp3-url", o.Mp3Url)
	if o.Artist != nil {
		enc.AddObject("artist", o.Artist)
	}
	if o.Album != nil {
		enc.AddObject("album", o.Album)
	}
	if o.Genre != nil {
		enc.AddObject("genre", o.Genre)
	}
	return nil
}
