package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/service"
	"go.uber.org/zap"
)

type Granter interface {
	Grant(req *http.Request) (*http.Request, error)
}

func Protect(granter Granter, nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		var err error
		if granter != nil {
			req, err = granter.Grant(req)
			if err != nil {
				encodeError(resp, http.StatusUnauthorized, err.Error())
				return
			}
		}
		nextHandler(resp, req)
	}
}

func WithPermission(logger *zap.Logger, sessionService service.SessionService, permission model.Permission) func(http.HandlerFunc) http.HandlerFunc {
	granter := NewPermissionGranter(logger, sessionService, permission)
	return func(nextHanlder http.HandlerFunc) http.HandlerFunc {
		return Protect(granter, nextHanlder)
	}
}

func NewPermissionGranter(logger *zap.Logger, sessionService service.SessionService, permission model.Permission) Granter {
	return &permissionGranter{
		logger:         logger,
		sessionService: sessionService,
		permission:     permission,
	}
}

type permissionGranter struct {
	logger         *zap.Logger
	sessionService service.SessionService
	permission     model.Permission
}

func (g *permissionGranter) Grant(req *http.Request) (*http.Request, error) {

	//
	// extract session token
	//

	token, err := extractSessionToken(req)
	if err != nil {
		g.logger.Info("unable to extract session token", zap.Error(err))
		return req, err
	}

	//
	// check session token
	//

	session, err := g.sessionService.IsGranted(req.Context(), token, g.permission)
	if err != nil {
		g.logger.Info(fmt.Sprintf("user not granted for permission %s", g.permission), zap.Error(err))
		return req, err
	}
	req = req.WithContext(model.WithSession(req.Context(), session))

	g.logger.Info(fmt.Sprintf("user granted for permission %s", g.permission), zap.Error(err))
	return req, nil
}

// //////////////////////////////////////////////////
// extract session token

func extractSessionToken(req *http.Request) (model.SessionToken, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", model.ErrMissingAuthorizationHeader
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", model.ErrInvalidAuthorizationHeader
	}
	return model.SessionToken(authHeader[len("Bearer "):]), nil
}
