package store

import (
	"context"
	"database/sql"
	"sync"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// theme memory store

func NewThemeMemoryStore() ThemeStore {
	return &themeMemoryStore{
		themes: make(map[model.ThemeId]*model.Theme),
	}
}

type themeMemoryStore struct {
	themes     map[model.ThemeId]*model.Theme
	themesLock sync.RWMutex
}

var (
	NextThemeId = 0
)

func (s *themeMemoryStore) Create(ctx context.Context, _ *sql.Tx, theme *model.Theme) *model.Theme {
	s.themesLock.Lock()
	defer s.themesLock.Unlock()

	NextThemeId++
	theme.Id = model.ThemeId(NextThemeId)
	s.themes[theme.Id] = theme.Copy()
	return s.themes[theme.Id].Copy()
}

func (s *themeMemoryStore) Retrieve(ctx context.Context, _ *sql.Tx, id model.ThemeId) *model.Theme {
	s.themesLock.Lock()
	defer s.themesLock.Unlock()

	theme, found := s.themes[id]
	if !found {
		return nil
	}
	return theme.Copy()
}

func (s *themeMemoryStore) Update(ctx context.Context, _ *sql.Tx, theme *model.Theme) *model.Theme {
	s.themesLock.Lock()
	defer s.themesLock.Unlock()

	_, found := s.themes[theme.Id]
	if !found {
		panic(model.ErrThemeNotFound)
	}
	s.themes[theme.Id] = theme.Copy()
	return s.themes[theme.Id].Copy()
}

func (s *themeMemoryStore) Delete(ctx context.Context, _ *sql.Tx, filter *model.ThemeFilter) {
	s.themesLock.Lock()
	defer s.themesLock.Unlock()

	for id, theme := range s.themes {
		if filter.IsMatching(theme) {
			delete(s.themes, id)
		}
	}
}

func (s *themeMemoryStore) List(ctx context.Context, _ *sql.Tx, filter *model.ThemeFilter) []*model.Theme {
	s.themesLock.Lock()
	defer s.themesLock.Unlock()

	themes := make([]*model.Theme, 0, len(s.themes))
	for _, theme := range s.themes {
		if filter.IsMatching(theme) {
			themes = append(themes, theme.Copy())
		}
	}

	return themes
}
