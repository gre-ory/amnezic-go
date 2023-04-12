package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gre-ory/amnezic-go/internal/model"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// music artist store

type MusicArtistStore interface {
	Create(ctx context.Context, music *model.MusicArtist) (*model.MusicArtist, error)
	Retrieve(ctx context.Context, id model.MusicArtistId) (*model.MusicArtist, error)
	RetrieveByDeezerId(ctx context.Context, deezerId model.DeezerArtistId) (*model.MusicArtist, error)
	Update(ctx context.Context, music *model.MusicArtist) (*model.MusicArtist, error)
	Delete(ctx context.Context, id model.MusicArtistId) error
}

func NewMusicArtistStore(logger *zap.Logger, db *sql.DB) MusicArtistStore {
	return &musicArtistStore{
		logger:      logger,
		db:          db,
		table:       "music_artist",
		columns:     "id,deezer_id,name,img_url",
		errNotFound: model.ErrMusicArtistNotFound,
	}
}

type musicArtistStore struct {
	logger      *zap.Logger
	db          *sql.DB
	table       string
	columns     string
	errNotFound error
}

// //////////////////////////////////////////////////
// adapter

func (s *musicArtistStore) toModel(row *sql.Rows) *model.MusicArtist {
	var id int64
	var deezerArtistId int64
	var name string
	var imgUrl string
	row.Scan(&id, &deezerArtistId, &name, &imgUrl)
	return &model.MusicArtist{
		Id:       model.MusicArtistId(id),
		DeezerId: model.DeezerArtistId(deezerArtistId),
		Name:     name,
		ImgUrl:   imgUrl,
	}
}

// //////////////////////////////////////////////////
// create

func (s *musicArtistStore) Create(ctx context.Context, artist *model.MusicArtist) (*model.MusicArtist, error) {

	query := fmt.Sprintf("INSERT INTO %s (deezer_id,name,img_url) VALUES ($1,$2,$3) RETURNING %s", s.table, s.columns)
	args := []any{artist.DeezerId, artist.Name, artist.ImgUrl}
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

func (s *musicArtistStore) Retrieve(ctx context.Context, id model.MusicArtistId) (*model.MusicArtist, error) {

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

func (s *musicArtistStore) RetrieveByDeezerId(ctx context.Context, deezerId model.DeezerArtistId) (*model.MusicArtist, error) {

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

func (s *musicArtistStore) Update(ctx context.Context, music *model.MusicArtist) (*model.MusicArtist, error) {
	return nil, model.ErrNotImplemented
}

// //////////////////////////////////////////////////
// delete

func (s *musicArtistStore) Delete(ctx context.Context, id model.MusicArtistId) error {
	return model.ErrNotImplemented
}
