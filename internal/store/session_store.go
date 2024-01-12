package store

import (
	"context"
	"database/sql"
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
	Delete(ctx context.Context, tx *sql.Tx, token model.SessionToken)
	CleanUp(ctx context.Context, tx *sql.Tx, refTime time.Time)
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
		Expiration: obj.Expiration.Unix(),
	}
}

func (s *sessionStore) DecodeRow(row *SessionRow) *model.Session {
	if row == nil {
		return nil
	}
	return &model.Session{
		Token:      model.SessionToken(row.Token),
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
// delete

func (s *sessionStore) Delete(ctx context.Context, tx *sql.Tx, token model.SessionToken) {
	s.DeleteRows(ctx, tx, s.matchingToken(token))
}

// //////////////////////////////////////////////////
// delete

func (s *sessionStore) CleanUp(ctx context.Context, tx *sql.Tx, refTime time.Time) {
	s.DeleteRows(ctx, tx, s.olderThan(refTime))
}

// //////////////////////////////////////////////////
// where clause

func (s *sessionStore) matchingToken(token model.SessionToken) util.SqlWhereClause {
	return util.NewSqlCondition("token = %s", token)
}

func (s *sessionStore) olderThan(refTime time.Time) util.SqlWhereClause {
	return util.NewSqlCondition("expiration < %d", refTime.Unix())
}
