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
// artist service

type ArtistService interface {
	ListArtist(ctx context.Context, search *model.MusicArtistFilter) ([]*model.MusicArtist, error)
	GetArtist(ctx context.Context, id model.MusicArtistId) (*model.MusicArtist, error)
	CreateArtist(ctx context.Context, music *model.MusicArtist) (*model.MusicArtist, error)
	UpdateArtist(ctx context.Context, music *model.MusicArtist) (*model.MusicArtist, error)
	DeleteArtist(ctx context.Context, id model.MusicArtistId) error
}

func NewArtistService(logger *zap.Logger, downloadClient client.DownloadClient, db *sql.DB, artistStore store.MusicArtistStore, musicStore store.MusicStore, imageFileValidator model.PathValidator) ArtistService {
	return &artistService{
		logger:             logger,
		downloadClient:     downloadClient,
		db:                 db,
		artistStore:        artistStore,
		musicStore:         musicStore,
		imageFileValidator: imageFileValidator,
	}
}

type artistService struct {
	logger             *zap.Logger
	downloadClient     client.DownloadClient
	db                 *sql.DB
	artistStore        store.MusicArtistStore
	musicStore         store.MusicStore
	imageFileValidator model.PathValidator
}

// //////////////////////////////////////////////////
// list artist

func (s *artistService) ListArtist(ctx context.Context, filter *model.MusicArtistFilter) ([]*model.MusicArtist, error) {
	var artists []*model.MusicArtist
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {
		artists = s.artistStore.List(ctx, tx, filter)
	})
	if err != nil {
		s.logger.Info("[ KO ] list artist", zap.Error(err))
		return nil, err
	}
	s.logger.Info("[ OK ] list artist")
	return artists, nil
}

// //////////////////////////////////////////////////
// create artist

func (s *artistService) CreateArtist(ctx context.Context, artist *model.MusicArtist) (*model.MusicArtist, error) {

	var err error

	//
	// validate
	//

	if artist == nil {
		return nil, model.ErrMissingArtist
	}
	if err = artist.Validate(s.imageFileValidator); err != nil {
		return nil, err
	}

	//
	// download remote files
	//

	s.DownloadRemoteFiles(ctx, artist)

	//
	// create
	//

	err = util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// check if artist exists
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve artist from name %q", artist.Name))
		other := s.artistStore.SearchByName(ctx, tx, artist.Name)
		if other != nil {
			panic(model.ErrExistingArtist)
		}

		//
		// create artist
		//

		artist = s.artistStore.Create(ctx, tx, artist)
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] create artist %q", artist.Name), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] create artist %q", artist.Name))
	return artist, nil

}

// //////////////////////////////////////////////////
// get artist

func (s *artistService) GetArtist(ctx context.Context, id model.MusicArtistId) (*model.MusicArtist, error) {

	var artist *model.MusicArtist
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// retrieve artist
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve artist %d", id))
		artist = s.artistStore.Retrieve(ctx, tx, id)

		//
		// retrieve related musics
		//

		artist.Musics = s.musicStore.List(ctx, tx, &model.MusicFilter{ArtistId: artist.Id})
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] retrieve artist %d", id), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] retrieve artist %d", id))
	return artist, nil
}

// //////////////////////////////////////////////////
// update music

func (s *artistService) UpdateArtist(ctx context.Context, artist *model.MusicArtist) (*model.MusicArtist, error) {

	//
	// pre-validate
	//

	if err := artist.Validate(nil); err != nil {
		return nil, err
	}

	//
	// download remote files
	//

	s.DownloadRemoteFiles(ctx, artist)

	//
	// update
	//

	var updated *model.MusicArtist
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {
		updated = s.artistStore.Update(ctx, tx, artist)

		if err := updated.Validate(s.imageFileValidator); err != nil {
			panic(err)
		}

		updated.Musics = s.musicStore.List(ctx, tx, &model.MusicFilter{ArtistId: artist.Id})
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] update artist %d - %s", artist.Id, artist.Name), zap.Object("artist", artist), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] update artist %d - %s", updated.Id, updated.Name), zap.Object("artist", updated))
	return updated, nil
}

// //////////////////////////////////////////////////
// delete music

func (s *artistService) DeleteArtist(ctx context.Context, id model.MusicArtistId) error {

	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// check if artist is used
		//

		if used := s.musicStore.IsArtistUsed(ctx, tx, id); used {
			panic(model.ErrArtistUsed)
		}

		//
		// delete artist
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] delete artist %d", id))
		s.artistStore.Delete(ctx, tx, id)
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] delete artist %d", id), zap.Error(err))
		return err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] delete artist %d", id))
	return nil
}

// //////////////////////////////////////////////////
// download remote files

func (s *artistService) DownloadRemoteFiles(ctx context.Context, artist *model.MusicArtist) {
	downloadArtistImage(s.logger, s.downloadClient, artist)
}

func downloadArtistImage(logger *zap.Logger, downloadClient client.DownloadClient, artist *model.MusicArtist) {
	if artist == nil {
		return
	}
	if artist.ImgUrl.IsRemote() {
		fileName := artist.GetImageFileName()
		url := artist.ImgUrl
		err := downloadClient.DownloadImage(url, fileName)
		if err != nil {
			logger.Info(fmt.Sprintf("[ KO ] download image %s <<< %q", fileName, url), zap.Error(err))
		} else {
			logger.Info(fmt.Sprintf("[ OK ] download image %s <<< %q", fileName, url))
			artist.ImgUrl = fileName
		}
	}
}
