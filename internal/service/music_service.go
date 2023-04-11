package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/gre-ory/amnezic-go/internal/client"
	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// music service

type MusicService interface {
	SearchMusic(ctx context.Context, query string, limit int) ([]*model.Music, error)
	AddDeezerMusic(ctx context.Context, deezerId model.DeezerMusicId) (*model.Music, error)
	GetMusic(ctx context.Context, id model.MusicId) (*model.Music, error)
	DeleteMusic(ctx context.Context, id model.MusicId) error
}

func NewMusicService(logger *zap.Logger, deezerClient client.DeezerClient, musicStore store.MusicStore, albumStore store.MusicAlbumStore, artistStore store.MusicArtistStore) MusicService {
	return &musicService{
		logger:       logger,
		deezerClient: deezerClient,
		musicStore:   musicStore,
		albumStore:   albumStore,
		artistStore:  artistStore,
	}
}

type musicService struct {
	logger       *zap.Logger
	deezerClient client.DeezerClient
	musicStore   store.MusicStore
	albumStore   store.MusicAlbumStore
	artistStore  store.MusicArtistStore
}

// //////////////////////////////////////////////////
// search music

func (s *musicService) SearchMusic(ctx context.Context, query string, limit int) ([]*model.Music, error) {
	return s.deezerClient.Search(query, limit)
}

// //////////////////////////////////////////////////
// add deezer music

func (s *musicService) AddDeezerMusic(ctx context.Context, deezerId model.DeezerMusicId) (*model.Music, error) {

	var music *model.Music
	var album *model.MusicAlbum
	var artist *model.MusicArtist
	var err error

	defer func() {
		if err == nil {
			s.logger.Info(fmt.Sprintf("[ OK ] add deezer music %d", deezerId))
		} else {
			s.logger.Info(fmt.Sprintf("[ KO ] add deezer music %d", deezerId), zap.Error(err))
		}
	}()

	//
	// check if music exists
	//

	if deezerId != 0 {
		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve music from deezer id %d", deezerId))
		music, err = s.musicStore.RetrieveByDeezerId(ctx, deezerId)
		if err == nil {

			//
			// retrieve album
			//

			if music.AlbumId != 0 {
				s.logger.Info(fmt.Sprintf("[DEBUG] retrieve album %d", music.AlbumId))
				music.Album, err = s.albumStore.Retrieve(ctx, music.AlbumId)
				if err != nil {
					return nil, err
				}
			}

			//
			// retrieve artist
			//

			if music.ArtistId != 0 {
				s.logger.Info(fmt.Sprintf("[DEBUG] retrieve artist %d", music.ArtistId))
				music.Artist, err = s.artistStore.Retrieve(ctx, music.ArtistId)
				if err != nil {
					return nil, err
				}
			}

			return music, nil
		} else if !errors.Is(err, model.ErrMusicNotFound) {
			return nil, err
		}
	}

	//
	// retrieve music from deezer
	//

	s.logger.Info(fmt.Sprintf("[DEBUG] search music %d from deezer", deezerId))
	music, err = s.deezerClient.GetTrack(deezerId)
	if err != nil {
		return nil, err
	}

	//
	// create album ( if necessary )
	//

	if music.Album != nil && music.Album.DeezerAlbumId != 0 {
		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve album from deezer id %d", music.Album.DeezerAlbumId))
		album, err = s.albumStore.RetrieveByDeezerId(ctx, music.Album.DeezerAlbumId)
		if err != nil {
			if !errors.Is(err, model.ErrMusicAlbumNotFound) {
				return nil, err
			} else {
				s.logger.Info(fmt.Sprintf("[DEBUG] create album: %#v", music.Album.Copy()))
				album, err = s.albumStore.Create(ctx, music.Album)
				if err != nil {
					return nil, err
				}
			}
		}
		music.AlbumId = album.Id
	}

	//
	// create artist ( if necessary )
	//

	if music.Artist != nil && music.Artist.DeezerArtistId != 0 {
		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve artist from deezer id %d", music.Artist.DeezerArtistId))
		artist, err = s.artistStore.RetrieveByDeezerId(ctx, music.Artist.DeezerArtistId)
		if err != nil {
			if !errors.Is(err, model.ErrMusicArtistNotFound) {
				return nil, err
			} else {
				s.logger.Info(fmt.Sprintf("[DEBUG] create artist: %#v", music.Artist.Copy()))
				artist, err = s.artistStore.Create(ctx, music.Artist)
				if err != nil {
					return nil, err
				}
			}
		}
		music.ArtistId = artist.Id
	}

	//
	// create music ( if necessary )
	//

	s.logger.Info(fmt.Sprintf("[DEBUG] create music: %#v", music.Copy()))
	music, err = s.musicStore.Create(ctx, music)
	if err != nil {
		return nil, err
	}

	music.Album = album
	music.Artist = artist

	return music, nil
}

// //////////////////////////////////////////////////
// get music

func (s *musicService) GetMusic(ctx context.Context, id model.MusicId) (*model.Music, error) {

	var music *model.Music
	var err error

	defer func() {
		if err == nil {
			s.logger.Info(fmt.Sprintf("[ OK ] retrieve music %d", id))
		} else {
			s.logger.Info(fmt.Sprintf("[ KO ] retrieve music %d", id), zap.Error(err))
		}
	}()

	//
	// retrieve music
	//

	s.logger.Info(fmt.Sprintf("[DEBUG] retrieve music %d", id))
	music, err = s.musicStore.Retrieve(ctx, id)
	if err != nil {
		return nil, err
	}

	//
	// retrieve album
	//

	if music.AlbumId != 0 {
		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve album %d", music.AlbumId))
		music.Album, err = s.albumStore.Retrieve(ctx, music.AlbumId)
		if err != nil {
			return nil, err
		}
	}

	//
	// retrieve artist
	//

	if music.ArtistId != 0 {
		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve artist %d", music.ArtistId))
		music.Artist, err = s.artistStore.Retrieve(ctx, music.ArtistId)
		if err != nil {
			return nil, err
		}
	}

	return music, nil
}

// //////////////////////////////////////////////////
// delete music

func (s *musicService) DeleteMusic(ctx context.Context, id model.MusicId) error {

	var err error

	defer func() {
		if err == nil {
			s.logger.Info(fmt.Sprintf("[ OK ] delete music %d", id))
		} else {
			s.logger.Info(fmt.Sprintf("[ KO ] delete music %d", id), zap.Error(err))
		}
	}()

	//
	// delete music
	//

	s.logger.Info(fmt.Sprintf("[DEBUG] retrieve music %d", id))
	music, err := s.musicStore.Retrieve(ctx, id)
	if err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("[DEBUG] delete music %d", id))
	err = s.musicStore.Delete(ctx, id)
	if err != nil {
		return err
	}

	//
	// delete album if no more used
	//

	if music.AlbumId != 0 {
		s.logger.Info(fmt.Sprintf("[DEBUG] check usage of album %d", music.AlbumId))
		used, err := s.musicStore.IsAlbumUsed(ctx, music.AlbumId)
		if err != nil {
			return err
		}
		if !used {
			s.logger.Info(fmt.Sprintf("[DEBUG] delete unused album %d", music.AlbumId))
			err = s.albumStore.Delete(ctx, music.AlbumId)
			if err != nil {
				return err
			}
		}
	}

	//
	// delete artist if no more used
	//

	if music.ArtistId != 0 {
		s.logger.Info(fmt.Sprintf("[DEBUG] check usage of artist %d", music.ArtistId))
		used, err := s.musicStore.IsArtistUsed(ctx, music.ArtistId)
		if err != nil {
			return err
		}
		if !used {
			s.logger.Info(fmt.Sprintf("[DEBUG] delete unused artist %d", music.ArtistId))
			err = s.artistStore.Delete(ctx, music.ArtistId)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
