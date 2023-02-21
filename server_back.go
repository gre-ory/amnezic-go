package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gre-ory/amnezic-go/internal/api"
	"go.uber.org/zap"
)

func main_back() {

	ctx := context.Background()

	logger := zap.L()

	// services
	// gameService := service.NewGameService()

	// api
	log.Print("starting server...")
	logger.Info("[DEBUG] starting server...")

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}
	logger.Info("[DEBUG] port", zap.String("port", port))

	// Start HTTP server.

	mux := http.NewServeMux()
	// mux.Handle("/game/", api.NewGameHandler(gameService))
	mux.Handle("/root", api.NewRootHandler())
	mux.Handle("/", http.FileServer(http.Dir("www")))

	address := ":" + port

	server := http.Server{
		Addr:    address,
		Handler: mux,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Info("[DEBUG] error", zap.Error(err))
	}
}
