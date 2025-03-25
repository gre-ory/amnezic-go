package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// health handler

func NewHealthhandler(logger *zap.Logger) Handler {
	return &healthHandler{
		logger: logger,
	}
}

type healthHandler struct {
	logger *zap.Logger
}

// //////////////////////////////////////////////////
// register

func (h *healthHandler) RegisterRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/ping", h.handlePing)
	router.HandlerFunc(http.MethodGet, "/health", h.handleHealth)
}

// //////////////////////////////////////////////////
// ping

func (h *healthHandler) handlePing(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusOK)
}

// //////////////////////////////////////////////////
// ping

func (h *healthHandler) handleHealth(resp http.ResponseWriter, req *http.Request) {
	// TODO implement specific health checks
	resp.WriteHeader(http.StatusOK)
}
