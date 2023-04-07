package service

import (
	"context"
	"fmt"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// game service

type GameService interface {
	CreateGame(ctx context.Context, settings model.GameSettings) (*model.Game, error)
	RetrieveGame(ctx context.Context, id model.GameId) (*model.Game, error)
	DeleteGame(ctx context.Context, id model.GameId) error
}

func NewGameService(logger *zap.Logger, gameStore store.GameStore, musicStore store.GameQuestionStore) GameService {
	return &gameService{
		logger:     logger,
		gameStore:  gameStore,
		musicStore: musicStore,
	}
}

type gameService struct {
	logger     *zap.Logger
	gameStore  store.GameStore
	musicStore store.GameQuestionStore
}

func (s *gameService) CreateGame(ctx context.Context, settings model.GameSettings) (*model.Game, error) {

	var questions []*model.GameQuestion
	var err error

	questions, err = s.musicStore.SelectRandomQuestions(ctx, settings)
	if err != nil {
		return nil, err
	}

	game := &model.Game{
		Settings:  &settings,
		Players:   s.createPlayers(settings.NbPlayer),
		Questions: questions,
	}

	game, err = s.gameStore.Create(ctx, game)
	if err != nil {
		return nil, err
	}

	for questionIndex, question := range game.Questions {
		question.Id = model.NewGameQuestionId(game.Id, questionIndex+1)
		for answerIndex, answer := range question.Answers {
			answer.Id = model.NewGameAnswerId(question.Id, answerIndex+1)
		}
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
	game, err := s.gameStore.Retrieve(ctx, id)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (s *gameService) DeleteGame(ctx context.Context, id model.GameId) error {
	return s.gameStore.Delete(ctx, id)
}
