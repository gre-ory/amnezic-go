package main

import (
	"context"
	"log"
	"path/filepath"

	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// config

type Config struct {
	Env string `env:"ENVIRONMENT,required"`
	App struct {
		Name    string `env:"NAME,required"`
		Version string `env:"VERSION,required"`
	} `env:",prefix=APP_"`
	Log struct {
		Encoder string `env:"ENCODER,required"`
		Level   string `env:"LEVEL,required"`
		File    string `env:"FILE,required"`
	} `env:",prefix=LOG_"`
	Server struct {
		KeyFile          string   `env:"KEY_FILE,required"`
		CrtFile          string   `env:"CRT_FILE,required"`
		Address          string   `env:"ADDRESS,required"`
		WhiteListOrigins []string `env:"WHITE_LIST_ORIGINS"`
	} `env:",prefix=SERVER_"`
	Sqlite struct {
		DataSource string `env:"DATA_SOURCE,required"`
	} `env:",prefix=SQLITE_"`
	DefaultAdmin struct {
		Name     string `env:"NAME,required"`
		Password string `env:"PASSWORD,required"`
	} `env:",prefix=DEFAULT_ADMIN_"`
	Session struct {
		SecretKey string `env:"SECRET_KEY,required"`
	} `env:",prefix=SESSION_"`
	Static struct {
		Directory string `env:"DIRECTORY,required"`
		Music     struct {
			Directory  string   `env:"DIRECTORY,required"`
			Extensions []string `env:"EXTENSIONS"`
		} `env:",prefix=MUSIC_"`
		Image struct {
			Directory  string   `env:"DIRECTORY,required"`
			Extensions []string `env:"EXTENSIONS"`
		} `env:",prefix=IMAGE_"`
	} `env:",prefix=STATIC_"`
}

func readConfig(ctx context.Context) *Config {
	var config Config
	if err := envconfig.Process(ctx, &config); err != nil {
		log.Fatalf("failed to instantiate internal config: %s\n", err)
	}
	return &config
}

func (c *Config) MusicFileFilter(logger *zap.Logger) *model.FileFilter {
	directory := filepath.Join(c.Static.Directory, c.Static.Music.Directory)
	filter, err := model.NewFileFilter(directory, c.Static.Music.Extensions)
	if err != nil {
		logger.Fatal("invalid music filter", zap.Error(err))
	}
	return filter
}

func (c *Config) ImageFileFilter(logger *zap.Logger) *model.FileFilter {
	directory := filepath.Join(c.Static.Directory, c.Static.Image.Directory)
	filter, err := model.NewFileFilter(directory, c.Static.Image.Extensions)
	if err != nil {
		logger.Fatal("invalid image filter", zap.Error(err))
	}
	return filter
}
