package service

import (
	"context"
	"sync"

	"github.com/gre-ory/amnezic-go/internal/model"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// game service

type GameService interface {
	CreateGame(ctx context.Context, logger *zap.Logger) (*model.Game, error)
	RetrieveGame(ctx context.Context, logger *zap.Logger, id int64) (*model.Game, error)
	DeleteGame(ctx context.Context, logger *zap.Logger, id int64) error
}

func NewGameService() GameService {
	return &gameService{
		games: make(map[int64]*model.Game),
	}
}

type gameService struct {
	games     map[int64]*model.Game
	gamesLock sync.RWMutex
}

func (s *gameService) CreateGame(ctx context.Context, logger *zap.Logger) (*model.Game, error) {
	s.gamesLock.Lock()
	defer s.gamesLock.Unlock()

	id := int64(len(s.games) + 1)
	s.games[id] = &model.Game{
		Id: id,
		Players: []*model.Player{
			{
				Id:     (id * 1000) + 1,
				Name:   "Player 01",
				Active: true,
				Score:  11,
			},
			{
				Id:     (id * 1000) + 2,
				Name:   "Player 02",
				Active: false,
				Score:  12,
			},
		},
	}
	return s.games[id], nil
}

func (s *gameService) RetrieveGame(ctx context.Context, logger *zap.Logger, id int64) (*model.Game, error) {
	s.gamesLock.Lock()
	defer s.gamesLock.Unlock()

	game, found := s.games[id]
	if !found {
		return nil, model.ErrGameNotFound
	}
	return game, nil
}

func (s *gameService) DeleteGame(ctx context.Context, logger *zap.Logger, id int64) error {
	s.gamesLock.Lock()
	defer s.gamesLock.Unlock()

	_, found := s.games[id]
	if !found {
		return model.ErrGameNotFound
	}
	delete(s.games, id)
	return nil
}
