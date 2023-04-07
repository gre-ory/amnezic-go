package store

import (
	"context"
	"sync"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// music album memory store

func NewMusicGenreMemoryStore() MusicGenreStore {
	return &musicGenreMemoryStore{
		musicGenres: make(map[model.MusicGenreId]*model.MusicGenre),
	}
}

type musicGenreMemoryStore struct {
	musicGenres     map[model.MusicGenreId]*model.MusicGenre
	musicGenresLock sync.RWMutex
}

func (s *musicGenreMemoryStore) Create(ctx context.Context, musicGenre *model.MusicGenre) (*model.MusicGenre, error) {
	s.musicGenresLock.Lock()
	defer s.musicGenresLock.Unlock()

	musicGenreNumber := len(s.musicGenres) + 1
	musicGenre.Id = model.MusicGenreId(musicGenreNumber)
	s.musicGenres[musicGenre.Id] = musicGenre.Copy()
	return s.musicGenres[musicGenre.Id].Copy(), nil
}

func (s *musicGenreMemoryStore) Retrieve(ctx context.Context, id model.MusicGenreId) (*model.MusicGenre, error) {
	s.musicGenresLock.Lock()
	defer s.musicGenresLock.Unlock()

	musicGenre, found := s.musicGenres[id]
	if !found {
		return nil, model.ErrMusicGenreNotFound
	}
	return musicGenre.Copy(), nil
}

func (s *musicGenreMemoryStore) RetrieveByDeezerId(ctx context.Context, deezerId model.DeezerGenreId) (*model.MusicGenre, error) {
	s.musicGenresLock.Lock()
	defer s.musicGenresLock.Unlock()

	if deezerId == 0 {
		return nil, model.ErrMusicGenreNotFound
	}

	for _, musicGenre := range s.musicGenres {
		if musicGenre.DeezerGenreId == deezerId {
			return musicGenre.Copy(), nil
		}
	}
	return nil, model.ErrMusicGenreNotFound
}

func (s *musicGenreMemoryStore) Update(ctx context.Context, musicGenre *model.MusicGenre) (*model.MusicGenre, error) {
	s.musicGenresLock.Lock()
	defer s.musicGenresLock.Unlock()

	_, found := s.musicGenres[musicGenre.Id]
	if !found {
		return nil, model.ErrMusicGenreNotFound
	}
	s.musicGenres[musicGenre.Id] = musicGenre.Copy()
	return s.musicGenres[musicGenre.Id].Copy(), nil
}

func (s *musicGenreMemoryStore) Delete(ctx context.Context, id model.MusicGenreId) error {
	s.musicGenresLock.Lock()
	defer s.musicGenresLock.Unlock()

	_, found := s.musicGenres[id]
	if !found {
		return model.ErrMusicGenreNotFound
	}
	delete(s.musicGenres, id)
	return nil
}
