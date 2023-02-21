package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"

	"github.com/gre-ory/amnezic-go/internal/api"
	"github.com/gre-ory/amnezic-go/internal/service"
	"github.com/gre-ory/amnezic-go/internal/store"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// main

func main() {

	var gameStore store.GameStore
	var musicStore store.MusicStore
	var gameService service.GameService
	var gameHandler api.Handler
	var router *httprouter.Router
	var staticFs, reactFs fs.FS
	var server http.Server
	var err error

	ctx := context.Background()

	//
	// config
	//

	env := os.Getenv("ENVIRONMENT")
	app := os.Getenv("APPLICATION_NAME")
	version := os.Getenv("APPLICATION_VERSION")
	address := os.Getenv("ADDRESS")

	fmt.Printf("env: %s, app: %s %s, address: %s \n", env, app, version, address)

	//
	// logger
	//

	logger, err := NewLogger()
	if err != nil {
		goto exit_on_error
	}
	logger = logger.With(zap.String("env", env), zap.String("app", app), zap.String("version", version))
	logger.Info(fmt.Sprintf("creating %s server...", env))

	//
	// store
	//

	gameStore = store.NewGameMemoryStore()
	musicStore = store.NewLegacyMusicStore(store.RootPath_FreeDotFr)

	//
	// client
	//

	//
	// service
	//

	gameService = service.NewGameService(logger, gameStore, musicStore)

	//
	// api
	//

	gameHandler = api.NewGamehandler(logger, gameService)

	//
	// router
	//

	router = httprouter.New()
	gameHandler.RegisterRoutes(router)

	reactFs, err = fs.Sub(react, reactSubPath)
	if err != nil {
		goto exit_on_error
	}
	router.ServeFiles(reactPath, http.FS(reactFs))

	staticFs, err = fs.Sub(static, staticSubPath)
	if err != nil {
		goto exit_on_error
	}
	router.ServeFiles(staticPath, http.FS(staticFs))

	//
	// server
	//

	server = http.Server{
		Addr:    address,
		Handler: WithRequestLogging(logger)(router),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	logger.Info(fmt.Sprintf("starting %s server on %s", env, server.Addr))
	err = server.ListenAndServe()
	if err != nil {
		goto exit_on_error
	}

	logger.Info("server completed")
	os.Exit(0)

exit_on_error:
	if logger != nil {
		logger.Error("server failed!", zap.Error(err))
	} else if err != nil {
		fmt.Printf("server failed: %s! \n", err.Error())
	} else {
		fmt.Printf("server failed without error! \n")
	}
	os.Exit(1)
}

// //////////////////////////////////////////////////
// file serve

var reactPath = "/react/*filepath"
var reactSubPath = "www/app/amnezic/react/1.0.0"

//go:embed www/app/amnezic/react/1.0.0
var react embed.FS

var staticPath = "/static/*filepath"
var staticSubPath = "www/static"

//go:embed www/static
var static embed.FS

// //////////////////////////////////////////////////
// logger

func NewLogger() (*zap.Logger, error) {

	var zapCfg zap.Config
	switch os.Getenv("LOG_CONFIG") {
	case "prd":
		zapCfg = zap.NewProductionConfig()
	case "dev":
		zapCfg = zap.NewDevelopmentConfig()
	default:
		fmt.Printf("invalid LOG_CONFIG >>> FALLBACK to 'prd'!\n")
		zapCfg = zap.NewProductionConfig()
	}

	switch os.Getenv("LOG_LEVEL") {
	case "err":
		zapCfg.Level.SetLevel(zap.ErrorLevel)
	case "warn":
		zapCfg.Level.SetLevel(zap.WarnLevel)
	case "info":
		zapCfg.Level.SetLevel(zap.InfoLevel)
	case "debug":
		zapCfg.Level.SetLevel(zap.DebugLevel)
	default:
		fmt.Printf("invalid LOG_LEVEL >>> FALLBACK to 'info'!\n")
		zapCfg.Level.SetLevel(zap.InfoLevel)
	}

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}
	logger.Info("logger has been initialized")
	return logger, nil
}

// //////////////////////////////////////////////////
// request logging

func WithRequestLogging(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				logger.Info(fmt.Sprintf("[DEBUG] %s %s - %s", r.Method, r.URL.Path, r.UserAgent()))
			}()
			next.ServeHTTP(w, r)
		})
	}
}
