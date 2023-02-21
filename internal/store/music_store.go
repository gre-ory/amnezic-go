package store

import (
	"context"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// music store

type MusicStore interface {
	SelectRandomQuestions(cxt context.Context, settings model.GameSettings) ([]*model.Question, error)
}
