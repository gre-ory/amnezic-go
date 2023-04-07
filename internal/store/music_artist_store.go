package store

import (
	"context"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// music artist store

type MusicArtistStore interface {
	Create(ctx context.Context, music *model.MusicArtist) (*model.MusicArtist, error)
	Retrieve(ctx context.Context, id model.MusicArtistId) (*model.MusicArtist, error)
	RetrieveByDeezerId(ctx context.Context, deezerId model.DeezerArtistId) (*model.MusicArtist, error)
	Update(ctx context.Context, music *model.MusicArtist) (*model.MusicArtist, error)
	Delete(ctx context.Context, id model.MusicArtistId) error
}
