package store

import (
	"context"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// music genre store

type MusicGenreStore interface {
	Create(ctx context.Context, music *model.MusicGenre) (*model.MusicGenre, error)
	Retrieve(ctx context.Context, id model.MusicGenreId) (*model.MusicGenre, error)
	RetrieveByDeezerId(ctx context.Context, deezerId model.DeezerGenreId) (*model.MusicGenre, error)
	Update(ctx context.Context, music *model.MusicGenre) (*model.MusicGenre, error)
	Delete(ctx context.Context, id model.MusicGenreId) error
}
