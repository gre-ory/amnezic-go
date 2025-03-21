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
	"github.com/gre-ory/amnezic-go/internal/model"
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

	config := readConfig(ctx)

	//
	// logger
	//

	logger := NewLogger(config)
	logger = logger.With(zap.String("env", config.Env), zap.String("app", fmt.Sprintf("%s %s", config.App.Name, config.App.Version)))
	logger.Info("starting app...", zap.Any("config", config))
	// logger.Info("secrets...", zap.Any("secrets", secrets))

	//
	// servers
	//

	NewServer(logger, config).Run(ctx)

	os.Exit(0)
}

// //////////////////////////////////////////////////
// server

func NewServer(logger *zap.Logger, config *Config) *Server {
	return &Server{
		logger: logger,
		config: config,
	}
}

type Server struct {
	logger *zap.Logger
	config *Config
}

func (s *Server) Run(ctx context.Context) {

	//
	// default admin user
	//

	defaultAdminUser := &model.User{
		Name:     s.config.DefaultAdmin.Name,
		Password: s.config.DefaultAdmin.Password,
	}

	//
	// store
	//

	db, _ := sql.Open("sqlite3", s.config.Sqlite.DataSource)
	defer db.Close()

	musicFilter := s.config.MusicFileFilter(s.logger)
	imageFilter := s.config.ImageFileFilter(s.logger)

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
	userStore := store.NewUserStore(s.logger)
	sessionStore := store.NewSessionStore(s.logger)
	fileStore := store.NewFileStore(s.logger)

	musicFileValidator := fileStore.PathValidator(ctx, musicFilter)
	imageFileValidator := fileStore.PathValidator(ctx, imageFilter)

	//
	// client
	//

	deezerClient := client.NewDeezerClient(s.logger)
	downloadClient := client.NewDownloadClient(s.logger, musicFilter, imageFilter)

	//
	// service
	//

	gameService := service.NewGameService(s.logger, db, gameStore, gameQuestionStore, musicStore, artistStore, albumStore, themeStore, themeQuestionStore, deezerClient)
	musicService := service.NewMusicService(s.logger, deezerClient, downloadClient, db, musicStore, albumStore, artistStore, themeStore, themeQuestionStore, musicFileValidator, imageFileValidator)
	artistService := service.NewArtistService(s.logger, downloadClient, db, artistStore, musicStore, imageFileValidator)
	albumService := service.NewAlbumService(s.logger, downloadClient, db, albumStore, musicStore, imageFileValidator)
	themeService := service.NewThemeService(s.logger, db, themeStore, themeQuestionStore, musicStore, artistStore, albumStore)
	userService := service.NewUserService(s.logger, db, userStore, defaultAdminUser)
	sessionService := service.NewSessionService(s.logger, s.config.Session.SecretKey, db, sessionStore, userStore)
	fileService := service.NewFileService(s.logger, fileStore)

	//
	// api
	//

	gameHandler := api.NewGamehandler(s.logger, gameService)
	playlistHandler := api.NewPlaylisthandler(s.logger, musicService)
	musicHandler := api.NewMusichandler(s.logger, musicService, sessionService)
	artistHandler := api.NewArtisthandler(s.logger, artistService, sessionService)
	albumHandler := api.NewAlbumhandler(s.logger, albumService, sessionService)
	themeHandler := api.NewThemehandler(s.logger, themeService, musicService, sessionService)
	userHandler := api.NewUserHandler(s.logger, userService, sessionService)
	sessionHandler := api.NewSessionhandler(s.logger, sessionService)
	fileHandler := api.NewFilehandler(s.logger, musicFilter, imageFilter, fileService, sessionService)

	//
	// router
	//

	router := httprouter.New()
	gameHandler.RegisterRoutes(router)
	musicHandler.RegisterRoutes(router)
	artistHandler.RegisterRoutes(router)
	albumHandler.RegisterRoutes(router)
	themeHandler.RegisterRoutes(router)
	playlistHandler.RegisterRoutes(router)
	userHandler.RegisterRoutes(router)
	sessionHandler.RegisterRoutes(router)
	fileHandler.RegisterRoutes(router)

	//
	// server
	//

	server := http.Server{
		Addr: s.config.Server.Address,
		Handler: AllowCORS(s.logger, s.config.Server.WhiteListOrigins)(
			WithRequestLogging(s.logger)(
				router,
			),
		),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	s.logger.Info(fmt.Sprintf("starting backend server on %s", server.Addr))
	err := server.ListenAndServeTLS(s.config.Server.CrtFile, s.config.Server.KeyFile)
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

func NewLogger(config *Config) *zap.Logger {

	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.Log.File,
		MaxSize:    100, // megabytes
		MaxBackups: 10,
		MaxAge:     30, // days
	})

	var zapEncoder zapcore.Encoder
	switch strings.ToLower(config.Log.Encoder) {
	case "prd":
		cfg := zap.NewProductionEncoderConfig()
		zapEncoder = zapcore.NewJSONEncoder(cfg)
	case "dev":
		cfg := zap.NewDevelopmentEncoderConfig()
		zapEncoder = zapcore.NewConsoleEncoder(cfg)
	default:
		fmt.Printf("invalid LOG_CONFIG >>> FALLBACK to 'prd'!\n")
		cfg := zap.NewProductionEncoderConfig()
		zapEncoder = zapcore.NewJSONEncoder(cfg)
	}

	var zapLevel zapcore.Level
	switch strings.ToLower(config.Log.Level) {
	case "err", "error":
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
		zapEncoder,
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

func AllowCORS(logger *zap.Logger, whitelistOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowed := util.Contains(whitelistOrigins, origin)
			if allowed {
				logger.Info(fmt.Sprintf("[COR] OK - Origin: %s", origin))
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "*")
				w.Header().Set("Access-Control-Allow-Headers", "content-type,authorization")
			} else {
				logger.Info(fmt.Sprintf("[COR] BLOCKED - Origin: %s", origin))
			}
			next.ServeHTTP(w, r)
		})
	}
}
