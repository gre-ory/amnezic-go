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
// album service

type AlbumService interface {
	ListAlbum(ctx context.Context, search *model.MusicAlbumFilter) ([]*model.MusicAlbum, error)
	GetAlbum(ctx context.Context, id model.MusicAlbumId) (*model.MusicAlbum, error)
	CreateAlbum(ctx context.Context, music *model.MusicAlbum) (*model.MusicAlbum, error)
	UpdateAlbum(ctx context.Context, music *model.MusicAlbum) (*model.MusicAlbum, error)
	DeleteAlbum(ctx context.Context, id model.MusicAlbumId) error
}

func NewAlbumService(logger *zap.Logger, downloadClient client.DownloadClient, db *sql.DB, albumStore store.MusicAlbumStore, musicStore store.MusicStore, imageFileValidator model.PathValidator) AlbumService {
	return &albumService{
		logger:             logger,
		downloadClient:     downloadClient,
		db:                 db,
		albumStore:         albumStore,
		musicStore:         musicStore,
		imageFileValidator: imageFileValidator,
	}
}

type albumService struct {
	logger             *zap.Logger
	downloadClient     client.DownloadClient
	db                 *sql.DB
	albumStore         store.MusicAlbumStore
	musicStore         store.MusicStore
	imageFileValidator model.PathValidator
}

// //////////////////////////////////////////////////
// list album

func (s *albumService) ListAlbum(ctx context.Context, filter *model.MusicAlbumFilter) ([]*model.MusicAlbum, error) {
	var albums []*model.MusicAlbum
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {
		albums = s.albumStore.List(ctx, tx, filter)
	})
	if err != nil {
		s.logger.Info("[ KO ] list album", zap.Error(err))
		return nil, err
	}
	s.logger.Info("[ OK ] list album")
	return albums, nil
}

// //////////////////////////////////////////////////
// create album

func (s *albumService) CreateAlbum(ctx context.Context, album *model.MusicAlbum) (*model.MusicAlbum, error) {

	var err error

	//
	// validate
	//

	if album == nil {
		return nil, model.ErrMissingAlbum
	}
	if err = album.Validate(s.imageFileValidator); err != nil {
		return nil, err
	}

	//
	// download remote files
	//

	s.DownloadRemoteFiles(ctx, album)

	//
	// create
	//

	err = util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// check if album exists
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve album from name %q", album.Name))
		other := s.albumStore.SearchByName(ctx, tx, album.Name)
		if other != nil {
			panic(model.ErrExistingAlbum)
		}

		//
		// create album
		//

		album = s.albumStore.Create(ctx, tx, album)
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] create album %q", album.Name), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] create album %q", album.Name))
	return album, nil

}

// //////////////////////////////////////////////////
// get album

func (s *albumService) GetAlbum(ctx context.Context, id model.MusicAlbumId) (*model.MusicAlbum, error) {

	var album *model.MusicAlbum
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// retrieve album
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve album %d", id))
		album = s.albumStore.Retrieve(ctx, tx, id)

		//
		// retrieve related musics
		//

		album.Musics = s.musicStore.List(ctx, tx, &model.MusicFilter{AlbumId: album.Id})
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] retrieve album %d", id), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] retrieve album %d", id))
	return album, nil
}

// //////////////////////////////////////////////////
// update music

func (s *albumService) UpdateAlbum(ctx context.Context, album *model.MusicAlbum) (*model.MusicAlbum, error) {

	//
	// pre-validate
	//

	if err := album.Validate(nil); err != nil {
		return nil, err
	}

	//
	// download remote files
	//

	s.DownloadRemoteFiles(ctx, album)

	//
	// update
	//

	var updated *model.MusicAlbum
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {
		updated = s.albumStore.Update(ctx, tx, album)

		if err := updated.Validate(s.imageFileValidator); err != nil {
			panic(err)
		}

		updated.Musics = s.musicStore.List(ctx, tx, &model.MusicFilter{AlbumId: album.Id})
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] update album %d - %s", album.Id, album.Name), zap.Object("album", album), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] update album %d - %s", updated.Id, updated.Name), zap.Object("album", updated))
	return updated, nil
}

// //////////////////////////////////////////////////
// delete music

func (s *albumService) DeleteAlbum(ctx context.Context, id model.MusicAlbumId) error {

	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// check if album is used
		//

		if used := s.musicStore.IsAlbumUsed(ctx, tx, id); used {
			panic(model.ErrAlbumUsed)
		}

		//
		// delete album
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] delete album %d", id))
		s.albumStore.Delete(ctx, tx, id)
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] delete album %d", id), zap.Error(err))
		return err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] delete album %d", id))
	return nil
}

// //////////////////////////////////////////////////
// download remote files

func (s *albumService) DownloadRemoteFiles(ctx context.Context, album *model.MusicAlbum) {
	downloadAlbumImage(s.logger, s.downloadClient, album)
}

func downloadAlbumImage(logger *zap.Logger, downloadClient client.DownloadClient, album *model.MusicAlbum) {
	if album == nil {
		return
	}
	if album.ImgUrl.IsRemote() {
		fileName := album.GetImageFileName()
		url := album.ImgUrl
		err := downloadClient.DownloadImage(url, fileName)
		if err != nil {
			logger.Info(fmt.Sprintf("[ KO ] download image %s <<< %q", fileName, url), zap.Error(err))
		} else {
			logger.Info(fmt.Sprintf("[ OK ] download image %s <<< %q", fileName, url))
			album.ImgUrl = fileName
		}
	}
}
