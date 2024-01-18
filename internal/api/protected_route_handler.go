package api

import (
	"fmt"
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
		permission:     permission,
		sessionService: sessionService,
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

	fmt.Printf("req: %+v\n", req)
	token, err := extractSessionToken(req)
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return err
	}
	fmt.Printf("token: %+v\n", token)

	//
	// check session token
	//

	fmt.Printf("permission: %+v\n", g.permission)
	err = g.sessionService.IsGranted(req.Context(), token, g.permission)
	if err != nil {
		fmt.Printf("err: %+v\n", err)
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
