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
	IsGranted(ctx context.Context, token model.SessionToken, permission model.Permission) error
	Logout(ctx context.Context, token model.SessionToken) error
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
// login

func (s *sessionService) Login(ctx context.Context, login *model.LoginRequest) (*model.Session, error) {

	now := time.Now()

	var user *model.User
	var created *model.Session
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
		user = s.userStore.RetrieveFromName(ctx, tx, login.Name)

		//
		// check password
		//

		err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(login.Password))
		if err != nil {
			s.logger.Info(fmt.Sprintf("[DEBUG] invalid password for user %d", user.Id), zap.Error(err))
			panic(model.ErrInvalidPassword)
		}

		//
		// create session
		//

		session, err := s.newSession(user.Id, 24*time.Hour)
		if err != nil {
			panic(err)
		}
		s.logger.Info(fmt.Sprintf("[DEBUG] create session: %#v", session))
		created = s.sessionStore.Create(ctx, tx, session)
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] create session for user %s", login.Name), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] create session for user %s", login.Name), zap.Object("session", created))
	return created, nil
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

	token, err := jwtToken.SignedString(s.secretKey)
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

func (s *sessionService) IsGranted(ctx context.Context, token model.SessionToken, permission model.Permission) error {

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
				return s.secretKey, nil
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
		return err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] is granted for permission %s", permission))
	return nil
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
