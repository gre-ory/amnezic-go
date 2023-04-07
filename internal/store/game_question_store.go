package store

import (
	"context"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// game question store

type GameQuestionStore interface {
	SelectRandomQuestions(cxt context.Context, settings model.GameSettings) ([]*model.GameQuestion, error)
}
