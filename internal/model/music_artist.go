package model

import (
	"fmt"
	"strings"

	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// music artist

type MusicArtistId int64

type DeezerArtistId int64

type MusicArtist struct {
	Id       MusicArtistId
	DeezerId DeezerArtistId
	Name     string
	ImgUrl   Url

	// consolidated data
	Musics []*Music
}

func (o *MusicArtist) Validate(imagePathValidator PathValidator) error {
	if o.Name == "" {
		return ErrInvalidArtistName
	}
	if o.ImgUrl != "" && !o.ImgUrl.IsValid(imagePathValidator) {
		return ErrInvalidImageUrl
	}
	return nil
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

func (o *MusicArtist) GetImageFileName() Url {
	parts := make([]string, 0)
	parts = append(parts, "artist")
	if o.DeezerId != 0 {
		parts = append(parts, fmt.Sprintf("deezer-%d", o.DeezerId))
	}
	if o.Name != "" {
		parts = append(parts, o.Name)
	}
	return Url(strings.Join(util.Convert(parts, util.SanitizeAlphaLower), "_") + ".jpg")
}

func (o *MusicArtist) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if o.Id != 0 {
		enc.AddInt64("id", int64(o.Id))
	}
	if o.DeezerId != 0 {
		enc.AddInt64("deezer-id", int64(o.DeezerId))
	}
	enc.AddString("name", o.Name)
	if o.ImgUrl != "" {
		enc.AddString("img-url", string(o.ImgUrl))
	}
	return nil
}
