package service

import (
	"context"
	"errors"

	"github.com/gre-ory/amnezic-go/internal/client"
	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// music service

type MusicService interface {
	SearchMusic(ctx context.Context, query string) ([]*model.Music, error)
	AddDeezerMusic(ctx context.Context, deezerId model.DeezerMusicId) (*model.Music, error)
	GetMusic(ctx context.Context, id model.MusicId) (*model.Music, error)
	DeleteMusic(ctx context.Context, id model.MusicId) error
}

func NewMusicService(logger *zap.Logger, deezerClient client.DeezerClient, musicStore store.MusicStore, albumStore store.MusicAlbumStore, artistStore store.MusicArtistStore, genreStore store.MusicGenreStore) MusicService {
	return &musicService{
		logger:       logger,
		deezerClient: deezerClient,
		musicStore:   musicStore,
		albumStore:   albumStore,
		artistStore:  artistStore,
		genreStore:   genreStore,
	}
}

type musicService struct {
	logger       *zap.Logger
	deezerClient client.DeezerClient
	musicStore   store.MusicStore
	albumStore   store.MusicAlbumStore
	artistStore  store.MusicArtistStore
	genreStore   store.MusicGenreStore
}

// //////////////////////////////////////////////////
// search music

func (s *musicService) SearchMusic(ctx context.Context, query string) ([]*model.Music, error) {
	return s.deezerClient.Search(query, 10)
}

// //////////////////////////////////////////////////
// add deezer music

func (s *musicService) AddDeezerMusic(ctx context.Context, deezerId model.DeezerMusicId) (*model.Music, error) {

	//
	// retrieve music from deezer
	//

	music, err := s.deezerClient.GetTrack(deezerId)
	if err != nil {
		return nil, err
	}

	//
	// create music ( if necessary )
	//

	var updated *model.Music
	if music.DeezerMusicId != 0 {
		orig, err := s.musicStore.RetrieveByDeezerId(ctx, music.DeezerMusicId)
		if err != nil {
			if !errors.Is(err, model.ErrMusicNotFound) {
				return nil, err
			} else {
				updated, err = s.musicStore.Create(ctx, music)
				if err != nil {
					return nil, err
				}
			}
		} else {
			updated = orig
		}
	}

	//
	// create album ( if necessary )
	//

	if music.Album != nil && music.Album.DeezerAlbumId != 0 {
		origAlbum, err := s.albumStore.RetrieveByDeezerId(ctx, music.Album.DeezerAlbumId)
		if err != nil {
			if !errors.Is(err, model.ErrMusicAlbumNotFound) {
				return nil, err
			} else {
				updated.Album, err = s.albumStore.Create(ctx, music.Album)
				if err != nil {
					return nil, err
				}
			}
		} else {
			updated.Album = origAlbum
		}
	}

	//
	// create artist ( if necessary )
	//

	if music.Artist != nil && music.Artist.DeezerArtistId != 0 {
		origArtist, err := s.artistStore.RetrieveByDeezerId(ctx, music.Artist.DeezerArtistId)
		if err != nil {
			if !errors.Is(err, model.ErrMusicArtistNotFound) {
				return nil, err
			} else {
				updated.Artist, err = s.artistStore.Create(ctx, music.Artist)
				if err != nil {
					return nil, err
				}
			}
		} else {
			updated.Artist = origArtist
		}
	}

	//
	// create genre ( if necessary )
	//

	if music.Genre != nil && music.Genre.DeezerGenreId != 0 {
		origGenre, err := s.genreStore.RetrieveByDeezerId(ctx, music.Genre.DeezerGenreId)
		if err != nil {
			if !errors.Is(err, model.ErrMusicGenreNotFound) {
				return nil, err
			} else {
				updated.Genre, err = s.genreStore.Create(ctx, music.Genre)
				if err != nil {
					return nil, err
				}
			}
		} else {
			updated.Genre = origGenre
		}
	}

	return updated, nil
}

// //////////////////////////////////////////////////
// get music

func (s *musicService) GetMusic(ctx context.Context, id model.MusicId) (*model.Music, error) {

	//
	// retrieve music
	//

	music, err := s.musicStore.Retrieve(ctx, id)
	if err != nil {
		return nil, err
	}

	//
	// retrieve album
	//

	if music.AlbumId != 0 {
		music.Album, err = s.albumStore.Retrieve(ctx, music.AlbumId)
		if err != nil {
			return nil, err
		}
	}

	//
	// retrieve artist
	//

	if music.ArtistId != 0 {
		music.Artist, err = s.artistStore.Retrieve(ctx, music.ArtistId)
		if err != nil {
			return nil, err
		}
	}

	//
	// retrieve genre
	//

	if music.GenreId != 0 {
		music.Genre, err = s.genreStore.Retrieve(ctx, music.GenreId)
		if err != nil {
			return nil, err
		}
	}

	return music, nil
}

// //////////////////////////////////////////////////
// delete music

func (s *musicService) DeleteMusic(ctx context.Context, id model.MusicId) error {

	//
	// delete music
	//

	_, err := s.musicStore.Retrieve(ctx, id)
	if err != nil {
		return err
	}
	err = s.musicStore.Delete(ctx, id)
	if err != nil {
		return err
	}

	//
	// delete album if no more used
	//

	// TODO

	//
	// delete artist if no more used
	//

	// TODO

	//
	// delete genre if no more used
	//

	// TODO

	return nil
}
