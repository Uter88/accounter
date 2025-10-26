package config

import (
	"context"
	_ "embed"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

//go:embed config.yaml
var configFile []byte

type Config struct {
	DebugMode bool
	AppMode   string

	SecretKey string `yaml:"secret_key"`

	Client struct {
		Port uint `yaml:"port"`
	} `yaml:"client"`

	HTTP struct {
		Host             string        `yaml:"host"`
		Port             uint          `yaml:"port"`
		AllowOrigins     []string      `yaml:"allow_origins"`
		AllowHeaders     []string      `yaml:"allow_headers"`
		AllowMethods     []string      `yaml:"allow_methods"`
		AllowWildcard    bool          `yaml:"allow_wildcard"`
		AllowCredentials bool          `yaml:"allow_credentials"`
		ExposeHeaders    []string      `yaml:"expose_headers"`
		MaxAge           time.Duration `yaml:"max_age"`
	} `yaml:"http"`
	DB struct {
		DSN string `yaml:"dsn"`
	} `yaml:"db"`
}

func InitConfig() (cfg Config) {
	flag.BoolVar(&cfg.DebugMode, "debug", true, "Debug mode")
	flag.StringVar(&cfg.AppMode, "mode", "dev", "App mode")
	flag.UintVar(&cfg.HTTP.Port, "backend-port", 8001, "HTTP BACKEND PORT")
	flag.UintVar(&cfg.Client.Port, "frontend-port", 8000, "HTTP FRONTEND PORT")

	flag.Parse()

	configs := make(map[string]Config)

	if err := yaml.Unmarshal(configFile, &configs); err != nil {
		log.Fatalf("Error init config: %s", err.Error())
	}

	if c, ok := configs[cfg.AppMode]; ok {
		c.AppMode = cfg.AppMode
		c.DebugMode = cfg.DebugMode
		c.HTTP.Port = cfg.HTTP.Port
		return c
	}

	log.Fatalln("Error init config: config not found")
	return
}

func InitGracefulShutdownCtx() (context.Context, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	return ctx, cancel
}
