package store

import (
	"context"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// music store

type MusicStore interface {
	Create(ctx context.Context, music *model.Music) (*model.Music, error)
	Retrieve(ctx context.Context, id model.MusicId) (*model.Music, error)
	RetrieveByDeezerId(ctx context.Context, deezerId model.DeezerMusicId) (*model.Music, error)
	Update(ctx context.Context, music *model.Music) (*model.Music, error)
	Delete(ctx context.Context, id model.MusicId) error
}
