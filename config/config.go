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

// Application config
type Config struct {

	// Debug mode flag
	DebugMode bool

	// Application mode: prod, dev, etc
	AppMode string

	// JWT secret key salt
	SecretKey string `yaml:"secret_key"`

	// Frontend config
	Client struct {
		Port uint `yaml:"port"`
	} `yaml:"client"`

	// Backend HTTP config
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

	// Database config
	DB struct {
		Driver string `yaml:"driver"`
		DSN    string `yaml:"dsn"`
	} `yaml:"db"`
}

// InitConfig parse args and load config from YAML file
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

// InitGracefulShutdownCtx creates graceful shutdown context and cancel function
func InitGracefulShutdownCtx() (context.Context, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	return ctx, cancel
}
