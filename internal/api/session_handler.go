package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/service"
	"github.com/gre-ory/amnezic-go/internal/util"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// session handler

func NewSessionhandler(logger *zap.Logger, sessionService service.SessionService) Handler {
	return &sessionHandler{
		logger:         logger,
		sessionService: sessionService,
	}
}

type sessionHandler struct {
	logger         *zap.Logger
	sessionService service.SessionService
}

// //////////////////////////////////////////////////
// register

func (h *sessionHandler) RegisterRoutes(router *httprouter.Router) {

	router.HandlerFunc(http.MethodPut, "/api/login", h.handleLogin)
	router.HandlerFunc(http.MethodDelete, "/api/logout", h.handleLogout)

	withSessionPermission := WithPermission(h.logger, h.sessionService, model.Permission_Session)

	router.HandlerFunc(http.MethodGet, "/api/session", withSessionPermission(h.handleListSession))
	router.HandlerFunc(http.MethodDelete, "/api/session", withSessionPermission(h.handleFlushSession))
}

// //////////////////////////////////////////////////
// list

func (h *sessionHandler) handleListSession(resp http.ResponseWriter, req *http.Request) {
	defer onPanic(resp)()

	ctx := req.Context()

	var sessions []*model.Session
	var err error

	switch {
	default:

		h.logger.Info("[api] list users")

		//
		// execute
		//

		sessions, err = h.sessionService.ListSessions(ctx)
		if err != nil {
			break
		}

		//
		// encode response
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonSessionsResponse(sessions))
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
// flush

func (h *sessionHandler) handleFlushSession(resp http.ResponseWriter, req *http.Request) {
	defer onPanic(resp)()

	ctx := req.Context()

	var err error

	switch {
	default:

		h.logger.Info("[api] flush all sessions")

		//
		// execute
		//

		err = h.sessionService.FlushSessions(ctx)
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
// login

func (h *sessionHandler) handleLogin(resp http.ResponseWriter, req *http.Request) {
	defer onPanic(resp)()

	ctx := req.Context()

	var loginRequest *model.LoginRequest
	var session *model.Session
	var err error

	switch {
	default:

		h.logger.Info("[api] login")

		//
		// extract login request
		//

		loginRequest, err = extractLoginRequestFromBody(req, h.logger)
		if err != nil {
			break
		}

		//
		// execute
		//

		session, err = h.sessionService.Login(ctx, loginRequest)
		if err != nil {
			break
		}

		//
		// encode response
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonSessionResponse(session))
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
// login

func (h *sessionHandler) handleLogout(resp http.ResponseWriter, req *http.Request) {
	defer onPanic(resp)()

	ctx := req.Context()

	var token model.SessionToken
	var err error

	switch {
	default:

		h.logger.Info("[api] logout")

		//
		// extract session token
		//

		token, err = extractSessionToken(req)
		if err != nil {
			break
		}

		//
		// execute
		//

		err = h.sessionService.Logout(ctx, token)
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
// decode

func extractLoginRequestFromBody(req *http.Request, logger *zap.Logger) (*model.LoginRequest, error) {
	var jsonBody JsonLoginBody
	jsonErr := json.NewDecoder(req.Body).Decode(&jsonBody)
	switch {
	case jsonErr == io.EOF:
		logger.Info("failed to decode login body: EOF")
		return nil, model.ErrInvalidBody
	case jsonErr != nil:
		logger.Info("failed to decode login body", zap.Error(jsonErr))
		return nil, model.ErrInvalidBody
	case jsonBody.Login == nil:
		logger.Info("missing login info", zap.Error(jsonErr))
		return nil, model.ErrInvalidBody
	}

	return toLoginRequest(jsonBody.Login), nil
}

func toLoginRequest(jsonLoginRequest *JsonLoginRequest) *model.LoginRequest {
	return &model.LoginRequest{
		Name:     jsonLoginRequest.Name,
		Password: jsonLoginRequest.Password,
	}
}

type JsonLoginBody struct {
	Login *JsonLoginRequest `json:"login,omitempty"`
}

type JsonLoginRequest struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

// //////////////////////////////////////////////////
// encode

func toJsonSessionsResponse(sessions []*model.Session) *JsonSessionsResponse {
	return &JsonSessionsResponse{
		Success:  true,
		Sessions: util.Convert(sessions, toJsonSession),
	}
}

func toJsonSessionResponse(session *model.Session) *JsonSessionResponse {
	return &JsonSessionResponse{
		Success: true,
		Session: toJsonSession(session),
	}
}

func toJsonSession(session *model.Session) *JsonSession {
	jsonSession := &JsonSession{
		Token:        session.Token.String(),
		ExpirationTs: session.Expiration.Unix(),
		Expiration:   session.Expiration.Format("2006-01-02T15:04:05Z"),
	}
	if session.User != nil {
		jsonSession.User = toJsonUser(session.User)
	} else {
		jsonSession.UserId = session.UserId.ToInt64()
	}
	return jsonSession
}

type JsonSessionsResponse struct {
	Success  bool           `json:"success,omitempty"`
	Sessions []*JsonSession `json:"sessions,omitempty"`
}

type JsonSessionResponse struct {
	Success bool         `json:"success,omitempty"`
	Session *JsonSession `json:"session,omitempty"`
}

type JsonSession struct {
	Token        string    `json:"token,omitempty"`
	UserId       int64     `json:"userId,omitempty"`
	User         *JsonUser `json:"user,omitempty"`
	ExpirationTs int64     `json:"expirationTs,omitempty"`
	Expiration   string    `json:"expiration,omitempty"`
}

// // //////////////////////////////////////////////////
// // jwt

// const (
// 	JwtExpiration = 24 * time.Hour
// )

// var secretKey = []byte("secret-key")

// func (h *sessionHandler) createJwtToken(username string) (string, error) {

// 	token := jwt.NewWithClaims(
// 		jwt.SigningMethodHS256,
// 		jwt.MapClaims{
// 			"username": username,
// 			"exp":      time.Now().Add(JwtExpiration).Unix(),
// 		},
// 	)

// 	tokenString, err := token.SignedString(secretKey)
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }

// func (h *sessionHandler) verifyToken(tokenString string) error {

// 	token, err := jwt.Parse(
// 		tokenString,
// 		func(token *jwt.Token) (interface{}, error) {
// 			return secretKey, nil
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	if !token.Valid {
// 		return fmt.Errorf("invalid token")
// 	}

// 	return nil
// }

// func LoginHandler(resp http.ResponseWriter, req *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	var u User
// 	json.NewDecoder(r.Body).Decode(&u)
// 	fmt.Printf("The user request value %v", u)

// 	if u.Username == "Chek" && u.Password == "123456" {
// 		tokenString, err := CreateToken(u.Username)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			fmt.Errorf("No username found")
// 		}
// 		w.WriteHeader(http.StatusOK)
// 		fmt.Fprint(w, tokenString)
// 		return
// 	} else {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		fmt.Fprint(w, "Invalid credentials")
// 	}
// }
