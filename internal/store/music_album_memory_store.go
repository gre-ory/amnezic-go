package store

import (
	"context"
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

func (s *musicAlbumMemoryStore) Create(ctx context.Context, musicAlbum *model.MusicAlbum) (*model.MusicAlbum, error) {
	s.musicAlbumsLock.Lock()
	defer s.musicAlbumsLock.Unlock()

	NextMusicAlbumId++
	musicAlbum.Id = model.MusicAlbumId(NextMusicAlbumId)
	s.musicAlbums[musicAlbum.Id] = musicAlbum.Copy()
	return s.musicAlbums[musicAlbum.Id].Copy(), nil
}

func (s *musicAlbumMemoryStore) Retrieve(ctx context.Context, id model.MusicAlbumId) (*model.MusicAlbum, error) {
	s.musicAlbumsLock.Lock()
	defer s.musicAlbumsLock.Unlock()

	musicAlbum, found := s.musicAlbums[id]
	if !found {
		return nil, model.ErrMusicAlbumNotFound
	}
	return musicAlbum.Copy(), nil
}

func (s *musicAlbumMemoryStore) RetrieveByDeezerId(ctx context.Context, deezerId model.DeezerAlbumId) (*model.MusicAlbum, error) {
	s.musicAlbumsLock.Lock()
	defer s.musicAlbumsLock.Unlock()

	if deezerId == 0 {
		return nil, model.ErrMusicAlbumNotFound
	}

	for _, musicAlbum := range s.musicAlbums {
		if musicAlbum.DeezerId == deezerId {
			return musicAlbum.Copy(), nil
		}
	}
	return nil, model.ErrMusicAlbumNotFound
}

func (s *musicAlbumMemoryStore) Update(ctx context.Context, musicAlbum *model.MusicAlbum) (*model.MusicAlbum, error) {
	s.musicAlbumsLock.Lock()
	defer s.musicAlbumsLock.Unlock()

	_, found := s.musicAlbums[musicAlbum.Id]
	if !found {
		return nil, model.ErrMusicAlbumNotFound
	}
	s.musicAlbums[musicAlbum.Id] = musicAlbum.Copy()
	return s.musicAlbums[musicAlbum.Id].Copy(), nil
}

func (s *musicAlbumMemoryStore) Delete(ctx context.Context, id model.MusicAlbumId) error {
	s.musicAlbumsLock.Lock()
	defer s.musicAlbumsLock.Unlock()

	_, found := s.musicAlbums[id]
	if !found {
		return model.ErrMusicAlbumNotFound
	}
	delete(s.musicAlbums, id)
	return nil
}
