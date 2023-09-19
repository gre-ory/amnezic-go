package store

import (
	"context"
	"database/sql"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// game question store

type GameQuestionStore interface {
	SelectRandomQuestions(cxt context.Context, tx *sql.Tx, settings model.GameSettings) []*model.GameQuestion
}
