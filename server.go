package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gre-ory/amnezic-go/internal/api"
	"github.com/gre-ory/amnezic-go/internal/client"
	"github.com/gre-ory/amnezic-go/internal/service"
	"github.com/gre-ory/amnezic-go/internal/store"
	"github.com/gre-ory/amnezic-go/internal/store/legacy"
	"github.com/gre-ory/amnezic-go/internal/store/memory"
	"github.com/gre-ory/amnezic-go/internal/util"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

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

	NewBackendServer(logger).Run(ctx)

	os.Exit(0)
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
	if address == "" {
		s.logger.Warn("missing BACKEND_ADDRESS")
		return
	}

	dataSource := os.Getenv("SQLITE_DATA_SOURCE")
	if address == "" {
		s.logger.Warn("missing SQLITE_DATA_SOURCE")
		return
	}

	//
	// store
	//

	db, _ := sql.Open("sqlite3", dataSource)
	defer db.Close()

	gameStore := memory.NewGameMemoryStore()
	gameQuestionStore := legacy.NewGameQuestionLegacyStore(s.logger, legacy.RootPath_FreeDotFr)

	// musicStore := store.NewMusicMemoryStore()
	// albumStore := store.NewMusicAlbumMemoryStore()
	// artistStore := store.NewMusicArtistMemoryStore()

	musicStore := store.NewMusicStore(s.logger)
	albumStore := store.NewMusicAlbumStore(s.logger)
	artistStore := store.NewMusicArtistStore(s.logger)

	themeStore := store.NewThemeStore(s.logger)
	themeQuestionStore := store.NewThemeQuestionStore(s.logger)

	//
	// client
	//

	deezerClient := client.NewDeezerClient(s.logger)

	//
	// service
	//

	gameService := service.NewGameService(s.logger, db, gameStore, gameQuestionStore, musicStore, artistStore, albumStore, themeStore, themeQuestionStore)
	musicService := service.NewMusicService(s.logger, deezerClient, db, musicStore, albumStore, artistStore, themeStore, themeQuestionStore)
	themeService := service.NewThemeService(s.logger, db, themeStore, themeQuestionStore, musicStore, artistStore, albumStore)

	//
	// api
	//

	gameHandler := api.NewGamehandler(s.logger, gameService)
	musicHandler := api.NewMusichandler(s.logger, musicService)
	themeHandler := api.NewThemehandler(s.logger, themeService, musicService)

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

	logFile := os.Getenv("LOG_FILE")
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100, // megabytes
		MaxBackups: 10,
		MaxAge:     30, // days
	})
	fmt.Printf("logFile: %s\n", logFile)

	var zapEncoderCfg zapcore.EncoderConfig
	switch os.Getenv("LOG_CONFIG") {
	case "prd":
		zapEncoderCfg = zap.NewProductionEncoderConfig()
	case "dev":
		zapEncoderCfg = zap.NewDevelopmentEncoderConfig()
	default:
		fmt.Printf("invalid LOG_CONFIG >>> FALLBACK to 'prd'!\n")
		zapEncoderCfg = zap.NewProductionEncoderConfig()
	}

	var zapLevel zapcore.Level
	switch os.Getenv("LOG_LEVEL") {
	case "err":
		zapLevel = zap.ErrorLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "debug":
		zapLevel = zap.DebugLevel
	default:
		fmt.Printf("invalid LOG_LEVEL >>> FALLBACK to 'info'!\n")
		zapLevel = zap.InfoLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapEncoderCfg),
		writer,
		zapLevel,
	)

	logger := zap.New(core)
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
	"http://158.178.206.68:8080",
	"http://158.178.206.68:8081",
	"http://158.178.206.68:8082",
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
				w.Header().Set("Access-Control-Allow-Headers", "content-type")
			} else {
				logger.Info(fmt.Sprintf("[COR] BLOCKED - Origin: %s", origin))
			}
			next.ServeHTTP(w, r)
		})
	}
}
