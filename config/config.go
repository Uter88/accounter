package config

import (
	"accounter/pkg/tools"
	"context"
	_ "embed"
	"flag"
	"log"
	"os/signal"
	"strings"
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
	SecretKey   string        `yaml:"secret_key"`
	TokenExpire time.Duration `yaml:"token_expire"`

	// Frontend config
	Client struct {
		Port uint `yaml:"port"`
	} `yaml:"client"`

	// Backend HTTP config
	HTTP HTTPConfig `yaml:"http"`

	// Database config
	DB DBConfig `yaml:"db"`

	Kafka struct {
		Brokers      []string      `yaml:"brokers"`
		Topic        string        `yaml:"topic"`
		Group        string        `yaml:"group"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
	} `yaml:"kafka"`
}

type DBConfig struct {
	Driver   string `yaml:"driver"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DbName   string `yaml:"dbname"`
	SSLMode  string `yaml:"ssl_mode"`

	DSN string `yaml:"dsn"`
}

// HTTP server config
type HTTPConfig struct {
	Host             string          `yaml:"host"`
	Port             uint            `yaml:"port"`
	AllowOrigins     []string        `yaml:"allow_origins"`
	AllowHeaders     []string        `yaml:"allow_headers"`
	AllowMethods     []string        `yaml:"allow_methods"`
	AllowWildcard    bool            `yaml:"allow_wildcard"`
	AllowCredentials bool            `yaml:"allow_credentials"`
	ExposeHeaders    []string        `yaml:"expose_headers"`
	MaxAge           time.Duration   `yaml:"max_age"`
	ReadTimeout      time.Duration   `yaml:"read_timeout"`
	WriteTimeout     time.Duration   `yaml:"write_timeout"`
	Websocket        WebsocketConfig `yaml:"websocket"`
}

// Websocket config
type WebsocketConfig struct {
	URL                  string        `yaml:"url"`
	ReadDeadline         time.Duration `yaml:"read_deadline"`
	WriteDeadline        time.Duration `yaml:"write_deadline"`
	ReconnectionInterval time.Duration `yaml:"reconnection_interval"`
	PingInterval         time.Duration `yaml:"ping_interval"`
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

		c.parseEnv()
		return c
	}

	log.Fatalln("Error init config: config not found")
	return
}

// parseEnv parse environment params
func (c *Config) parseEnv() {
	if port, ok := tools.GetEnv[uint]("server-port"); ok {
		c.HTTP.Port = port
	}

	if port, ok := tools.GetEnv[uint]("client-port"); ok {
		c.Client.Port = port
	}

	if user, ok := tools.GetEnv[string]("dbuser"); ok {
		c.DB.User = user
	}

	if password, ok := tools.GetEnv[string]("dbpassword"); ok {
		c.DB.Password = password
	}

	if host, ok := tools.GetEnv[string]("dbhost"); ok {
		c.DB.Host = host
	}

	if port, ok := tools.GetEnv[int]("dbport"); ok {
		c.DB.Port = port
	}

	if dbName, ok := tools.GetEnv[string]("dbname"); ok {
		c.DB.DbName = dbName
	}

	if brokers, ok := tools.GetEnv[string]("qbrokers"); ok {
		c.Kafka.Brokers = strings.Split(brokers, ",")
	}

	if topic, ok := tools.GetEnv[string]("qtopic"); ok {
		c.Kafka.Topic = topic
	}

	if group, ok := tools.GetEnv[string]("qgroup"); ok {
		c.Kafka.Group = group
	}
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
