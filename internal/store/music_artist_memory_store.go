package store

import (
	"context"
	"database/sql"
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

var (
	NextMusicArtistId = 0
)

func (s *musicArtistMemoryStore) Create(ctx context.Context, _ *sql.Tx, musicArtist *model.MusicArtist) *model.MusicArtist {
	s.musicArtistsLock.Lock()
	defer s.musicArtistsLock.Unlock()

	NextMusicArtistId++
	musicArtist.Id = model.MusicArtistId(NextMusicArtistId)
	s.musicArtists[musicArtist.Id] = musicArtist.Copy()
	return s.musicArtists[musicArtist.Id].Copy()
}

func (s *musicArtistMemoryStore) Retrieve(ctx context.Context, _ *sql.Tx, id model.MusicArtistId) *model.MusicArtist {
	s.musicArtistsLock.Lock()
	defer s.musicArtistsLock.Unlock()

	musicArtist, found := s.musicArtists[id]
	if !found {
		panic(model.ErrMusicArtistNotFound)
	}
	return musicArtist.Copy()
}

func (s *musicArtistMemoryStore) SearchByDeezerId(ctx context.Context, _ *sql.Tx, deezerId model.DeezerArtistId) *model.MusicArtist {
	s.musicArtistsLock.Lock()
	defer s.musicArtistsLock.Unlock()

	if deezerId == 0 {
		panic(model.ErrInvalidDeezerId)
	}

	for _, musicArtist := range s.musicArtists {
		if musicArtist.DeezerId == deezerId {
			return musicArtist.Copy()
		}
	}
	return nil
}

func (s *musicArtistMemoryStore) Update(ctx context.Context, _ *sql.Tx, musicArtist *model.MusicArtist) *model.MusicArtist {
	s.musicArtistsLock.Lock()
	defer s.musicArtistsLock.Unlock()

	_, found := s.musicArtists[musicArtist.Id]
	if !found {
		panic(model.ErrMusicArtistNotFound)
	}
	s.musicArtists[musicArtist.Id] = musicArtist.Copy()
	return s.musicArtists[musicArtist.Id].Copy()
}

func (s *musicArtistMemoryStore) Delete(ctx context.Context, _ *sql.Tx, id model.MusicArtistId) {
	s.musicArtistsLock.Lock()
	defer s.musicArtistsLock.Unlock()

	_, found := s.musicArtists[id]
	if !found {
		panic(model.ErrMusicArtistNotFound)
	}
	delete(s.musicArtists, id)
}
