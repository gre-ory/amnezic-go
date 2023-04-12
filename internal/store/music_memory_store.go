package store

import (
	"context"
	"sync"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// music memory store

func NewMusicMemoryStore() MusicStore {
	return &musicMemoryStore{
		musics: make(map[model.MusicId]*model.Music),
	}
}

type musicMemoryStore struct {
	musics     map[model.MusicId]*model.Music
	musicsLock sync.RWMutex
}

var (
	NextMusicId = 0
)

func (s *musicMemoryStore) Create(ctx context.Context, music *model.Music) (*model.Music, error) {
	s.musicsLock.Lock()
	defer s.musicsLock.Unlock()

	NextMusicId++
	music.Id = model.MusicId(NextMusicId)
	s.musics[music.Id] = music.Copy()
	return s.musics[music.Id].Copy(), nil
}

func (s *musicMemoryStore) Retrieve(ctx context.Context, id model.MusicId) (*model.Music, error) {
	s.musicsLock.Lock()
	defer s.musicsLock.Unlock()

	music, found := s.musics[id]
	if !found {
		return nil, model.ErrMusicNotFound
	}
	return music.Copy(), nil
}

func (s *musicMemoryStore) RetrieveByDeezerId(ctx context.Context, deezerId model.DeezerMusicId) (*model.Music, error) {
	s.musicsLock.Lock()
	defer s.musicsLock.Unlock()

	if deezerId == 0 {
		return nil, model.ErrMusicNotFound
	}

	for _, music := range s.musics {
		if music.DeezerId == deezerId {
			return music.Copy(), nil
		}
	}
	return nil, model.ErrMusicNotFound
}

func (s *musicMemoryStore) Update(ctx context.Context, music *model.Music) (*model.Music, error) {
	s.musicsLock.Lock()
	defer s.musicsLock.Unlock()

	_, found := s.musics[music.Id]
	if !found {
		return nil, model.ErrMusicNotFound
	}
	s.musics[music.Id] = music.Copy()
	return s.musics[music.Id].Copy(), nil
}

func (s *musicMemoryStore) Delete(ctx context.Context, id model.MusicId) error {
	s.musicsLock.Lock()
	defer s.musicsLock.Unlock()

	_, found := s.musics[id]
	if !found {
		return model.ErrMusicNotFound
	}
	delete(s.musics, id)
	return nil
}

func (s *musicMemoryStore) IsAlbumUsed(ctx context.Context, albumId model.MusicAlbumId) (bool, error) {
	s.musicsLock.Lock()
	defer s.musicsLock.Unlock()

	for _, music := range s.musics {
		if music.AlbumId == albumId {
			return true, nil
		}
	}
	return false, nil
}

func (s *musicMemoryStore) IsArtistUsed(ctx context.Context, artistId model.MusicArtistId) (bool, error) {
	s.musicsLock.Lock()
	defer s.musicsLock.Unlock()

	for _, music := range s.musics {
		if music.ArtistId == artistId {
			return true, nil
		}
	}
	return false, nil
}
