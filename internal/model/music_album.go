package model

import (
	"fmt"
	"strings"

	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// music album

type MusicAlbumId int64

type DeezerAlbumId int64

type MusicAlbum struct {
	Id       MusicAlbumId
	DeezerId DeezerAlbumId
	Name     string
	ImgUrl   Url

	// consolidated data
	Musics []*Music
}

func (o *MusicAlbum) Validate(imagePathValidator PathValidator) error {
	if o.Name == "" {
		return ErrInvalidAlbumName
	}
	if o.ImgUrl != "" && !o.ImgUrl.IsValid(imagePathValidator) {
		return ErrInvalidImageUrl
	}
	return nil
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

func (o *MusicAlbum) GetImageFileName() Url {
	parts := make([]string, 0)
	parts = append(parts, "album")
	if o.DeezerId != 0 {
		parts = append(parts, fmt.Sprintf("deezer-%d", o.DeezerId))
	}
	if o.Name != "" {
		parts = append(parts, o.Name)
	}
	return Url(strings.Join(util.Convert(parts, util.SanitizeAlphaLower), "_") + ".jpg")
}

func (o *MusicAlbum) MarshalLogObject(enc zapcore.ObjectEncoder) error {
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
