package store

import (
	"context"
	"database/sql"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// themeQuestion store

type ThemeQuestionStore interface {
	Create(ctx context.Context, tx *sql.Tx, themeQuestion *model.ThemeQuestion) *model.ThemeQuestion
	Retrieve(ctx context.Context, tx *sql.Tx, id model.ThemeQuestionId) *model.ThemeQuestion
	Update(ctx context.Context, tx *sql.Tx, themeQuestion *model.ThemeQuestion) *model.ThemeQuestion
	Delete(ctx context.Context, tx *sql.Tx, filter *model.ThemeQuestionFilter)
	List(ctx context.Context, tx *sql.Tx, filter *model.ThemeQuestionFilter) []*model.ThemeQuestion
	CountByTheme(ctx context.Context, tx *sql.Tx) map[model.ThemeId]int
	IsMusicInTheme(ctx context.Context, tx *sql.Tx, themeId model.ThemeId, musicId model.MusicId) bool
	IsMusicUsed(ctx context.Context, tx *sql.Tx, musicId model.MusicId) bool
}

func NewThemeQuestionStore(logger *zap.Logger) ThemeQuestionStore {
	return &themeQuestionStore{
		SqlTable: util.NewSqlTable[ThemeQuestionRow](logger, ThemeQuestionTable, model.ErrThemeQuestionNotFound),
	}
}

type themeQuestionStore struct {
	util.SqlTable[ThemeQuestionRow]
	util.SqlEncoder[model.ThemeQuestion, ThemeQuestionRow]
	util.SqlDecoder[ThemeQuestionRow, model.ThemeQuestion]
}

// //////////////////////////////////////////////////
// table

const ThemeQuestionTable = "theme_question"

// //////////////////////////////////////////////////
// row

type ThemeQuestionRow struct {
	Id      int64  `sql:"id,auto-generated"`
	ThemeId int64  `sql:"theme_id"`
	MusicId int64  `sql:"music_id"`
	Text    string `sql:"text"`
	Hint    string `sql:"hint"`
}

func (s *themeQuestionStore) EncodeRow(obj *model.ThemeQuestion) *ThemeQuestionRow {
	return &ThemeQuestionRow{
		Id:      int64(obj.Id),
		ThemeId: int64(obj.ThemeId),
		MusicId: int64(obj.MusicId),
		Text:    obj.Text,
		Hint:    obj.Hint,
	}
}

func (s *themeQuestionStore) DecodeRow(row *ThemeQuestionRow) *model.ThemeQuestion {
	if row == nil {
		return nil
	}
	return &model.ThemeQuestion{
		Id:      model.ThemeQuestionId(row.Id),
		ThemeId: model.ThemeId(row.ThemeId),
		MusicId: model.MusicId(row.MusicId),
		Text:    row.Text,
		Hint:    row.Hint,
	}
}

// //////////////////////////////////////////////////
// create

func (s *themeQuestionStore) Create(ctx context.Context, tx *sql.Tx, obj *model.ThemeQuestion) *model.ThemeQuestion {
	return s.DecodeRow(s.InsertRow(ctx, tx, s.EncodeRow(obj)))
}

// //////////////////////////////////////////////////
// retrieve

func (s *themeQuestionStore) Retrieve(ctx context.Context, tx *sql.Tx, id model.ThemeQuestionId) *model.ThemeQuestion {
	row, err := s.SelectRow(ctx, tx, s.matchingId(id))
	if err != nil {
		panic(err)
	}
	return s.DecodeRow(row)
}

// //////////////////////////////////////////////////
// update

func (s *themeQuestionStore) Update(ctx context.Context, tx *sql.Tx, obj *model.ThemeQuestion) *model.ThemeQuestion {
	return s.DecodeRow(s.UpdateRow(ctx, tx, s.EncodeRow(obj), s.matchingId(obj.Id)))
}

// //////////////////////////////////////////////////
// delete

func (s *themeQuestionStore) Delete(ctx context.Context, tx *sql.Tx, filter *model.ThemeQuestionFilter) {
	s.DeleteRows(ctx, tx, s.whereClause(filter))
}

// //////////////////////////////////////////////////
// list

func (s *themeQuestionStore) List(ctx context.Context, tx *sql.Tx, filter *model.ThemeQuestionFilter) []*model.ThemeQuestion {
	return util.Convert(s.ListRows(ctx, tx, s.whereClause(filter)), s.DecodeRow)
}

// //////////////////////////////////////////////////
// list

func (s *themeQuestionStore) CountByTheme(ctx context.Context, tx *sql.Tx) map[model.ThemeId]int {
	result := make(map[model.ThemeId]int, 0)
	util.SqlScan(
		util.SqlQuery(ctx, tx, "SELECT theme_id, count(1) AS count FROM "+ThemeQuestionTable+" GROUP BY theme_id"),
		func(rows *sql.Rows) {
			var themeId int64
			var count int
			rows.Scan(&themeId, &count)
			result[model.ThemeId(themeId)] = count
		},
	)
	return result
}

// //////////////////////////////////////////////////
// is music used

func (s *themeQuestionStore) IsMusicInTheme(ctx context.Context, tx *sql.Tx, themeId model.ThemeId, musicId model.MusicId) bool {
	return s.ExistsRow(ctx, tx,
		util.NewSqlWhereClause().
			WithCondition("theme_id = %s", themeId).
			WithCondition("music_id = %s", musicId),
	)
}

// //////////////////////////////////////////////////
// is music used

func (s *themeQuestionStore) IsMusicUsed(ctx context.Context, tx *sql.Tx, musicId model.MusicId) bool {
	return s.ExistsRow(ctx, tx,
		util.NewSqlCondition("music_id = %s", musicId),
	)
}

// //////////////////////////////////////////////////
// where clause

func (s *themeQuestionStore) matchingId(id model.ThemeQuestionId) util.SqlWhereClause {
	return util.NewSqlCondition("id = %s", id)
}

func (s *themeQuestionStore) whereClause(filter *model.ThemeQuestionFilter) util.SqlWhereClause {
	wc := util.NewSqlWhereClause()
	if filter != nil {
		if filter.ThemeQuestionId != 0 {
			wc.WithCondition("id = %s", filter.ThemeQuestionId)
		}
		if len(filter.ThemeIds) == 1 {
			wc.WithCondition("theme_id = %s", filter.ThemeIds[0])
		} else if len(filter.ThemeIds) > 1 {
			wc.WithCondition("theme_id IN (%s)", util.Join(filter.ThemeIds, util.IntToStr[model.ThemeId]))
		}
		if filter.MusicId != 0 {
			wc.WithCondition("music_id = %s", filter.MusicId)
		}
		if filter.Random {
			wc.WithRandomOrder()
		}
		if filter.Limit != 0 {
			wc.WithLimit(filter.Limit)
		}
	}
	return wc
}
