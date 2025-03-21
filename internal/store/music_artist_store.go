package store

import (
	"context"
	"database/sql"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// music artist store

type MusicArtistStore interface {
	List(ctx context.Context, tx *sql.Tx, filter *model.MusicArtistFilter) []*model.MusicArtist
	Create(ctx context.Context, tx *sql.Tx, music *model.MusicArtist) *model.MusicArtist
	Retrieve(ctx context.Context, tx *sql.Tx, id model.MusicArtistId) *model.MusicArtist
	SearchByDeezerId(ctx context.Context, tx *sql.Tx, deezerId model.DeezerArtistId) *model.MusicArtist
	SearchByName(ctx context.Context, tx *sql.Tx, name string) *model.MusicArtist
	Update(ctx context.Context, tx *sql.Tx, music *model.MusicArtist) *model.MusicArtist
	Delete(ctx context.Context, tx *sql.Tx, id model.MusicArtistId)
}

func NewMusicArtistStore(logger *zap.Logger) MusicArtistStore {
	return &musicArtistStore{
		SqlTable: util.NewSqlTable[MusicArtistRow](logger, "music_artist", model.ErrMusicArtistNotFound),
	}
}

type musicArtistStore struct {
	util.SqlTable[MusicArtistRow]
	util.SqlEncoder[model.MusicArtist, MusicArtistRow]
	util.SqlDecoder[MusicArtistRow, model.MusicArtist]
}

// //////////////////////////////////////////////////
// row

type MusicArtistRow struct {
	Id       int64  `sql:"id,auto-generated"`
	DeezerId int64  `sql:"deezer_id"`
	Name     string `sql:"name"`
	ImgUrl   string `sql:"img_url"`
}

func (s *musicArtistStore) EncodeRow(obj *model.MusicArtist) *MusicArtistRow {
	return &MusicArtistRow{
		Id:       int64(obj.Id),
		DeezerId: int64(obj.DeezerId),
		Name:     obj.Name,
		ImgUrl:   string(obj.ImgUrl),
	}
}

func (s *musicArtistStore) DecodeRow(row *MusicArtistRow) *model.MusicArtist {
	if row == nil {
		return nil
	}
	return &model.MusicArtist{
		Id:       model.MusicArtistId(row.Id),
		DeezerId: model.DeezerArtistId(row.DeezerId),
		Name:     row.Name,
		ImgUrl:   model.Url(row.ImgUrl),
	}
}

// //////////////////////////////////////////////////
// list

func (s *musicArtistStore) List(ctx context.Context, tx *sql.Tx, filter *model.MusicArtistFilter) []*model.MusicArtist {
	return util.Convert(s.ListRows(ctx, tx, s.whereClause(filter)), s.DecodeRow)
}

func (s *musicArtistStore) whereClause(filter *model.MusicArtistFilter) util.SqlWhereClause {
	wc := util.NewSqlWhereClause()
	if filter != nil {
		if filter.Name != "" {
			wc.WithCondition("name LIKE '%' || $_ || '%'", filter.Name)
		}
		if filter.Limit > 0 {
			wc.WithLimit(filter.Limit)
		}
	}
	return wc
}

// //////////////////////////////////////////////////
// create

func (s *musicArtistStore) Create(ctx context.Context, tx *sql.Tx, obj *model.MusicArtist) *model.MusicArtist {
	return s.DecodeRow(s.InsertRow(ctx, tx, s.EncodeRow(obj)))
}

// //////////////////////////////////////////////////
// retrieve

func (s *musicArtistStore) Retrieve(ctx context.Context, tx *sql.Tx, id model.MusicArtistId) *model.MusicArtist {
	row, err := s.SelectRow(ctx, tx, s.matchingId(id))
	if err != nil {
		panic(err)
	}
	return s.DecodeRow(row)
}

// //////////////////////////////////////////////////
// retrieve by deezer id

func (s *musicArtistStore) SearchByDeezerId(ctx context.Context, tx *sql.Tx, deezerId model.DeezerArtistId) *model.MusicArtist {
	row, _ := s.SelectRow(ctx, tx, util.NewSqlCondition("deezer_id = $_", deezerId))
	return s.DecodeRow(row)
}

// //////////////////////////////////////////////////
// retrieve by deezer id

func (s *musicArtistStore) SearchByName(ctx context.Context, tx *sql.Tx, name string) *model.MusicArtist {
	row, _ := s.SelectRow(ctx, tx, util.NewSqlCondition("name = $_", name))
	return s.DecodeRow(row)
}

// //////////////////////////////////////////////////
// update

func (s *musicArtistStore) Update(ctx context.Context, tx *sql.Tx, obj *model.MusicArtist) *model.MusicArtist {
	return s.DecodeRow(s.UpdateRow(ctx, tx, s.EncodeRow(obj), s.matchingId(obj.Id)))
}

// //////////////////////////////////////////////////
// delete

func (s *musicArtistStore) Delete(ctx context.Context, tx *sql.Tx, id model.MusicArtistId) {
	s.DeleteRows(ctx, tx, s.matchingId(id))
}

// //////////////////////////////////////////////////
// where clause

func (s *musicArtistStore) matchingId(id model.MusicArtistId) util.SqlWhereClause {
	return util.NewSqlCondition("id = $_", id)
}
