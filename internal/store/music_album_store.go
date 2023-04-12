package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gre-ory/amnezic-go/internal/model"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// music album store

type MusicAlbumStore interface {
	Create(ctx context.Context, music *model.MusicAlbum) (*model.MusicAlbum, error)
	Retrieve(ctx context.Context, id model.MusicAlbumId) (*model.MusicAlbum, error)
	RetrieveByDeezerId(ctx context.Context, deezerId model.DeezerAlbumId) (*model.MusicAlbum, error)
	Update(ctx context.Context, music *model.MusicAlbum) (*model.MusicAlbum, error)
	Delete(ctx context.Context, id model.MusicAlbumId) error
}

func NewMusicAlbumStore(logger *zap.Logger, db *sql.DB) MusicAlbumStore {
	return &musicAlbumStore{
		logger:      logger,
		db:          db,
		table:       "music_album",
		columns:     "id,deezer_id,name,img_url",
		errNotFound: model.ErrMusicAlbumNotFound,
	}
}

type musicAlbumStore struct {
	logger      *zap.Logger
	db          *sql.DB
	table       string
	columns     string
	errNotFound error
}

// //////////////////////////////////////////////////
// adapter

func (s *musicAlbumStore) toModel(row *sql.Rows) *model.MusicAlbum {
	var id int64
	var deezerAlbumId int64
	var name string
	var imgUrl string
	row.Scan(&id, &deezerAlbumId, &name, &imgUrl)
	return &model.MusicAlbum{
		Id:       model.MusicAlbumId(id),
		DeezerId: model.DeezerAlbumId(deezerAlbumId),
		Name:     name,
		ImgUrl:   imgUrl,
	}
}

// //////////////////////////////////////////////////
// create

func (s *musicAlbumStore) Create(ctx context.Context, album *model.MusicAlbum) (*model.MusicAlbum, error) {

	query := fmt.Sprintf("INSERT INTO %s (deezer_id,name,img_url) VALUES ($1,$2,$3) RETURNING %s", s.table, s.columns)
	args := []any{album.DeezerId, album.Name, album.ImgUrl}
	s.logger.Info(fmt.Sprintf("[DEBUG] query: %s, args: %#v", query, args))

	statement, err := s.db.Prepare(query) // to avoid SQL injection
	if err != nil {
		return nil, err
	}

	rows, err := statement.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return s.toModel(rows), nil
	}

	return nil, s.errNotFound

}

// //////////////////////////////////////////////////
// retrieve

func (s *musicAlbumStore) Retrieve(ctx context.Context, id model.MusicAlbumId) (*model.MusicAlbum, error) {

	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", s.columns, s.table)
	args := []any{id}
	s.logger.Info(fmt.Sprintf("[DEBUG] query: %s, args: %#v", query, args))

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return s.toModel(rows), nil
	}

	return nil, s.errNotFound
}

// //////////////////////////////////////////////////
// retrieve by deezer id

func (s *musicAlbumStore) RetrieveByDeezerId(ctx context.Context, deezerId model.DeezerAlbumId) (*model.MusicAlbum, error) {

	query := fmt.Sprintf("SELECT %s FROM %s WHERE deezer_id = $1", s.columns, s.table)
	args := []any{deezerId}
	s.logger.Info(fmt.Sprintf("[DEBUG] query: %s, args: %#v", query, args))

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return s.toModel(rows), nil
	}

	return nil, s.errNotFound
}

// //////////////////////////////////////////////////
// update

func (s *musicAlbumStore) Update(ctx context.Context, music *model.MusicAlbum) (*model.MusicAlbum, error) {
	return nil, model.ErrNotImplemented
}

// //////////////////////////////////////////////////
// delete

func (s *musicAlbumStore) Delete(ctx context.Context, id model.MusicAlbumId) error {
	return model.ErrNotImplemented
}
