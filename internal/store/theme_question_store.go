package store

import (
	"context"
	"database/sql"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// themeQuestion store

type ThemeQuestionStore interface {
	Create(ctx context.Context, tx *sql.Tx, themeQuestion *model.ThemeQuestion) *model.ThemeQuestion
	Retrieve(ctx context.Context, tx *sql.Tx, id model.ThemeQuestionId) *model.ThemeQuestion
	Update(ctx context.Context, tx *sql.Tx, themeQuestion *model.ThemeQuestion) *model.ThemeQuestion
	Delete(ctx context.Context, tx *sql.Tx, filter *model.ThemeQuestionFilter)
	List(ctx context.Context, tx *sql.Tx, filter *model.ThemeQuestionFilter) []*model.ThemeQuestion
	CountByTheme(ctx context.Context, tx *sql.Tx) map[model.ThemeId]int
}
