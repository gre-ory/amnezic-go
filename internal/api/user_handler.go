package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/service"
	"github.com/gre-ory/amnezic-go/internal/util"
)

// //////////////////////////////////////////////////
// user handler

func NewUserHandler(logger *zap.Logger, userService service.UserService, sessionService service.SessionService) Handler {
	return &userHandler{
		logger:         logger,
		userService:    userService,
		sessionService: sessionService,
	}
}

type userHandler struct {
	logger         *zap.Logger
	userService    service.UserService
	sessionService service.SessionService
}

// //////////////////////////////////////////////////
// register

func (h *userHandler) RegisterRoutes(router *httprouter.Router) {

	router.HandlerFunc(http.MethodPut, "/api/user/set-up", h.handleUserSetUp)

	withUserPermission := WithPermission(h.logger, h.sessionService, model.Permission_User)

	router.HandlerFunc(http.MethodGet, "/api/user", withUserPermission(h.handleListUser))
	router.HandlerFunc(http.MethodGet, "/api/user/:user_id", withUserPermission(h.handleRetrieveUser))
	router.HandlerFunc(http.MethodPut, "/api/user/new", withUserPermission(h.handleCreateUser))
	router.HandlerFunc(http.MethodPost, "/api/user/:user_id", withUserPermission(h.handleUpdateUser))
	router.HandlerFunc(http.MethodDelete, "/api/user/:user_id", withUserPermission(h.handleDeleteUser))
	router.HandlerFunc(http.MethodPut, "/api/user-permission/:user_id/:permission", withUserPermission(h.handleAddPermission))
	router.HandlerFunc(http.MethodDelete, "/api/user-permission/:user_id/:permission", withUserPermission(h.handleRemovePermission))
}

// //////////////////////////////////////////////////
// set-up

func (h *userHandler) handleUserSetUp(resp http.ResponseWriter, req *http.Request) {
	defer onPanic(resp)()

	ctx := req.Context()

	var err error

	switch {
	default:

		h.logger.Info("[api] user set-up")

		//
		// execute
		//

		err = h.userService.CreateDefaultAdminUser(ctx)
		if err != nil {
			break
		}

		//
		// encode response
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonSuccess())
		if err != nil {
			break
		}
		return

	}

	//
	// encode error
	//

	// TODO status code
	encodeError(resp, http.StatusBadRequest, err.Error())
}

// //////////////////////////////////////////////////
// list

func (h *userHandler) handleListUser(resp http.ResponseWriter, req *http.Request) {
	defer onPanic(resp)()

	ctx := req.Context()

	var users []*model.User
	var err error

	switch {
	default:

		h.logger.Info("[api] list users")

		//
		// execute
		//

		users, err = h.userService.ListUsers(ctx)
		if err != nil {
			break
		}

		//
		// encode response
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonUsersResponse(users))
		if err != nil {
			break
		}
		return

	}

	//
	// encode error
	//

	// TODO status code
	encodeError(resp, http.StatusBadRequest, err.Error())
}

// //////////////////////////////////////////////////
// create

func (h *userHandler) handleCreateUser(resp http.ResponseWriter, req *http.Request) {
	defer onPanic(resp)()

	ctx := req.Context()

	var user *model.User
	var err error

	switch {
	default:

		//
		// decode request
		//

		user, err = extractUserFromBody(req, h.logger)
		if err != nil {
			break
		}
		h.logger.Info(fmt.Sprintf("[api] create user %s", user.Name))

		//
		// execute
		//

		user, err = h.userService.CreateUser(ctx, user)
		if err != nil {
			break
		}
		if user == nil {
			err = model.ErrUserNotFound
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonUserResponse(user))
		if err != nil {
			break
		}
		return
	}

	//
	// encode error
	//

	// TODO status code
	encodeError(resp, http.StatusBadRequest, err.Error())
}

// //////////////////////////////////////////////////
// retrieve

func (h *userHandler) handleRetrieveUser(resp http.ResponseWriter, req *http.Request) {
	defer onPanic(resp)()

	ctx := req.Context()

	var userId model.UserId
	var user *model.User
	var err error

	switch {
	default:

		//
		// decode request
		//

		userId = model.UserId(toInt64(extractPathParameter(req, "user_id")))
		if userId == 0 {
			err = model.ErrInvalidUserId
			break
		}
		h.logger.Info(fmt.Sprintf("[api] retrieve user %d", userId))

		//
		// execute
		//

		user, err = h.userService.RetrieveUser(ctx, userId)
		if err != nil {
			break
		}
		if user == nil {
			err = model.ErrUserNotFound
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonUserResponse(user))
		if err != nil {
			break
		}
		return
	}

	//
	// encode error
	//

	// TODO status code
	encodeError(resp, http.StatusBadRequest, err.Error())
}

// //////////////////////////////////////////////////
// update

func (h *userHandler) handleUpdateUser(resp http.ResponseWriter, req *http.Request) {
	defer onPanic(resp)()

	ctx := req.Context()

	var userId model.UserId
	var user *model.User
	var err error

	switch {
	default:

		//
		// decode request
		//

		userId = model.UserId(toInt64(extractPathParameter(req, "user_id")))
		if userId == 0 {
			err = model.ErrInvalidUserId
			break
		}
		user, err = extractUserFromBody(req, h.logger)
		if err != nil {
			break
		}
		if user.Id == 0 {
			user.Id = userId
		} else if user.Id != userId {
			err = model.ErrInvalidUserId
			break
		}
		h.logger.Info(fmt.Sprintf("[api] update user %d", userId))

		//
		// execute
		//

		user, err = h.userService.UpdateUser(ctx, user)
		if err != nil {
			break
		}
		if user == nil {
			err = model.ErrUserNotFound
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonUserResponse(user))
		if err != nil {
			break
		}
		return
	}

	//
	// encode error
	//

	// TODO status code
	encodeError(resp, http.StatusBadRequest, err.Error())
}

// //////////////////////////////////////////////////
// delete

func (h *userHandler) handleDeleteUser(resp http.ResponseWriter, req *http.Request) {
	defer onPanic(resp)()

	ctx := req.Context()

	var userId model.UserId
	var err error

	switch {
	default:

		//
		// decode request
		//

		userId = model.UserId(toInt64(extractPathParameter(req, "user_id")))
		if userId == 0 {
			err = model.ErrInvalidUserId
			break
		}
		h.logger.Info(fmt.Sprintf("[api] delete user %d", userId))

		//
		// execute
		//

		err = h.userService.DeleteUser(ctx, userId)
		if err != nil {
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonSuccess())
		if err != nil {
			break
		}
		return
	}

	//
	// encode error
	//

	// TODO status code
	encodeError(resp, http.StatusBadRequest, err.Error())
}

// //////////////////////////////////////////////////
// add permission

func (h *userHandler) handleAddPermission(resp http.ResponseWriter, req *http.Request) {
	defer onPanic(resp)()

	ctx := req.Context()

	var userId model.UserId
	var permission model.Permission
	var user *model.User
	var err error

	switch {
	default:

		//
		// decode request
		//

		userId = model.UserId(toInt64(extractPathParameter(req, "user_id")))
		if userId == 0 {
			err = model.ErrInvalidUserId
			break
		}
		permission = model.ToPermission(extractPathParameter(req, "permission"))
		if permission == "" {
			err = model.ErrInvalidPermission
			break
		}
		h.logger.Info(fmt.Sprintf("[api] add permission %s to user %d", permission, userId))

		//
		// execute
		//

		user, err = h.userService.AddPermission(ctx, userId, permission)
		if err != nil {
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonUserResponse(user))
		if err != nil {
			break
		}
		return
	}

	//
	// encode error
	//

	// TODO status code
	encodeError(resp, http.StatusBadRequest, err.Error())
}

// //////////////////////////////////////////////////
// remove permission

func (h *userHandler) handleRemovePermission(resp http.ResponseWriter, req *http.Request) {
	defer onPanic(resp)()

	ctx := req.Context()

	var userId model.UserId
	var permission model.Permission
	var user *model.User
	var err error

	switch {
	default:

		//
		// decode request
		//

		userId = model.UserId(toInt64(extractPathParameter(req, "user_id")))
		if userId == 0 {
			err = model.ErrInvalidUserId
			break
		}
		permission = model.ToPermission(extractPathParameter(req, "permission"))
		if permission == "" {
			err = model.ErrInvalidPermission
			break
		}
		h.logger.Info(fmt.Sprintf("[api] remove permission %s to user %d", permission, userId))

		//
		// execute
		//

		user, err = h.userService.RemovePermission(ctx, userId, permission)
		if err != nil {
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonUserResponse(user))
		if err != nil {
			break
		}
		return
	}

	//
	// encode error
	//

	// TODO status code
	encodeError(resp, http.StatusBadRequest, err.Error())
}

// //////////////////////////////////////////////////
// decode

func extractUserFromBody(req *http.Request, logger *zap.Logger) (*model.User, error) {
	var jsonBody JsonUserBody
	jsonErr := json.NewDecoder(req.Body).Decode(&jsonBody)
	switch {
	case jsonErr == io.EOF:
		logger.Info("failed to decode user body: EOF")
		return nil, model.ErrInvalidBody
	case jsonErr != nil:
		logger.Info("failed to decode user body", zap.Error(jsonErr))
		return nil, model.ErrInvalidBody
	}

	return toUser(jsonBody.User), nil
}

func toUser(jsonUser *JsonUser) *model.User {
	return &model.User{
		Id:          model.UserId(jsonUser.Id),
		Name:        jsonUser.Name,
		OldPassword: jsonUser.OldPassword,
		Password:    jsonUser.Password,
		Permissions: util.Convert(jsonUser.Permissions, model.ToPermission),
	}
}

type JsonUserBody struct {
	User *JsonUser `json:"user,omitempty"`
}

// //////////////////////////////////////////////////
// encode

func toJsonUsersResponse(users []*model.User) *JsonUsersResponse {
	return &JsonUsersResponse{
		Success: true,
		Users:   util.Convert(users, toJsonUser),
	}
}

func toJsonUserResponse(user *model.User) *JsonUserResponse {
	return &JsonUserResponse{
		Success: true,
		User:    toJsonUser(user),
	}
}

func toJsonUser(user *model.User) *JsonUser {
	return &JsonUser{
		Id:          user.Id.ToInt64(),
		Name:        user.Name,
		Permissions: util.Convert(user.Permissions, model.Permission.String),
	}
}

type JsonUsersResponse struct {
	Success bool        `json:"success,omitempty"`
	Users   []*JsonUser `json:"users,omitempty"`
}

type JsonUserResponse struct {
	Success bool      `json:"success,omitempty"`
	User    *JsonUser `json:"user,omitempty"`
}

type JsonUser struct {
	Id          int64    `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	OldPassword string   `json:"oldPassword,omitempty"`
	Password    string   `json:"password,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}
