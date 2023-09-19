package store

import (
	"context"
	"database/sql"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// game store

type GameStore interface {
	Create(ctx context.Context, tx *sql.Tx, game *model.Game) *model.Game
	Retrieve(ctx context.Context, tx *sql.Tx, id model.GameId) *model.Game
	Update(ctx context.Context, tx *sql.Tx, game *model.Game) *model.Game
	Delete(ctx context.Context, tx *sql.Tx, id model.GameId)
}
