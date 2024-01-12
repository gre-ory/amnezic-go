package api

import (
	"net/http"
	"strings"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/service"
)

type Granter interface {
	Grant(req *http.Request) error
}

func Protect(granter Granter, nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if granter != nil {
			err := granter.Grant(req)
			if err != nil {
				encodeError(resp, http.StatusUnauthorized, err.Error())
				return
			}
		}
		nextHandler(resp, req)
	}
}

func NewPermissionGranter(permission model.Permission, sessionService service.SessionService) Granter {
	return &permissionGranter{
		permission: permission,
	}
}

type permissionGranter struct {
	permission     model.Permission
	sessionService service.SessionService
}

func (g *permissionGranter) Grant(req *http.Request) error {

	//
	// extract session token
	//

	token, err := extractSessionToken(req)
	if err != nil {
		return err
	}

	//
	// check session token
	//

	err = g.sessionService.IsGranted(req.Context(), token, g.permission)
	if err != nil {
		return err
	}

	return nil
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
