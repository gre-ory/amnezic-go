package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
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
	yaml "gopkg.in/yaml.v2"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

// //////////////////////////////////////////////////
// main

func main() {

	ctx := context.Background()

	//
	// config
	//

	config := readConfig()
	secrets := readSecrets()

	//
	// logger
	//

	logger := NewLogger(config.Log)
	logger = logger.With(zap.String("env", config.Env), zap.String("app", config.App), zap.String("version", config.Version))
	logger.Info("starting app...", zap.Any("config", config))
	// logger.Info("secrets...", zap.Any("secrets", secrets))

	//
	// servers
	//

	NewBackendServer(logger, config, secrets).Run(ctx)

	os.Exit(0)
}

// //////////////////////////////////////////////////
// backend server

func NewBackendServer(logger *zap.Logger, config *Config, secrets *Secrets) *BackendServer {
	return &BackendServer{
		logger:  logger.With(zap.String("server", "backend")),
		config:  config,
		secrets: secrets,
	}
}

type BackendServer struct {
	logger  *zap.Logger
	config  *Config
	secrets *Secrets
}

func (s *BackendServer) Run(ctx context.Context) {

	//
	// default admin user
	//

	defaultAdminUser := &model.User{
		Name:     s.secrets.DefaultAdmin.Name,
		Password: s.secrets.DefaultAdmin.Password,
	}

	//
	// store
	//

	db, _ := sql.Open("sqlite3", s.config.Store.SqliteDataSource)
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
	userStore := store.NewUserStore(s.logger)
	sessionStore := store.NewSessionStore(s.logger)

	//
	// client
	//

	deezerClient := client.NewDeezerClient(s.logger)

	//
	// service
	//

	gameService := service.NewGameService(s.logger, db, gameStore, gameQuestionStore, musicStore, artistStore, albumStore, themeStore, themeQuestionStore, deezerClient)
	musicService := service.NewMusicService(s.logger, deezerClient, db, musicStore, albumStore, artistStore, themeStore, themeQuestionStore)
	themeService := service.NewThemeService(s.logger, db, themeStore, themeQuestionStore, musicStore, artistStore, albumStore)
	userService := service.NewUserService(s.logger, db, userStore, defaultAdminUser)
	sessionService := service.NewSessionService(s.logger, s.secrets.SessionSecretKey, db, sessionStore, userStore)

	//
	// api
	//

	gameHandler := api.NewGamehandler(s.logger, gameService)
	playlistHandler := api.NewPlaylisthandler(s.logger, musicService)

	musicHandler := api.NewMusichandler(s.logger, musicService, sessionService)
	themeHandler := api.NewThemehandler(s.logger, themeService, musicService, sessionService)
	userHandler := api.NewUserHandler(s.logger, userService, sessionService)
	sessionHandler := api.NewSessionhandler(s.logger, sessionService)

	//
	// router
	//

	router := httprouter.New()
	gameHandler.RegisterRoutes(router)
	musicHandler.RegisterRoutes(router)
	themeHandler.RegisterRoutes(router)
	playlistHandler.RegisterRoutes(router)
	userHandler.RegisterRoutes(router)
	sessionHandler.RegisterRoutes(router)

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

func NewLogger(config LogConfig) *zap.Logger {

	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.File,
		MaxSize:    100, // megabytes
		MaxBackups: 10,
		MaxAge:     30, // days
	})

	var zapEncoderCfg zapcore.EncoderConfig
	switch config.Encoder {
	case "prd":
		zapEncoderCfg = zap.NewProductionEncoderConfig()
	case "dev":
		zapEncoderCfg = zap.NewDevelopmentEncoderConfig()
	default:
		fmt.Printf("invalid LOG_CONFIG >>> FALLBACK to 'prd'!\n")
		zapEncoderCfg = zap.NewProductionEncoderConfig()
	}

	var zapLevel zapcore.Level
	switch config.Level {
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

// var whitelistOrigins []string = []string{
// 	"http://localhost:3000",
// 	"http://localhost:9090",
// 	"http://158.178.206.68:8080",
// 	"http://158.178.206.68:8081",
// 	"http://158.178.206.68:8082",
// }

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

// //////////////////////////////////////////////////
// config

type Config struct {
	Env     string       `yaml:"env"`
	App     string       `yaml:"app"`
	Version string       `yaml:"version"`
	Log     LogConfig    `yaml:"log"`
	Server  ServerConfig `yaml:"server"`
	Store   StoreConfig  `yaml:"store"`
}

type LogConfig struct {
	Encoder string `yaml:"encoder"`
	Level   string `yaml:"level"`
	File    string `yaml:"file"`
}

type ServerConfig struct {
	Address          string   `yaml:"address"`
	WhiteListOrigins []string `yaml:"white-list-origins"`
}

type StoreConfig struct {
	SqliteDataSource string `yaml:"sqlite-data-source"`
}

func readConfig() *Config {

	path := os.Getenv("CONFIG_FILE")
	if path == "" {
		panic(fmt.Errorf("missing CONFIG_FILE env variable"))
	}

	_, err := os.Stat(path)
	if err != nil {
		panic(fmt.Errorf("missing config file: %s", err.Error()))
	}

	file, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("unable to read config file %s: %s", path, err.Error()))
	}
	fmt.Printf("\n\n ----- %s ----- \n%s\n\n", path, string(file))

	config := Config{}
	err = yaml.UnmarshalStrict(file, &config)
	if err != nil {
		panic(fmt.Errorf("unable to decode config file %s: %s", path, err.Error()))
	}

	// replace env variables
	config.Log.File = replaceEnvVariables(config.Log.File)
	config.Store.SqliteDataSource = replaceEnvVariables(config.Store.SqliteDataSource)

	return &config
}

// //////////////////////////////////////////////////
// secrets

type Secrets struct {
	SessionSecretKey string       `yaml:"session-secret-key"`
	DefaultAdmin     DefaultAdmin `yaml:"default-admin"`
}

type DefaultAdmin struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

func readSecrets() *Secrets {

	path := os.Getenv("SECRET_FILE")
	if path == "" {
		panic(fmt.Errorf("missing SECRET_FILE env variable"))
	}

	stats, err := os.Stat(path)
	if err != nil {
		panic(fmt.Errorf("missing secret file: %s", err.Error()))
	}

	permissions := stats.Mode().Perm()
	if permissions != 0o600 {
		panic(fmt.Errorf("incorrect permissions for secret file %s (0%o), must be 0600 for '%s'", permissions, permissions, path))
	}

	file, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("unable to read secret file %s: %s", path, err.Error()))
	}

	secrets := &Secrets{}
	err = yaml.UnmarshalStrict(file, secrets)
	if err != nil {
		panic(fmt.Errorf("unable to decode secret file %s: %s", path, err.Error()))
	}

	return secrets
}

// //////////////////////////////////////////////////
// env variable

func replaceEnvVariables(value string) string {
	regexp := regexp.MustCompile(`\$([A-Z_]+)\$`)
	matches := regexp.FindAllStringSubmatch(value, -1)
	for _, match := range matches {
		envVariable := match[1]
		envValue := os.Getenv(envVariable)
		value = strings.ReplaceAll(value, match[0], envValue)
	}
	return value
}
