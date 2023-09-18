package store

import (
	"context"
	"database/sql"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// theme store

type ThemeStore interface {
	Create(ctx context.Context, tx *sql.Tx, theme *model.Theme) *model.Theme
	Retrieve(ctx context.Context, tx *sql.Tx, id model.ThemeId) *model.Theme
	Update(ctx context.Context, tx *sql.Tx, theme *model.Theme) *model.Theme
	Delete(ctx context.Context, tx *sql.Tx, filter *model.ThemeFilter)
	List(ctx context.Context, tx *sql.Tx, filter *model.ThemeFilter) []*model.Theme
}
