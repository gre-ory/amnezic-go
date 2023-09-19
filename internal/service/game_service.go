package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// game service

type GameService interface {
	CreateGame(ctx context.Context, settings model.GameSettings) (*model.Game, error)
	RetrieveGame(ctx context.Context, id model.GameId) (*model.Game, error)
	DeleteGame(ctx context.Context, id model.GameId) error
}

func NewGameService(logger *zap.Logger, db *sql.DB, gameStore store.GameStore, musicStore store.GameQuestionStore) GameService {
	return &gameService{
		logger:     logger,
		db:         db,
		gameStore:  gameStore,
		musicStore: musicStore,
	}
}

type gameService struct {
	logger     *zap.Logger
	db         *sql.DB
	gameStore  store.GameStore
	musicStore store.GameQuestionStore
}

func (s *gameService) CreateGame(ctx context.Context, settings model.GameSettings) (*model.Game, error) {

	var game *model.Game
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		questions := s.musicStore.SelectRandomQuestions(ctx, tx, settings)

		game = &model.Game{
			Settings:  &settings,
			Players:   s.createPlayers(settings.NbPlayer),
			Questions: questions,
		}

		game = s.gameStore.Create(ctx, tx, game)

		for questionIndex, question := range game.Questions {
			question.Id = model.NewGameQuestionId(game.Id, questionIndex+1)
			for answerIndex, answer := range question.Answers {
				answer.Id = model.NewGameAnswerId(question.Id, answerIndex+1)
			}
		}
	})

	if err != nil {
		return nil, err
	}
	return game, nil
}

func (s *gameService) createPlayers(nbPlayer int) []*model.GamePlayer {
	players := make([]*model.GamePlayer, 0, nbPlayer)
	for playerNumber := 1; playerNumber <= nbPlayer; playerNumber++ {
		players = append(players, s.createPlayer(playerNumber))
	}
	return players
}

func (s *gameService) createPlayer(playerNumber int) *model.GamePlayer {
	return &model.GamePlayer{
		Id:     model.NewGamePlayerId(playerNumber),
		Name:   fmt.Sprintf("Player %02d", playerNumber),
		Active: true,
		Score:  0,
	}
}

func (s *gameService) RetrieveGame(ctx context.Context, id model.GameId) (*model.Game, error) {

	var game *model.Game
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {
		game = s.gameStore.Retrieve(ctx, tx, id)
	})
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (s *gameService) DeleteGame(ctx context.Context, id model.GameId) error {
	return util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {
		s.gameStore.Delete(ctx, tx, id)
	})
}
