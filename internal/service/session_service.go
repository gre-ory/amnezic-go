package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt"
	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
	"github.com/gre-ory/amnezic-go/internal/util"
)

// //////////////////////////////////////////////////
// session service

type SessionService interface {
	Login(ctx context.Context, login *model.LoginRequest) (*model.Session, error)
	IsGranted(ctx context.Context, token model.SessionToken, permission model.Permission) (*model.Session, error)
	Logout(ctx context.Context, token model.SessionToken) error

	ListSessions(ctx context.Context) ([]*model.Session, error)
	FlushSessions(ctx context.Context) error
}

func NewSessionService(logger *zap.Logger, secretKey string, db *sql.DB, sessionStore store.SessionStore, userStore store.UserStore) SessionService {
	return &sessionService{
		logger:       logger,
		secretKey:    secretKey,
		db:           db,
		sessionStore: sessionStore,
		userStore:    userStore,
	}
}

type sessionService struct {
	logger       *zap.Logger
	secretKey    string
	db           *sql.DB
	sessionStore store.SessionStore
	userStore    store.UserStore
}

// //////////////////////////////////////////////////
// list

func (s *sessionService) ListSessions(ctx context.Context) ([]*model.Session, error) {

	var sessions []*model.Session

	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {
		//
		// list users
		//

		s.logger.Info("[DEBUG] list sessions")
		sessions = s.sessionStore.List(ctx, tx)
	})

	if err != nil {
		s.logger.Info("[ KO ] list sessions", zap.Error(err))
		return nil, err
	}
	s.logger.Info("[ OK ] list sessions")
	return sessions, nil
}

// //////////////////////////////////////////////////
// list

func (s *sessionService) FlushSessions(ctx context.Context) error {

	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {
		//
		// list users
		//

		s.logger.Info("[DEBUG] flush all sessions")
		s.sessionStore.Flush(ctx, tx)
	})

	if err != nil {
		s.logger.Info("[ KO ] flush all sessions", zap.Error(err))
		return err
	}
	s.logger.Info("[ OK ] flush all sessions")
	return nil
}

// //////////////////////////////////////////////////
// login

func (s *sessionService) Login(ctx context.Context, login *model.LoginRequest) (*model.Session, error) {

	now := time.Now()

	var session *model.Session
	var user *model.User
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// clean-up expired sessions
		//

		s.logger.Info("[DEBUG] clean-up expired sessions")
		s.sessionStore.CleanUp(ctx, tx, now)

		//
		// retrieve user from name
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve user %s", login.Name))
		user = s.userStore.SearchByName(ctx, tx, login.Name)
		if user == nil {
			s.logger.Info(fmt.Sprintf("[DEBUG] user %s not found", login.Name))
			panic(model.ErrUserNotFound)
		}

		//
		// check for any previous session
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] check for any existing session for user %d", user.Id))
		session = s.sessionStore.SearchByUserId(ctx, tx, user.Id)
		if session != nil {
			s.logger.Info(fmt.Sprintf("[DEBUG] existing session for user %d", user.Id))
			return
		}

		//
		// check password
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] check password for user %d: %s >>> %s", user.Id, login.Password, user.Hash))
		err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(login.Password))
		if err != nil {
			s.logger.Info(fmt.Sprintf("[DEBUG] invalid password for user %d", user.Id), zap.Error(err))
			panic(model.ErrInvalidPassword)
		}

		//
		// create session
		//

		newSession, err := s.newSession(user.Id, 24*time.Hour)
		if err != nil {
			s.logger.Info(fmt.Sprintf("[DEBUG] unable to create new session for user %d", user.Id), zap.Error(err))
			panic(err)
		}
		s.logger.Info(fmt.Sprintf("[DEBUG] create session: %#v", newSession))
		session = s.sessionStore.Create(ctx, tx, newSession)
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] login for user %s", login.Name), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] login for user %s", login.Name), zap.Object("session", session), zap.Object("user", user))
	session.User = user
	return session, nil
}

func (s *sessionService) newSession(userId model.UserId, ttl time.Duration) (*model.Session, error) {

	expiration := time.Now().Add(ttl)

	jwtToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": userId,
			"exp":      expiration.Unix(),
		},
	)

	token, err := jwtToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, err
	}

	return &model.Session{
		Token:      model.SessionToken(token),
		UserId:     userId,
		Expiration: expiration,
	}, nil
}

// //////////////////////////////////////////////////
// retrieve

func (s *sessionService) IsGranted(ctx context.Context, token model.SessionToken, permission model.Permission) (*model.Session, error) {

	now := time.Now()

	var session *model.Session
	var user *model.User
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// clean-up expired sessions
		//

		s.logger.Info("[DEBUG] clean-up expired sessions")
		s.sessionStore.CleanUp(ctx, tx, now)

		//
		// retrieve session
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve session %s", token))
		session = s.sessionStore.Retrieve(ctx, tx, token)

		//
		// validate session
		//

		jwtToken, err := jwt.Parse(
			session.Token.String(),
			func(token *jwt.Token) (interface{}, error) {
				return []byte(s.secretKey), nil
			},
		)
		if err != nil {
			s.logger.Info(fmt.Sprintf("[DEBUG] invalid session token: %s", token), zap.Error(err))
			panic(model.ErrInvalidSessionToken)
		}
		if !jwtToken.Valid {
			s.logger.Info(fmt.Sprintf("[DEBUG] invalid session token: %s", token))
			panic(model.ErrInvalidSessionToken)
		}
		// TODO need more validation

		//
		// retrieve user
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve user %s", session.UserId))
		user = s.userStore.Retrieve(ctx, tx, session.UserId)

		//
		// check permission
		//

		if !user.HasPermission(permission) {
			panic(model.ErrUserNotGranted)
		}
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] is granted for permission %s", permission), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] is granted for permission %s", permission))
	session.User = user
	return session, nil
}

// //////////////////////////////////////////////////
// retrieve

func (s *sessionService) Logout(ctx context.Context, token model.SessionToken) error {

	now := time.Now()

	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// clean-up expired sessions
		//

		s.logger.Info("[DEBUG] clean-up expired sessions")
		s.sessionStore.CleanUp(ctx, tx, now)

		//
		// delete session
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] delete session %s", token))
		s.sessionStore.Delete(ctx, tx, token)

	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] delete session %s", token), zap.Error(err))
		return err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] delete session %s", token))
	return nil
}
