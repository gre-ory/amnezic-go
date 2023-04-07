package store

import (
	"context"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// music album store

type MusicAlbumStore interface {
	Create(ctx context.Context, music *model.MusicAlbum) (*model.MusicAlbum, error)
	Retrieve(ctx context.Context, id model.MusicAlbumId) (*model.MusicAlbum, error)
	RetrieveByDeezerId(ctx context.Context, deezerId model.DeezerAlbumId) (*model.MusicAlbum, error)
	Update(ctx context.Context, music *model.MusicAlbum) (*model.MusicAlbum, error)
	Delete(ctx context.Context, id model.MusicAlbumId) error
}
