package model

import (
	"fmt"
	"strings"

	"github.com/gre-ory/amnezic-go/internal/util"
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
	IsLocal  bool
	Mp3Url   Url
	ArtistId MusicArtistId
	AlbumId  MusicAlbumId

	// consolidated data
	Artist    *MusicArtist
	Album     *MusicAlbum
	Questions []*ThemeQuestion
}

func (o *Music) Validate(musicPathValidator, imagePathValidator PathValidator) error {
	if o == nil {
		return ErrConcurrentUpdate
	}
	if o.Name == "" {
		return ErrInvalidMusicName
	}
	if o.Mp3Url == "" || !o.Mp3Url.IsValid(musicPathValidator) {
		return ErrInvalidMusicUrl
	}
	if o.Artist == nil {
		return ErrMissingArtist
	} else {
		if err := o.Artist.Validate(imagePathValidator); err != nil {
			return err
		}
	}
	if o.Album != nil {
		if err := o.Album.Validate(imagePathValidator); err != nil {
			return err
		}
	}
	return nil
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

func (o *Music) GetMp3FileName() Url {
	parts := make([]string, 0)
	parts = append(parts, "music")
	if o.DeezerId != 0 {
		parts = append(parts, fmt.Sprintf("deezer-%d", o.DeezerId))
	}
	if o.Artist != nil && o.Artist.Name != "" {
		parts = append(parts, o.Artist.Name)
	}
	if o.Name != "" {
		parts = append(parts, o.Name)
	}
	return Url(strings.Join(util.Convert(parts, util.SanitizeAlphaLower), "_") + ".mp3")
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
	if o.Id != 0 {
		enc.AddInt64("id", int64(o.Id))
	}
	if o.DeezerId != 0 {
		enc.AddInt64("deezer-id", int64(o.DeezerId))
	}
	enc.AddString("name", o.Name)
	if o.Mp3Url != "" {
		enc.AddString("mp3-url", string(o.Mp3Url))
	}
	if o.Artist != nil {
		enc.AddObject("artist", o.Artist)
	}
	if o.Album != nil {
		enc.AddObject("album", o.Album)
	}
	return nil
}
