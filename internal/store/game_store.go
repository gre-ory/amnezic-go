package store

import (
	"context"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// game store

type GameStore interface {
	Create(ctx context.Context, game *model.Game) (*model.Game, error)
	Retrieve(ctx context.Context, id model.GameId) (*model.Game, error)
	Update(ctx context.Context, game *model.Game) (*model.Game, error)
	Delete(ctx context.Context, id model.GameId) error
}
