package store

import (
	"context"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// theme store

type ThemeStore interface {
	Create(ctx context.Context, theme *model.Theme) (*model.Theme, error)
	Retrieve(ctx context.Context, id model.ThemeId) (*model.Theme, error)
	Update(ctx context.Context, theme *model.Theme) (*model.Theme, error)
	Delete(ctx context.Context, filter *model.ThemeFilter) error
	List(ctx context.Context, filter *model.ThemeFilter) ([]*model.Theme, error)
}
