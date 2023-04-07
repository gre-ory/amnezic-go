package store

import (
	"context"
	"sync"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// music album memory store

func NewMusicArtistMemoryStore() MusicArtistStore {
	return &musicArtistMemoryStore{
		musicArtists: make(map[model.MusicArtistId]*model.MusicArtist),
	}
}

type musicArtistMemoryStore struct {
	musicArtists     map[model.MusicArtistId]*model.MusicArtist
	musicArtistsLock sync.RWMutex
}

func (s *musicArtistMemoryStore) Create(ctx context.Context, musicArtist *model.MusicArtist) (*model.MusicArtist, error) {
	s.musicArtistsLock.Lock()
	defer s.musicArtistsLock.Unlock()

	musicArtistNumber := len(s.musicArtists) + 1
	musicArtist.Id = model.MusicArtistId(musicArtistNumber)
	s.musicArtists[musicArtist.Id] = musicArtist.Copy()
	return s.musicArtists[musicArtist.Id].Copy(), nil
}

func (s *musicArtistMemoryStore) Retrieve(ctx context.Context, id model.MusicArtistId) (*model.MusicArtist, error) {
	s.musicArtistsLock.Lock()
	defer s.musicArtistsLock.Unlock()

	musicArtist, found := s.musicArtists[id]
	if !found {
		return nil, model.ErrMusicArtistNotFound
	}
	return musicArtist.Copy(), nil
}

func (s *musicArtistMemoryStore) RetrieveByDeezerId(ctx context.Context, deezerId model.DeezerArtistId) (*model.MusicArtist, error) {
	s.musicArtistsLock.Lock()
	defer s.musicArtistsLock.Unlock()

	if deezerId == 0 {
		return nil, model.ErrMusicArtistNotFound
	}

	for _, musicArtist := range s.musicArtists {
		if musicArtist.DeezerArtistId == deezerId {
			return musicArtist.Copy(), nil
		}
	}
	return nil, model.ErrMusicArtistNotFound
}

func (s *musicArtistMemoryStore) Update(ctx context.Context, musicArtist *model.MusicArtist) (*model.MusicArtist, error) {
	s.musicArtistsLock.Lock()
	defer s.musicArtistsLock.Unlock()

	_, found := s.musicArtists[musicArtist.Id]
	if !found {
		return nil, model.ErrMusicArtistNotFound
	}
	s.musicArtists[musicArtist.Id] = musicArtist.Copy()
	return s.musicArtists[musicArtist.Id].Copy(), nil
}

func (s *musicArtistMemoryStore) Delete(ctx context.Context, id model.MusicArtistId) error {
	s.musicArtistsLock.Lock()
	defer s.musicArtistsLock.Unlock()

	_, found := s.musicArtists[id]
	if !found {
		return model.ErrMusicArtistNotFound
	}
	delete(s.musicArtists, id)
	return nil
}
