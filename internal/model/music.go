package model

import (
	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// music

type MusicId int64

type DeezerMusicId int64

type Music struct {
	Id       MusicId
	DeezerId DeezerMusicId
	Name     string
	Mp3Url   string
	ArtistId MusicArtistId
	AlbumId  MusicAlbumId

	// consolidated data
	Artist    *MusicArtist
	Album     *MusicAlbum
	Questions []*ThemeQuestion
}

func (o *Music) Copy() *Music {
	if o == nil {
		return nil
	}
	return &Music{
		Id:       o.Id,
		DeezerId: o.DeezerId,
		Name:     o.Name,
		Mp3Url:   o.Mp3Url,
		ArtistId: o.ArtistId,
		AlbumId:  o.AlbumId,
	}
}

func (o *Music) GetDefaultAnswerText() string {
	if o.Artist != nil {
		return o.Artist.Name
	}
	return ""
}

func (o *Music) GetDefaultAnswerHint() string {
	hint := o.Name
	if hint == "" && o.Album != nil && o.Album.Name != "" {
		hint = o.Album.Name
	}
	return hint
}

func (o *Music) ToGameAnswer(correct bool) *GameAnswer {
	return &GameAnswer{
		Text:    o.GetDefaultAnswerText(),
		Hint:    o.GetDefaultAnswerHint(),
		Correct: correct,
	}
}

func (o *Music) ToThemeQuestion(themeId ThemeId) *ThemeQuestion {
	return &ThemeQuestion{
		ThemeId: themeId,
		MusicId: o.Id,
		Text:    o.GetDefaultAnswerText(),
		Hint:    o.GetDefaultAnswerHint(),
	}
}

func (o *Music) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("id", int64(o.Id))
	if o.DeezerId != 0 {
		enc.AddInt64("deezer-id", int64(o.DeezerId))
	}
	enc.AddString("name", o.Name)
	enc.AddString("mp3-url", o.Mp3Url)
	if o.Artist != nil {
		enc.AddObject("artist", o.Artist)
	}
	if o.Album != nil {
		enc.AddObject("album", o.Album)
	}
	return nil
}
