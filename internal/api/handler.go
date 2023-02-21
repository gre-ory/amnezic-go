package api

import (
	"github.com/julienschmidt/httprouter"
)

// //////////////////////////////////////////////////
// handler

type Handler interface {
	RegisterRoutes(router *httprouter.Router)
}
