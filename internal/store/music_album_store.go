package store

import (
	"context"
	"database/sql"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// music album store

type MusicAlbumStore interface {
	Create(ctx context.Context, tx *sql.Tx, music *model.MusicAlbum) *model.MusicAlbum
	Retrieve(ctx context.Context, tx *sql.Tx, id model.MusicAlbumId) *model.MusicAlbum
	SearchByDeezerId(ctx context.Context, tx *sql.Tx, deezerId model.DeezerAlbumId) *model.MusicAlbum
	Update(ctx context.Context, tx *sql.Tx, music *model.MusicAlbum) *model.MusicAlbum
	Delete(ctx context.Context, tx *sql.Tx, id model.MusicAlbumId)
}

func NewMusicAlbumStore(logger *zap.Logger) MusicAlbumStore {
	return &musicAlbumStore{
		SqlTable: util.NewSqlTable[MusicAlbumRow](logger, "music_album", model.ErrMusicAlbumNotFound),
	}
}

type musicAlbumStore struct {
	util.SqlTable[MusicAlbumRow]
	util.SqlEncoder[model.MusicArtist, MusicArtistRow]
	util.SqlDecoder[MusicArtistRow, model.MusicArtist]
}

// //////////////////////////////////////////////////
// row

type MusicAlbumRow struct {
	Id       int64  `sql:"id,auto-generated"`
	DeezerId int64  `sql:"deezer_id"`
	Name     string `sql:"name"`
	ImgUrl   string `sql:"img_url"`
}

func (s *musicAlbumStore) EncodeRow(obj *model.MusicAlbum) *MusicAlbumRow {
	return &MusicAlbumRow{
		Id:       int64(obj.Id),
		DeezerId: int64(obj.DeezerId),
		Name:     obj.Name,
		ImgUrl:   obj.ImgUrl,
	}
}

func (s *musicAlbumStore) DecodeRow(row *MusicAlbumRow) *model.MusicAlbum {
	if row == nil {
		return nil
	}
	return &model.MusicAlbum{
		Id:       model.MusicAlbumId(row.Id),
		DeezerId: model.DeezerAlbumId(row.DeezerId),
		Name:     row.Name,
		ImgUrl:   row.ImgUrl,
	}
}

// //////////////////////////////////////////////////
// create

func (s *musicAlbumStore) Create(ctx context.Context, tx *sql.Tx, obj *model.MusicAlbum) *model.MusicAlbum {
	return s.DecodeRow(s.InsertRow(ctx, tx, s.EncodeRow(obj)))
}

// //////////////////////////////////////////////////
// retrieve

func (s *musicAlbumStore) Retrieve(ctx context.Context, tx *sql.Tx, id model.MusicAlbumId) *model.MusicAlbum {
	row, err := s.SelectRow(ctx, tx, "WHERE id = %s", id)
	if err != nil {
		panic(err)
	}
	return s.DecodeRow(row)
}

// //////////////////////////////////////////////////
// retrieve by deezer id

func (s *musicAlbumStore) SearchByDeezerId(ctx context.Context, tx *sql.Tx, deezerId model.DeezerAlbumId) *model.MusicAlbum {
	row, _ := s.SelectRow(ctx, tx, "WHERE deezer_id = %s", deezerId)
	return s.DecodeRow(row)
}

// //////////////////////////////////////////////////
// update

func (s *musicAlbumStore) Update(ctx context.Context, tx *sql.Tx, obj *model.MusicAlbum) *model.MusicAlbum {
	return s.DecodeRow(s.UpdateRow(ctx, tx, s.EncodeRow(obj), "WHERE id = %s", obj.Id))
}

// //////////////////////////////////////////////////
// delete

func (s *musicAlbumStore) Delete(ctx context.Context, tx *sql.Tx, id model.MusicAlbumId) {
	s.DeleteRow(ctx, tx, "WHERE id = %s", id)
}
