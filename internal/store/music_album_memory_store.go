package store

import (
	"context"
	"database/sql"
	"sync"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// music album memory store

func NewMusicAlbumMemoryStore() MusicAlbumStore {
	return &musicAlbumMemoryStore{
		musicAlbums: make(map[model.MusicAlbumId]*model.MusicAlbum),
	}
}

type musicAlbumMemoryStore struct {
	musicAlbums     map[model.MusicAlbumId]*model.MusicAlbum
	musicAlbumsLock sync.RWMutex
}

var (
	NextMusicAlbumId = 0
)

func (s *musicAlbumMemoryStore) Create(ctx context.Context, _ *sql.Tx, musicAlbum *model.MusicAlbum) *model.MusicAlbum {
	s.musicAlbumsLock.Lock()
	defer s.musicAlbumsLock.Unlock()

	NextMusicAlbumId++
	musicAlbum.Id = model.MusicAlbumId(NextMusicAlbumId)
	s.musicAlbums[musicAlbum.Id] = musicAlbum.Copy()
	return s.musicAlbums[musicAlbum.Id].Copy()
}

func (s *musicAlbumMemoryStore) Retrieve(ctx context.Context, _ *sql.Tx, id model.MusicAlbumId) *model.MusicAlbum {
	s.musicAlbumsLock.Lock()
	defer s.musicAlbumsLock.Unlock()

	musicAlbum, found := s.musicAlbums[id]
	if !found {
		panic(model.ErrMusicAlbumNotFound)
	}
	return musicAlbum.Copy()
}

func (s *musicAlbumMemoryStore) SearchByDeezerId(ctx context.Context, _ *sql.Tx, deezerId model.DeezerAlbumId) *model.MusicAlbum {
	s.musicAlbumsLock.Lock()
	defer s.musicAlbumsLock.Unlock()

	if deezerId == 0 {
		panic(model.ErrInvalidDeezerId)
	}

	for _, musicAlbum := range s.musicAlbums {
		if musicAlbum.DeezerId == deezerId {
			return musicAlbum.Copy()
		}
	}
	return nil
}

func (s *musicAlbumMemoryStore) Update(ctx context.Context, _ *sql.Tx, musicAlbum *model.MusicAlbum) *model.MusicAlbum {
	s.musicAlbumsLock.Lock()
	defer s.musicAlbumsLock.Unlock()

	_, found := s.musicAlbums[musicAlbum.Id]
	if !found {
		panic(model.ErrMusicAlbumNotFound)
	}
	s.musicAlbums[musicAlbum.Id] = musicAlbum.Copy()
	return s.musicAlbums[musicAlbum.Id].Copy()
}

func (s *musicAlbumMemoryStore) Delete(ctx context.Context, _ *sql.Tx, id model.MusicAlbumId) {
	s.musicAlbumsLock.Lock()
	defer s.musicAlbumsLock.Unlock()

	_, found := s.musicAlbums[id]
	if !found {
		panic(model.ErrMusicAlbumNotFound)
	}
	delete(s.musicAlbums, id)
}
