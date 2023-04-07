package store

import (
	"context"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// themeQuestion store

type ThemeQuestionStore interface {
	Create(ctx context.Context, themeQuestion *model.ThemeQuestion) (*model.ThemeQuestion, error)
	Retrieve(ctx context.Context, id model.ThemeQuestionId) (*model.ThemeQuestion, error)
	Update(ctx context.Context, themeQuestion *model.ThemeQuestion) (*model.ThemeQuestion, error)
	Delete(ctx context.Context, filter *model.ThemeQuestionFilter) error
	List(ctx context.Context, filter *model.ThemeQuestionFilter) ([]*model.ThemeQuestion, error)
}
