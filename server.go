package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gre-ory/amnezic-go/internal/api"
	"github.com/gre-ory/amnezic-go/internal/client"
	"github.com/gre-ory/amnezic-go/internal/service"
	"github.com/gre-ory/amnezic-go/internal/store"
	"github.com/gre-ory/amnezic-go/internal/util"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

// //////////////////////////////////////////////////
// main

func main() {

	ctx := context.Background()

	//
	// config
	//

	env := os.Getenv("ENVIRONMENT")
	app := os.Getenv("APPLICATION_NAME")
	version := os.Getenv("APPLICATION_VERSION")

	//
	// logger
	//

	logger := NewLogger()
	logger = logger.With(zap.String("env", env), zap.String("app", app), zap.String("version", version))
	logger.Info("starting app...")

	//
	// servers
	//

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		NewBackendServer(logger).Run(ctx)
		wg.Done()
	}()

	go func() {
		NewFrontendServer(logger).Run(ctx)
		wg.Done()
	}()

	wg.Wait()

	logger.Info("app completed")
	os.Exit(0)
}

// //////////////////////////////////////////////////
// frontend server

var frontendRootPath = "www"

//go:embed www
var frontendFS embed.FS

type FrontendServer struct {
	logger *zap.Logger
}

func NewFrontendServer(logger *zap.Logger) *FrontendServer {
	return &FrontendServer{
		logger: logger.With(zap.String("server", "frontend")),
	}
}

func (s *FrontendServer) Run(ctx context.Context) {

	//
	// config
	//

	address := os.Getenv("FRONTEND_ADDRESS")

	//
	// file server
	//

	rootFS, err := fs.Sub(frontendFS, frontendRootPath)
	if err != nil {
		s.logger.Fatal("failed to serve frontend files", zap.Error(err))
	}
	fileServer := http.FileServer(http.FS(rootFS))

	//
	// server
	//

	server := http.Server{
		Addr: address,
		Handler: WithRequestLogging(s.logger)(
			fileServer,
		),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	s.logger.Info(fmt.Sprintf("starting frontend server on %s", server.Addr))
	err = server.ListenAndServe()
	if err != nil {
		s.logger.Fatal("failed to start backend server", zap.Error(err))
	}
}

// //////////////////////////////////////////////////
// backend server

func NewBackendServer(logger *zap.Logger) *BackendServer {
	return &BackendServer{
		logger: logger.With(zap.String("server", "backend")),
	}
}

type BackendServer struct {
	logger *zap.Logger
}

func (s *BackendServer) Run(ctx context.Context) {

	//
	// config
	//

	address := os.Getenv("BACKEND_ADDRESS")

	//
	// store
	//

	db, _ := sql.Open("sqlite3", "./db/amnezic.db")
	defer db.Close()

	gameStore := store.NewGameMemoryStore()
	gameQuestionStore := store.NewGameQuestionLegacyStore(s.logger, store.RootPath_FreeDotFr)

	// musicStore := store.NewMusicMemoryStore()
	// albumStore := store.NewMusicAlbumMemoryStore()
	// artistStore := store.NewMusicArtistMemoryStore()

	musicStore := store.NewMusicStore(s.logger)
	albumStore := store.NewMusicAlbumStore(s.logger)
	artistStore := store.NewMusicArtistStore(s.logger)

	themeStore := store.NewThemeMemoryStore()
	themeQuestionStore := store.NewThemeQuestionMemoryStore()

	//
	// client
	//

	deezerClient := client.NewDeezerClient(s.logger)

	//
	// service
	//

	gameService := service.NewGameService(s.logger, db, gameStore, gameQuestionStore)
	musicService := service.NewMusicService(s.logger, deezerClient, db, musicStore, albumStore, artistStore)
	themeService := service.NewThemeService(s.logger, db, themeStore, themeQuestionStore, musicStore)

	//
	// api
	//

	gameHandler := api.NewGamehandler(s.logger, gameService)
	musicHandler := api.NewMusichandler(s.logger, musicService)
	themeHandler := api.NewThemehandler(s.logger, themeService)

	//
	// router
	//

	router := httprouter.New()
	gameHandler.RegisterRoutes(router)
	musicHandler.RegisterRoutes(router)
	themeHandler.RegisterRoutes(router)

	//
	// server
	//

	server := http.Server{
		Addr: address,
		Handler: AllowCORS(s.logger)(
			WithRequestLogging(s.logger)(
				router,
			),
		),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	s.logger.Info(fmt.Sprintf("starting backend server on %s", server.Addr))
	err := server.ListenAndServe()
	if err != nil {
		s.logger.Fatal("backend server failed", zap.Error(err))
	}
}

// //////////////////////////////////////////////////
// basic dispatcher

type dispatcher struct {
	extraHandlers  map[string]http.Handler
	defaultHandler http.Handler
}

func NewDispatcher(defaultHandler http.Handler) *dispatcher {
	return &dispatcher{
		extraHandlers:  make(map[string]http.Handler, 0),
		defaultHandler: defaultHandler,
	}
}

func (d *dispatcher) Register(path string, handler http.Handler) {
	d.extraHandlers[path] = handler
}

func (d *dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for path, handler := range d.extraHandlers {
		if strings.HasPrefix(r.URL.Path, path) {
			handler.ServeHTTP(w, r)
			return
		}
	}
	d.defaultHandler.ServeHTTP(w, r)
}

// //////////////////////////////////////////////////
// logger

func NewLogger() *zap.Logger {

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
		fmt.Printf("failed to initialize logger: %s \n", err.Error())
		os.Exit(1)
	}
	logger.Info("logger has been initialized")
	return logger
}

// //////////////////////////////////////////////////
// request logging

func WithRequestLogging(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(fmt.Sprintf("[DEBUG] %s %s - %s", r.Method, r.URL.Path, r.UserAgent()))
			next.ServeHTTP(w, r)
		})
	}
}

// //////////////////////////////////////////////////
// cors

var whitelistOrigins []string = []string{
	"http://localhost:3000",
	"http://localhost:9090",
}

func AllowCORS(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowed := util.Contains(whitelistOrigins, origin)
			if allowed {
				logger.Info(fmt.Sprintf("[COR] OK - Origin: %s", origin))
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "*")
			} else {
				logger.Warn(fmt.Sprintf("[COR] BLOCKED - Origin: %s", origin))
			}
			next.ServeHTTP(w, r)
		})
	}
}
