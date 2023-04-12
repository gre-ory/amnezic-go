package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gre-ory/amnezic-go/internal/model"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// music store

type MusicStore interface {
	Create(ctx context.Context, music *model.Music) (*model.Music, error)
	Retrieve(ctx context.Context, id model.MusicId) (*model.Music, error)
	RetrieveByDeezerId(ctx context.Context, deezerId model.DeezerMusicId) (*model.Music, error)
	Update(ctx context.Context, music *model.Music) (*model.Music, error)
	Delete(ctx context.Context, id model.MusicId) error
	IsAlbumUsed(ctx context.Context, albumId model.MusicAlbumId) (bool, error)
	IsArtistUsed(ctx context.Context, artistId model.MusicArtistId) (bool, error)
}

func NewMusicStore(logger *zap.Logger, db *sql.DB) MusicStore {
	return &musicStore{
		logger:      logger,
		db:          db,
		table:       "music",
		columns:     "id,deezer_id,artist_id,album_id,name,mp3_url",
		errNotFound: model.ErrMusicNotFound,
	}
}

type musicStore struct {
	logger      *zap.Logger
	db          *sql.DB
	table       string
	columns     string
	errNotFound error
}

// //////////////////////////////////////////////////
// adapter

func (s *musicStore) toModel(row *sql.Rows) *model.Music {
	var id int64
	var deezerMusicId int64
	var artistId int64
	var albumId int64
	var name string
	var mp3Url string
	row.Scan(&id, &deezerMusicId, &artistId, &albumId, &name, &mp3Url)
	return &model.Music{
		Id:       model.MusicId(id),
		DeezerId: model.DeezerMusicId(deezerMusicId),
		ArtistId: model.MusicArtistId(artistId),
		AlbumId:  model.MusicAlbumId(albumId),
		Name:     name,
		Mp3Url:   mp3Url,
	}
}

// //////////////////////////////////////////////////
// create

func (s *musicStore) Create(ctx context.Context, music *model.Music) (*model.Music, error) {

	query := fmt.Sprintf("INSERT INTO %s (deezer_id,artist_id,album_id,name,mp3_url) VALUES ($1,$2,$3,$4,$5) RETURNING %s", s.table, s.columns)
	args := []any{music.DeezerId, music.ArtistId, music.AlbumId, music.Name, music.Mp3Url}
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

func (s *musicStore) Retrieve(ctx context.Context, id model.MusicId) (*model.Music, error) {

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

func (s *musicStore) RetrieveByDeezerId(ctx context.Context, deezerId model.DeezerMusicId) (*model.Music, error) {

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

func (s *musicStore) Update(ctx context.Context, music *model.Music) (*model.Music, error) {
	return nil, model.ErrNotImplemented
}

// //////////////////////////////////////////////////
// delete

func (s *musicStore) Delete(ctx context.Context, id model.MusicId) error {
	return model.ErrNotImplemented
}

// //////////////////////////////////////////////////
// album usage

func (s *musicStore) IsAlbumUsed(ctx context.Context, albumId model.MusicAlbumId) (bool, error) {
	if albumId == 0 {
		return false, nil
	}

	return false, model.ErrNotImplemented
}

// //////////////////////////////////////////////////
// artist usage

func (s *musicStore) IsArtistUsed(ctx context.Context, artistId model.MusicArtistId) (bool, error) {
	if artistId == 0 {
		return false, nil
	}

	return false, model.ErrNotImplemented
}
