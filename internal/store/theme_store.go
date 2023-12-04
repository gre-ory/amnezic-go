package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// theme store

type ThemeStore interface {
	Create(ctx context.Context, tx *sql.Tx, theme *model.Theme) *model.Theme
	Retrieve(ctx context.Context, tx *sql.Tx, id model.ThemeId) *model.Theme
	Update(ctx context.Context, tx *sql.Tx, theme *model.Theme) *model.Theme
	Delete(ctx context.Context, tx *sql.Tx, filter *model.ThemeFilter)
	List(ctx context.Context, tx *sql.Tx, filter *model.ThemeFilter) []*model.Theme
}

func NewThemeStore(logger *zap.Logger) ThemeStore {
	return &themeStore{
		SqlTable: util.NewSqlTable[ThemeRow](logger, ThemeTable, model.ErrThemeNotFound),
	}
}

type themeStore struct {
	util.SqlTable[ThemeRow]
	util.SqlEncoder[model.Theme, ThemeRow]
	util.SqlDecoder[ThemeRow, model.Theme]
}

// //////////////////////////////////////////////////
// table

const ThemeTable = "theme"

// //////////////////////////////////////////////////
// row

type ThemeRow struct {
	Id     int64  `sql:"id,auto-generated"`
	Title  string `sql:"title"`
	ImgUrl string `sql:"img_url"`
	Labels string `sql:"labels"`
}

func (s *themeStore) EncodeRow(obj *model.Theme) *ThemeRow {
	return &ThemeRow{
		Id:     int64(obj.Id),
		Title:  obj.Title,
		ImgUrl: obj.ImgUrl,
		Labels: s.EncodeLabels(obj.Labels),
	}
}

func (s *themeStore) EncodeLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return ""
	}
	return util.JoinMap(labels, func(key string, value string) string {
		return fmt.Sprintf("%s=%s", key, value)
	})
}

func (s *themeStore) DecodeRow(row *ThemeRow) *model.Theme {
	if row == nil {
		return nil
	}
	return &model.Theme{
		Id:     model.ThemeId(row.Id),
		Title:  row.Title,
		ImgUrl: row.ImgUrl,
		Labels: s.DecodeLabels(row.Labels),
	}
}

func (s *themeStore) DecodeLabels(labels string) map[string]string {
	if labels == "" {
		return nil
	}
	return util.ConvertToMap(strings.Split(labels, ","), func(item string) (string, string) {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) > 1 {
			return parts[0], parts[1]
		}
		return parts[0], ""
	})
}

// //////////////////////////////////////////////////
// create

func (s *themeStore) Create(ctx context.Context, tx *sql.Tx, obj *model.Theme) *model.Theme {
	return s.DecodeRow(s.InsertRow(ctx, tx, s.EncodeRow(obj)))
}

// //////////////////////////////////////////////////
// retrieve

func (s *themeStore) Retrieve(ctx context.Context, tx *sql.Tx, id model.ThemeId) *model.Theme {
	row, err := s.SelectRow(ctx, tx, s.matchingId(id))
	if err != nil {
		panic(err)
	}
	return s.DecodeRow(row)
}

// //////////////////////////////////////////////////
// update

func (s *themeStore) Update(ctx context.Context, tx *sql.Tx, obj *model.Theme) *model.Theme {
	return s.DecodeRow(s.UpdateRow(ctx, tx, s.EncodeRow(obj), s.matchingId(obj.Id)))
}

// //////////////////////////////////////////////////
// delete

func (s *themeStore) Delete(ctx context.Context, tx *sql.Tx, filter *model.ThemeFilter) {
	s.DeleteRows(ctx, tx, s.whereClause(filter))
}

// //////////////////////////////////////////////////
// list

func (s *themeStore) List(ctx context.Context, tx *sql.Tx, filter *model.ThemeFilter) []*model.Theme {
	return util.Convert(s.ListRows(ctx, tx, s.whereClause(filter)), s.DecodeRow)
}

// //////////////////////////////////////////////////
// where clause

func (s *themeStore) matchingId(id model.ThemeId) util.SqlWhereClause {
	return util.NewSqlCondition("id = %s", id)
}

func (s *themeStore) whereClause(filter *model.ThemeFilter) util.SqlWhereClause {
	wc := util.NewSqlWhereClause()
	if filter != nil {
		if filter.ThemeId != 0 {
			wc.WithCondition("id = %s", filter.ThemeId)
		}
	}
	return wc
}
