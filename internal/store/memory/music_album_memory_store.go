package memory

import (
	"context"
	"database/sql"
	"sync"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
)

// //////////////////////////////////////////////////
// music album memory store

func NewMusicAlbumMemoryStore() store.MusicAlbumStore {
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

func (s *musicAlbumMemoryStore) List(ctx context.Context, tx *sql.Tx, filter *model.MusicAlbumFilter) []*model.MusicAlbum {
	filtered := make([]*model.MusicAlbum, 0, len(s.musicAlbums))
	for _, musicAlbum := range s.musicAlbums {
		if filter.IsMatching(len(filtered), musicAlbum) {
			filtered = append(filtered, musicAlbum)
		}
	}
	return filtered
}

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

func (s *musicAlbumMemoryStore) SearchByName(ctx context.Context, _ *sql.Tx, name string) *model.MusicAlbum {
	s.musicAlbumsLock.Lock()
	defer s.musicAlbumsLock.Unlock()

	if name == "" {
		panic(model.ErrInvalidMusicName)
	}

	for _, musicAlbum := range s.musicAlbums {
		if musicAlbum.Name == name {
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
