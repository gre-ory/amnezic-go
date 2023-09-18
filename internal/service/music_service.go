package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gre-ory/amnezic-go/internal/client"
	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
	"github.com/gre-ory/amnezic-go/internal/util"
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

func NewMusicService(logger *zap.Logger, deezerClient client.DeezerClient, db *sql.DB, musicStore store.MusicStore, albumStore store.MusicAlbumStore, artistStore store.MusicArtistStore) MusicService {
	return &musicService{
		logger:       logger,
		deezerClient: deezerClient,
		db:           db,
		musicStore:   musicStore,
		albumStore:   albumStore,
		artistStore:  artistStore,
	}
}

type musicService struct {
	logger       *zap.Logger
	deezerClient client.DeezerClient
	db           *sql.DB
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

	if deezerId == 0 {
		return nil, model.ErrInvalidDeezerId
	}

	var music *model.Music
	var album *model.MusicAlbum
	var artist *model.MusicArtist
	var err error

	//
	// retrieve music from deezer
	//

	// TODO NOT optimal to search on deezer if it already exists in DB
	s.logger.Info(fmt.Sprintf("[DEBUG] search music %d from deezer", deezerId))
	music, err = s.deezerClient.GetTrack(deezerId)
	if err != nil {
		return nil, err
	}

	err = util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// check if music exists
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve music from deezer id %d", deezerId))
		music = s.musicStore.SearchByDeezerId(ctx, tx, deezerId)
		if music != nil {

			//
			// retrieve album
			//

			if music.AlbumId != 0 {
				s.logger.Info(fmt.Sprintf("[DEBUG] retrieve album %d", music.AlbumId))
				music.Album = s.albumStore.Retrieve(ctx, tx, music.AlbumId)
			}

			//
			// retrieve artist
			//

			if music.ArtistId != 0 {
				s.logger.Info(fmt.Sprintf("[DEBUG] retrieve artist %d", music.ArtistId))
				music.Artist = s.artistStore.Retrieve(ctx, tx, music.ArtistId)
			}

			// stop as it already exists
			return
		}

		//
		// create album ( if necessary )
		//

		if music.Album != nil && music.Album.DeezerId != 0 {
			s.logger.Info(fmt.Sprintf("[DEBUG] retrieve album from deezer id %d", music.Album.DeezerId))
			album = s.albumStore.SearchByDeezerId(ctx, tx, music.Album.DeezerId)
			if album == nil {
				s.logger.Info(fmt.Sprintf("[DEBUG] create album: %#v", music.Album.Copy()))
				album = s.albumStore.Create(ctx, tx, music.Album)
			}
			music.AlbumId = album.Id
		}

		//
		// create artist ( if necessary )
		//

		if music.Artist != nil && music.Artist.DeezerId != 0 {
			s.logger.Info(fmt.Sprintf("[DEBUG] retrieve artist from deezer id %d", music.Artist.DeezerId))
			artist = s.artistStore.SearchByDeezerId(ctx, tx, music.Artist.DeezerId)
			if artist != nil {
				s.logger.Info(fmt.Sprintf("[DEBUG] create artist: %#v", music.Artist.Copy()))
				artist = s.artistStore.Create(ctx, tx, music.Artist)
			}
			music.ArtistId = artist.Id
		}

		//
		// create music
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] create music: %#v", music.Copy()))
		music = s.musicStore.Create(ctx, tx, music)
		music.Album = album
		music.Artist = artist
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] add deezer music %d", deezerId), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] add deezer music %d", deezerId))
	return music, nil

}

// //////////////////////////////////////////////////
// get music

func (s *musicService) GetMusic(ctx context.Context, id model.MusicId) (*model.Music, error) {

	var music *model.Music
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// retrieve music
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve music %d", id))
		music = s.musicStore.Retrieve(ctx, tx, id)

		//
		// retrieve album
		//

		if music.AlbumId != 0 {
			s.logger.Info(fmt.Sprintf("[DEBUG] retrieve album %d", music.AlbumId))
			music.Album = s.albumStore.Retrieve(ctx, tx, music.AlbumId)
		}

		//
		// retrieve artist
		//

		if music.ArtistId != 0 {
			s.logger.Info(fmt.Sprintf("[DEBUG] retrieve artist %d", music.ArtistId))
			music.Artist = s.artistStore.Retrieve(ctx, tx, music.ArtistId)
		}

	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] retrieve music %d", id), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] retrieve music %d", id))
	return music, nil
}

// //////////////////////////////////////////////////
// delete music

func (s *musicService) DeleteMusic(ctx context.Context, id model.MusicId) error {

	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// delete music
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve music %d", id))
		music := s.musicStore.Retrieve(ctx, tx, id)

		s.logger.Info(fmt.Sprintf("[DEBUG] delete music %d", id))
		s.musicStore.Delete(ctx, tx, id)

		//
		// delete album if no more used
		//

		if music.AlbumId != 0 {
			s.logger.Info(fmt.Sprintf("[DEBUG] check usage of album %d", music.AlbumId))
			if used := s.musicStore.IsAlbumUsed(ctx, tx, music.AlbumId); !used {
				s.logger.Info(fmt.Sprintf("[DEBUG] delete unused album %d", music.AlbumId))
				s.albumStore.Delete(ctx, tx, music.AlbumId)
			}
		}

		//
		// delete artist if no more used
		//

		if music.ArtistId != 0 {
			s.logger.Info(fmt.Sprintf("[DEBUG] check usage of artist %d", music.ArtistId))
			if used := s.musicStore.IsArtistUsed(ctx, tx, music.ArtistId); !used {
				s.logger.Info(fmt.Sprintf("[DEBUG] delete unused artist %d", music.ArtistId))
				s.artistStore.Delete(ctx, tx, music.ArtistId)
			}
		}
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] delete music %d", id), zap.Error(err))
		return err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] delete music %d", id))
	return nil
}
