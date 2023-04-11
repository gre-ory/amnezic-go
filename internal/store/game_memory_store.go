package store

import (
	"context"
	"sync"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// game memory store

func NewGameMemoryStore() GameStore {
	return &gameMemoryStore{
		games: make(map[model.GameId]*model.Game),
	}
}

type gameMemoryStore struct {
	games     map[model.GameId]*model.Game
	gamesLock sync.RWMutex
}

var (
	NextGameId = 0
)

func (s *gameMemoryStore) Create(ctx context.Context, game *model.Game) (*model.Game, error) {
	s.gamesLock.Lock()
	defer s.gamesLock.Unlock()

	NextGameId++
	game.Id = model.NewGameId(NextGameId)
	game.Version = 1
	s.games[game.Id] = game
	return s.games[game.Id], nil
}

func (s *gameMemoryStore) Retrieve(ctx context.Context, id model.GameId) (*model.Game, error) {
	s.gamesLock.Lock()
	defer s.gamesLock.Unlock()

	game, found := s.games[id]
	if !found {
		return nil, model.ErrGameNotFound
	}
	return game, nil
}

func (s *gameMemoryStore) Update(ctx context.Context, game *model.Game) (*model.Game, error) {
	s.gamesLock.Lock()
	defer s.gamesLock.Unlock()

	orig, found := s.games[game.Id]
	if !found {
		return nil, model.ErrGameNotFound
	}
	if orig.Version != game.Version {
		return nil, model.ErrConcurrentUpdate
	}
	game.Version++
	s.games[game.Id] = game
	return s.games[game.Id], nil
}

func (s *gameMemoryStore) Delete(ctx context.Context, id model.GameId) error {
	s.gamesLock.Lock()
	defer s.gamesLock.Unlock()

	_, found := s.games[id]
	if !found {
		return model.ErrGameNotFound
	}
	delete(s.games, id)
	return nil
}
