package memory

import (
	"context"
	"database/sql"
	"sync"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
)

// //////////////////////////////////////////////////
// game memory store

func NewGameMemoryStore() store.GameStore {
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

func (s *gameMemoryStore) Create(ctx context.Context, _ *sql.Tx, game *model.Game) *model.Game {
	s.gamesLock.Lock()
	defer s.gamesLock.Unlock()

	NextGameId++
	game.Id = model.NewGameId(NextGameId)
	game.Version = 1
	s.games[game.Id] = game
	return s.games[game.Id]
}

func (s *gameMemoryStore) Retrieve(ctx context.Context, _ *sql.Tx, id model.GameId) *model.Game {
	s.gamesLock.Lock()
	defer s.gamesLock.Unlock()

	game, found := s.games[id]
	if !found {
		panic(model.ErrGameNotFound)
	}
	return game
}

func (s *gameMemoryStore) Update(ctx context.Context, _ *sql.Tx, game *model.Game) *model.Game {
	s.gamesLock.Lock()
	defer s.gamesLock.Unlock()

	orig, found := s.games[game.Id]
	if !found {
		panic(model.ErrGameNotFound)
	}
	if orig.Version != game.Version {
		panic(model.ErrConcurrentUpdate)
	}
	game.Version++
	s.games[game.Id] = game
	return s.games[game.Id]
}

func (s *gameMemoryStore) Delete(ctx context.Context, _ *sql.Tx, id model.GameId) {
	s.gamesLock.Lock()
	defer s.gamesLock.Unlock()

	_, found := s.games[id]
	if !found {
		panic(model.ErrGameNotFound)
	}
	delete(s.games, id)
}
