package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// user store

type UserStore interface {
	Create(ctx context.Context, tx *sql.Tx, user *model.User) *model.User
	Retrieve(ctx context.Context, tx *sql.Tx, id model.UserId) *model.User
	SearchByName(ctx context.Context, tx *sql.Tx, name string) *model.User
	Update(ctx context.Context, tx *sql.Tx, user *model.User) *model.User
	Delete(ctx context.Context, tx *sql.Tx, filter *model.UserFilter)
	List(ctx context.Context, tx *sql.Tx, filter *model.UserFilter) []*model.User
}

func NewUserStore(logger *zap.Logger) UserStore {
	return &userStore{
		SqlTable: util.NewSqlTable[UserRow](logger, UserTable, model.ErrUserNotFound),
	}
}

type userStore struct {
	util.SqlTable[UserRow]
	util.SqlEncoder[model.User, UserRow]
	util.SqlDecoder[UserRow, model.User]
}

// //////////////////////////////////////////////////
// table

const UserTable = "user"

// //////////////////////////////////////////////////
// row

type UserRow struct {
	Id          int64  `sql:"id,auto-generated"`
	Name        string `sql:"name"`
	Hash        string `sql:"hash"`
	Permissions string `sql:"permissions"`
}

func (s *userStore) EncodeRow(obj *model.User) *UserRow {
	return &UserRow{
		Id:          obj.Id.ToInt64(),
		Name:        obj.Name,
		Hash:        obj.Hash,
		Permissions: s.EncodePermissions(obj.Permissions),
	}
}

func (s *userStore) EncodePermissions(permissions []model.Permission) string {
	if len(permissions) == 0 {
		return ""
	}
	return util.Join(permissions, ",")
}

func (s *userStore) DecodeRow(row *UserRow) *model.User {
	if row == nil {
		return nil
	}
	return &model.User{
		Id:          model.UserId(row.Id),
		Name:        row.Name,
		Hash:        row.Hash,
		Permissions: s.DecodePermissions(row.Permissions),
	}
}

func (s *userStore) DecodePermissions(permissions string) []model.Permission {
	if permissions == "" {
		return nil
	}
	return util.Convert(strings.Split(permissions, ","), model.ToPermission)
}

// //////////////////////////////////////////////////
// create

func (s *userStore) Create(ctx context.Context, tx *sql.Tx, obj *model.User) *model.User {
	return s.DecodeRow(s.InsertRow(ctx, tx, s.EncodeRow(obj)))
}

// //////////////////////////////////////////////////
// retrieve

func (s *userStore) Retrieve(ctx context.Context, tx *sql.Tx, id model.UserId) *model.User {
	row, err := s.SelectRow(ctx, tx, s.matchingId(id))
	if err != nil {
		panic(err)
	}
	return s.DecodeRow(row)
}

// //////////////////////////////////////////////////
// search by name

func (s *userStore) SearchByName(ctx context.Context, tx *sql.Tx, name string) *model.User {
	row, err := s.SelectRow(ctx, tx, s.matchingName(name))
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil
		} else {
			panic(err)
		}
	}
	return s.DecodeRow(row)
}

// //////////////////////////////////////////////////
// update

func (s *userStore) Update(ctx context.Context, tx *sql.Tx, obj *model.User) *model.User {
	return s.DecodeRow(s.UpdateRow(ctx, tx, s.EncodeRow(obj), s.matchingId(obj.Id)))
}

// //////////////////////////////////////////////////
// delete

func (s *userStore) Delete(ctx context.Context, tx *sql.Tx, filter *model.UserFilter) {
	s.DeleteRows(ctx, tx, s.whereClause(filter))
}

// //////////////////////////////////////////////////
// list

func (s *userStore) List(ctx context.Context, tx *sql.Tx, filter *model.UserFilter) []*model.User {
	return util.Convert(s.ListRows(ctx, tx, s.whereClause(filter)), s.DecodeRow)
}

// //////////////////////////////////////////////////
// where clause

func (s *userStore) matchingId(id model.UserId) util.SqlWhereClause {
	return util.NewSqlCondition("id = $_", id)
}

func (s *userStore) matchingName(name string) util.SqlWhereClause {
	return util.NewSqlCondition("name = $_", name)
}

func (s *userStore) whereClause(filter *model.UserFilter) util.SqlWhereClause {
	wc := util.NewSqlWhereClause()
	if filter != nil {
		if filter.UserId != 0 {
			wc.WithCondition("id = $_", filter.UserId)
		}
	}
	return wc
}
