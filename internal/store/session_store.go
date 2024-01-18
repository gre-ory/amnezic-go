package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// session store

type SessionStore interface {
	Create(ctx context.Context, tx *sql.Tx, session *model.Session) *model.Session
	Retrieve(ctx context.Context, tx *sql.Tx, token model.SessionToken) *model.Session
	SearchByUserId(ctx context.Context, tx *sql.Tx, id model.UserId) *model.Session
	Delete(ctx context.Context, tx *sql.Tx, token model.SessionToken)
	CleanUp(ctx context.Context, tx *sql.Tx, refTime time.Time)
	List(ctx context.Context, tx *sql.Tx) []*model.Session
	Flush(ctx context.Context, tx *sql.Tx)
}

func NewSessionStore(logger *zap.Logger) SessionStore {
	return &sessionStore{
		SqlTable: util.NewSqlTable[SessionRow](logger, SessionTable, model.ErrSessionNotFound),
	}
}

type sessionStore struct {
	util.SqlTable[SessionRow]
	util.SqlEncoder[model.Session, SessionRow]
	util.SqlDecoder[SessionRow, model.Session]
}

// //////////////////////////////////////////////////
// table

const SessionTable = "user_session"

// //////////////////////////////////////////////////
// row

type SessionRow struct {
	Token      string `sql:"token"`
	UserId     int64  `sql:"user_id"`
	Expiration int64  `sql:"expiration"`
}

func (s *sessionStore) EncodeRow(obj *model.Session) *SessionRow {
	return &SessionRow{
		Token:      string(obj.Token),
		UserId:     int64(obj.UserId),
		Expiration: obj.Expiration.Unix(),
	}
}

func (s *sessionStore) DecodeRow(row *SessionRow) *model.Session {
	if row == nil {
		return nil
	}
	return &model.Session{
		Token:      model.SessionToken(row.Token),
		UserId:     model.UserId(row.UserId),
		Expiration: time.Unix(row.Expiration, 0),
	}
}

// //////////////////////////////////////////////////
// create

func (s *sessionStore) Create(ctx context.Context, tx *sql.Tx, obj *model.Session) *model.Session {
	return s.DecodeRow(s.InsertRow(ctx, tx, s.EncodeRow(obj)))
}

// //////////////////////////////////////////////////
// retrieve

func (s *sessionStore) Retrieve(ctx context.Context, tx *sql.Tx, token model.SessionToken) *model.Session {
	row, err := s.SelectRow(ctx, tx, s.matchingToken(token))
	if err != nil {
		panic(err)
	}
	return s.DecodeRow(row)
}

// //////////////////////////////////////////////////
// search by id

func (s *sessionStore) SearchByUserId(ctx context.Context, tx *sql.Tx, id model.UserId) *model.Session {
	row, err := s.SelectRow(ctx, tx, s.matchingUserId(id))
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			return nil
		} else {
			panic(err)
		}
	}
	return s.DecodeRow(row)
}

// //////////////////////////////////////////////////
// delete

func (s *sessionStore) Delete(ctx context.Context, tx *sql.Tx, token model.SessionToken) {
	s.DeleteRows(ctx, tx, s.matchingToken(token))
}

// //////////////////////////////////////////////////
// clean-up

func (s *sessionStore) CleanUp(ctx context.Context, tx *sql.Tx, refTime time.Time) {
	s.DeleteRows(ctx, tx, s.olderThan(refTime))
}

// //////////////////////////////////////////////////
// list

func (s *sessionStore) List(ctx context.Context, tx *sql.Tx) []*model.Session {
	return util.Convert(s.ListRows(ctx, tx, util.NoWhereClause), s.DecodeRow)
}

// //////////////////////////////////////////////////
// flush

func (s *sessionStore) Flush(ctx context.Context, tx *sql.Tx) {
	s.DeleteRows(ctx, tx, util.NoWhereClause)
}

// //////////////////////////////////////////////////
// where clause

func (s *sessionStore) matchingUserId(id model.UserId) util.SqlWhereClause {
	return util.NewSqlCondition("user_id = %s", id)
}

func (s *sessionStore) matchingToken(token model.SessionToken) util.SqlWhereClause {
	return util.NewSqlCondition("token = %s", token)
}

func (s *sessionStore) olderThan(refTime time.Time) util.SqlWhereClause {
	return util.NewSqlCondition("expiration < %s", refTime.Unix())
}
