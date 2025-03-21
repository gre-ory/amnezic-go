package store

import (
	"context"
	"database/sql"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// music store

type MusicStore interface {
	List(ctx context.Context, tx *sql.Tx, filter *model.MusicFilter) []*model.Music
	Create(ctx context.Context, tx *sql.Tx, music *model.Music) *model.Music
	Retrieve(ctx context.Context, tx *sql.Tx, id model.MusicId) *model.Music
	SearchByDeezerId(ctx context.Context, tx *sql.Tx, deezerId model.DeezerMusicId) *model.Music
	SearchByName(ctx context.Context, tx *sql.Tx, name string) *model.Music
	Update(ctx context.Context, tx *sql.Tx, music *model.Music) *model.Music
	Delete(ctx context.Context, tx *sql.Tx, id model.MusicId)
	IsAlbumUsed(ctx context.Context, tx *sql.Tx, albumId model.MusicAlbumId) bool
	IsArtistUsed(ctx context.Context, tx *sql.Tx, artistId model.MusicArtistId) bool
}

func NewMusicStore(logger *zap.Logger) MusicStore {
	return &musicStore{
		SqlTable: util.NewSqlTable[MusicRow](logger, "music", model.ErrMusicNotFound),
	}
}

type musicStore struct {
	util.SqlTable[MusicRow]
	util.SqlEncoder[model.Music, MusicRow]
	util.SqlDecoder[MusicRow, model.Music]
}

// //////////////////////////////////////////////////
// row

type MusicRow struct {
	Id       int64  `sql:"id,auto-generated"`
	DeezerId int64  `sql:"deezer_id"`
	ArtistId int64  `sql:"artist_id"`
	AlbumId  int64  `sql:"album_id"`
	Name     string `sql:"name"`
	Mp3Url   string `sql:"mp3_url"`
}

func (s *musicStore) EncodeRow(obj *model.Music) *MusicRow {
	return &MusicRow{
		Id:       int64(obj.Id),
		DeezerId: int64(obj.DeezerId),
		Name:     obj.Name,
		Mp3Url:   string(obj.Mp3Url),
		ArtistId: int64(obj.ArtistId),
		AlbumId:  int64(obj.AlbumId),
	}
}

func (s *musicStore) DecodeRow(row *MusicRow) *model.Music {
	if row == nil {
		return nil
	}
	return &model.Music{
		Id:       model.MusicId(row.Id),
		DeezerId: model.DeezerMusicId(row.DeezerId),
		Name:     row.Name,
		Mp3Url:   model.Url(row.Mp3Url),
		ArtistId: model.MusicArtistId(row.ArtistId),
		AlbumId:  model.MusicAlbumId(row.AlbumId),
	}
}

// //////////////////////////////////////////////////
// list

func (s *musicStore) List(ctx context.Context, tx *sql.Tx, filter *model.MusicFilter) []*model.Music {
	return util.Convert(s.ListRows(ctx, tx, s.whereClause(filter)), s.DecodeRow)
}

func (s *musicStore) whereClause(filter *model.MusicFilter) util.SqlWhereClause {
	wc := util.NewSqlWhereClause()
	if filter != nil {
		if filter.Name != "" {
			wc.WithCondition("name LIKE '%' || $_ || '%'", filter.Name)
		}
		if filter.ArtistId > 0 {
			wc.WithCondition("artist_id = $_", filter.ArtistId)
		}
		if filter.AlbumId > 0 {
			wc.WithCondition("artist_id = $_", filter.ArtistId)
		}
		if filter.Limit > 0 {
			wc.WithLimit(filter.Limit)
		}
	}
	return wc
}

// //////////////////////////////////////////////////
// create

func (s *musicStore) Create(ctx context.Context, tx *sql.Tx, obj *model.Music) *model.Music {
	return s.DecodeRow(s.InsertRow(ctx, tx, s.EncodeRow(obj)))
}

// //////////////////////////////////////////////////
// retrieve

func (s *musicStore) Retrieve(ctx context.Context, tx *sql.Tx, id model.MusicId) *model.Music {
	row, err := s.SelectRow(ctx, tx, s.matchingId(id))
	if err != nil {
		panic(err)
	}
	return s.DecodeRow(row)
}

// //////////////////////////////////////////////////
// search by deezer id

func (s *musicStore) SearchByDeezerId(ctx context.Context, tx *sql.Tx, deezerId model.DeezerMusicId) *model.Music {
	row, _ := s.SelectRow(ctx, tx, util.NewSqlCondition("deezer_id = $_", deezerId))
	return s.DecodeRow(row)
}

// //////////////////////////////////////////////////
// search by deezer id

func (s *musicStore) SearchByName(ctx context.Context, tx *sql.Tx, name string) *model.Music {
	row, _ := s.SelectRow(ctx, tx, util.NewSqlCondition("name = $_", name))
	return s.DecodeRow(row)
}

// //////////////////////////////////////////////////
// update

func (s *musicStore) Update(ctx context.Context, tx *sql.Tx, obj *model.Music) *model.Music {
	return s.DecodeRow(s.UpdateRow(ctx, tx, s.EncodeRow(obj), s.matchingId(obj.Id)))
}

// //////////////////////////////////////////////////
// delete

func (s *musicStore) Delete(ctx context.Context, tx *sql.Tx, id model.MusicId) {
	s.DeleteRows(ctx, tx, s.matchingId(id))
}

// //////////////////////////////////////////////////
// album usage

func (s *musicStore) IsAlbumUsed(ctx context.Context, tx *sql.Tx, albumId model.MusicAlbumId) bool {
	return s.ExistsRow(ctx, tx, util.NewSqlCondition("album_id = $_", albumId))
}

// //////////////////////////////////////////////////
// artist usage

func (s *musicStore) IsArtistUsed(ctx context.Context, tx *sql.Tx, artistId model.MusicArtistId) bool {
	return s.ExistsRow(ctx, tx, util.NewSqlCondition("artist_id = $_", artistId))
}

// //////////////////////////////////////////////////
// where clause

func (s *musicStore) matchingId(id model.MusicId) util.SqlWhereClause {
	return util.NewSqlCondition("id = $_", id)
}
