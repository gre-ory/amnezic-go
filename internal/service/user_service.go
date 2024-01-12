package service

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
	"github.com/gre-ory/amnezic-go/internal/util"
)

// //////////////////////////////////////////////////
// user service

type UserService interface {
	ListUsers(ctx context.Context) ([]*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	RetrieveUser(ctx context.Context, id model.UserId) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id model.UserId) error

	AddPermission(ctx context.Context, id model.UserId, permission model.Permission) (*model.User, error)
	RemovePermission(ctx context.Context, id model.UserId, permission model.Permission) (*model.User, error)

	CheckUser(ctx context.Context, login *model.LoginRequest) (*model.User, error)
}

func NewUserService(logger *zap.Logger, db *sql.DB, userStore store.UserStore) UserService {
	return &userService{
		logger:    logger,
		db:        db,
		userStore: userStore,
	}
}

type userService struct {
	logger    *zap.Logger
	db        *sql.DB
	userStore store.UserStore
}

func (s *userService) ListUsers(ctx context.Context) ([]*model.User, error) {

	var users []*model.User

	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {
		//
		// list users
		//

		s.logger.Info("[DEBUG] list users")
		users = s.userStore.List(ctx, tx, nil)
	})

	if err != nil {
		s.logger.Info("[ KO ] list users", zap.Error(err))
		return nil, err
	}
	s.logger.Info("[ OK ] list users")
	return users, nil
}

// //////////////////////////////////////////////////
// create

func (s *userService) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {

	var created *model.User
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// compute password hash
		//

		if user.Password == "" {
			panic(model.ErrInvalidPassword)
		}
		newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		user.Hash = string(newHash)

		//
		// create user
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] create user: %#v", user))
		created = s.userStore.Create(ctx, tx, user)

	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] create user: %#v", user), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] create user: %#v", created))
	return created, nil
}

// //////////////////////////////////////////////////
// retrieve

func (s *userService) RetrieveUser(ctx context.Context, id model.UserId) (*model.User, error) {

	var user *model.User
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {
		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve user %d", id))
		user = s.userStore.Retrieve(ctx, tx, id)
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] retrieve user %d", id), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] retrieve user %d", id))
	return user, nil
}

// //////////////////////////////////////////////////
// delete

func (s *userService) DeleteUser(ctx context.Context, id model.UserId) error {

	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// delete user
		//

		s.userStore.Delete(ctx, tx, &model.UserFilter{UserId: id})

	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] delete user: %#v", id), zap.Error(err))
		return err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] delete user: %#v", id))
	return nil
}

// //////////////////////////////////////////////////
// update

func (s *userService) UpdateUser(ctx context.Context, values *model.User) (*model.User, error) {

	id := values.Id
	var user *model.User
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// retrieve user
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve user %d", id))
		orig := s.userStore.Retrieve(ctx, tx, id)

		//
		// update
		//

		target := orig.Copy()
		if values.Name != "" {
			target.Name = values.Name
		}
		if values.Password != "" {
			if values.OldPassword == "" {
				panic(model.ErrInvalidPassword)
			}
			if err := bcrypt.CompareHashAndPassword([]byte(orig.Hash), []byte(values.OldPassword)); err != nil {
				s.logger.Info(fmt.Sprintf("[DEBUG] invalid password for user %d", values.Id), zap.Error(err))
				panic(model.ErrInvalidPassword)
			}
			newHash, err := bcrypt.GenerateFromPassword([]byte(values.Password), bcrypt.DefaultCost)
			if err != nil {
				panic(err)
			}
			target.Hash = string(newHash)
		}

		//
		// update user
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] update user: %#v", target))
		user = s.userStore.Update(ctx, tx, target)

	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] update user %d", id), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] update user %d - %s", id, user.Name), zap.Object("user", user))
	return user, nil
}

// //////////////////////////////////////////////////
// add permission

func (s *userService) AddPermission(ctx context.Context, id model.UserId, permission model.Permission) (*model.User, error) {

	var user *model.User
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// retrieve user
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve user %d", id))
		orig := s.userStore.Retrieve(ctx, tx, id)

		//
		// add permission
		//

		target := orig.Copy()
		target.AddPermission(permission)

		//
		// update user
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] update user: %#v", target))
		user = s.userStore.Update(ctx, tx, target)

	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] add permission %s to user %d", permission, user.Id), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] add permission %s to user %d", permission, user.Id), zap.Object("user", user))
	return user, nil
}

// //////////////////////////////////////////////////
// remove permission

func (s *userService) RemovePermission(ctx context.Context, id model.UserId, permission model.Permission) (*model.User, error) {

	var user *model.User
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// retrieve user
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve user %d", id))
		orig := s.userStore.Retrieve(ctx, tx, id)

		//
		// add permission
		//

		target := orig.Copy()
		target.RemovePermission(permission)

		//
		// update user
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] update user: %#v", target))
		user = s.userStore.Update(ctx, tx, target)

	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] remove permission %s to user %d", permission, user.Id), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] remove permission %s to user %d", permission, user.Id), zap.Object("user", user))
	return user, nil
}

// //////////////////////////////////////////////////
// check

func (s *userService) CheckUser(ctx context.Context, login *model.LoginRequest) (*model.User, error) {

	var user *model.User

	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// retrieve user from name
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] retrieve user %s", login.Name))
		user = s.userStore.RetrieveFromName(ctx, tx, login.Name)

		//
		// check password
		//

		if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(login.Password)); err != nil {
			s.logger.Info(fmt.Sprintf("[DEBUG] invalid password for user %d", user.Id), zap.Error(err))
			panic(model.ErrInvalidPassword)
		}
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] retrieve user %s", login.Name), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] retrieve user %s", login.Name), zap.Object("user", user))
	return user, nil
}
