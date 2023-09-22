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
	UpdateMusic(ctx context.Context, music *model.Music) (*model.Music, error)
	DeleteMusic(ctx context.Context, id model.MusicId) error
}

func NewMusicService(logger *zap.Logger, deezerClient client.DeezerClient, db *sql.DB, musicStore store.MusicStore, albumStore store.MusicAlbumStore, artistStore store.MusicArtistStore, themeStore store.ThemeStore, themeQuestionStore store.ThemeQuestionStore) MusicService {
	return &musicService{
		logger:             logger,
		deezerClient:       deezerClient,
		db:                 db,
		musicStore:         musicStore,
		albumStore:         albumStore,
		artistStore:        artistStore,
		themeStore:         themeStore,
		themeQuestionStore: themeQuestionStore,
	}
}

type musicService struct {
	logger             *zap.Logger
	deezerClient       client.DeezerClient
	db                 *sql.DB
	musicStore         store.MusicStore
	albumStore         store.MusicAlbumStore
	artistStore        store.MusicArtistStore
	themeStore         store.ThemeStore
	themeQuestionStore store.ThemeQuestionStore
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
	s.logger.Info("[DEBUG] music... 1", zap.Object("music", music))

	err = util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// check if music exists
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve music from deezer id %d", deezerId))
		orig := s.musicStore.SearchByDeezerId(ctx, tx, deezerId)
		if orig != nil {
			music = orig

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
		s.logger.Info("[DEBUG] music... 2", zap.Object("music", music))
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
			s.logger.Info("[DEBUG] music... 2.a", zap.Object("music", music))
			s.logger.Info(fmt.Sprintf("[DEBUG] retrieve artist from deezer id %d", music.Artist.DeezerId))
			artist = s.artistStore.SearchByDeezerId(ctx, tx, music.Artist.DeezerId)
			if artist == nil {
				s.logger.Info(fmt.Sprintf("[DEBUG] create artist: %#v", music.Artist.Copy()))
				artist = s.artistStore.Create(ctx, tx, music.Artist)
			}
			music.ArtistId = artist.Id
		}

		//
		// create music
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] create music: %#v", music.Copy()))
		s.logger.Info("[DEBUG] music... 3", zap.Object("music", music))
		music = s.musicStore.Create(ctx, tx, music)
		s.logger.Info("[DEBUG] music... 4", zap.Object("music", music))
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

		//
		// retrieve related questions
		//

		music.Questions = s.themeQuestionStore.List(ctx, tx, &model.ThemeQuestionFilter{MusicId: music.Id})
		music.Questions = util.Convert(music.Questions, s.AttachTheme(ctx, tx))
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] retrieve music %d", id), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] retrieve music %d", id))
	return music, nil
}

// //////////////////////////////////////////////////
// update music

func (s *musicService) UpdateMusic(ctx context.Context, music *model.Music) (*model.Music, error) {

	var updated *model.Music
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {
		updated = s.musicStore.Update(ctx, tx, music)
		updated.Artist = s.artistStore.Update(ctx, tx, music.Artist)
		updated.Album = s.albumStore.Update(ctx, tx, music.Album)
		updated.Questions = s.themeQuestionStore.List(ctx, tx, &model.ThemeQuestionFilter{MusicId: music.Id})
		updated.Questions = util.Convert(updated.Questions, s.AttachTheme(ctx, tx))
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] update music %d - %s", music.Id, music.Name), zap.Object("music", music), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] update music %d - %s", updated.Id, updated.Name), zap.Object("music", updated))
	return updated, nil
}

// //////////////////////////////////////////////////
// delete music

func (s *musicService) DeleteMusic(ctx context.Context, id model.MusicId) error {

	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// check if music is used
		//

		if used := s.themeQuestionStore.IsMusicUsed(ctx, tx, id); used {
			panic(model.ErrMusicUsed)
		}

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

// //////////////////////////////////////////////////
// attach theme

func (s *musicService) AttachTheme(ctx context.Context, tx *sql.Tx) func(question *model.ThemeQuestion) *model.ThemeQuestion {
	return func(question *model.ThemeQuestion) *model.ThemeQuestion {
		if question.ThemeId != 0 {
			question.Theme = s.themeStore.Retrieve(ctx, tx, question.ThemeId)
		}
		return question
	}
}
