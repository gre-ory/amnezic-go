package memory

import (
	"context"
	"database/sql"
	"sync"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
)

// //////////////////////////////////////////////////
// music memory store

func NewMusicMemoryStore() store.MusicStore {
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

func (s *musicMemoryStore) Create(ctx context.Context, _ *sql.Tx, music *model.Music) *model.Music {
	s.musicsLock.Lock()
	defer s.musicsLock.Unlock()

	NextMusicId++
	music.Id = model.MusicId(NextMusicId)
	s.musics[music.Id] = music.Copy()
	return s.musics[music.Id].Copy()
}

func (s *musicMemoryStore) Retrieve(ctx context.Context, _ *sql.Tx, id model.MusicId) *model.Music {
	s.musicsLock.Lock()
	defer s.musicsLock.Unlock()

	music, found := s.musics[id]
	if !found {
		panic(model.ErrMusicNotFound)
	}
	return music.Copy()
}

func (s *musicMemoryStore) SearchByDeezerId(ctx context.Context, _ *sql.Tx, deezerId model.DeezerMusicId) *model.Music {
	s.musicsLock.Lock()
	defer s.musicsLock.Unlock()

	if deezerId == 0 {
		panic(model.ErrInvalidDeezerId)
	}

	for _, music := range s.musics {
		if music.DeezerId == deezerId {
			return music.Copy()
		}
	}
	return nil
}

func (s *musicMemoryStore) Update(ctx context.Context, _ *sql.Tx, music *model.Music) *model.Music {
	s.musicsLock.Lock()
	defer s.musicsLock.Unlock()

	_, found := s.musics[music.Id]
	if !found {
		panic(model.ErrMusicNotFound)
	}
	s.musics[music.Id] = music.Copy()
	return s.musics[music.Id].Copy()
}

func (s *musicMemoryStore) Delete(ctx context.Context, _ *sql.Tx, id model.MusicId) {
	s.musicsLock.Lock()
	defer s.musicsLock.Unlock()

	_, found := s.musics[id]
	if !found {
		panic(model.ErrMusicNotFound)
	}
	delete(s.musics, id)
}

func (s *musicMemoryStore) IsAlbumUsed(ctx context.Context, _ *sql.Tx, albumId model.MusicAlbumId) bool {
	s.musicsLock.Lock()
	defer s.musicsLock.Unlock()

	for _, music := range s.musics {
		if music.AlbumId == albumId {
			return true
		}
	}
	return false
}

func (s *musicMemoryStore) IsArtistUsed(ctx context.Context, _ *sql.Tx, artistId model.MusicArtistId) bool {
	s.musicsLock.Lock()
	defer s.musicsLock.Unlock()

	for _, music := range s.musics {
		if music.ArtistId == artistId {
			return true
		}
	}
	return false
}
