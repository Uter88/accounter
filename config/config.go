package config

import (
	"accounter/pkg/utils"
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
	SecretKey   string        `yaml:"secret_key,omitempty"`
	TokenExpire time.Duration `yaml:"token_expire,omitempty"`

	// Frontend config
	Client struct {
		Port uint `yaml:"port"`
	} `yaml:"client"`

	// Backend HTTP config
	HTTP HTTPConfig `yaml:"http"`

	// Database config
	DB DBConfig `yaml:"db"`

	Kafka struct {
		Brokers      []string      `yaml:"brokers,omitempty"`
		Topic        string        `yaml:"topic,omitempty"`
		Group        string        `yaml:"group,omitempty"`
		ReadTimeout  time.Duration `yaml:"read_timeout,omitempty"`
		WriteTimeout time.Duration `yaml:"write_timeout,omitempty"`
		AutoCommit   bool          `yaml:"auto_commit,omitempty"`
	} `yaml:"kafka"`
}

type DBConfig struct {
	Driver   string `yaml:"driver,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
	Host     string `yaml:"host,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	DbName   string `yaml:"dbname,omitempty"`
	SSLMode  string `yaml:"ssl_mode,omitempty"`

	DSN string `yaml:"dsn"`
}

// HTTP server config
type HTTPConfig struct {
	Host             string          `yaml:"host,omitempty"`
	Port             uint            `yaml:"port,omitempty"`
	AllowOrigins     []string        `yaml:"allow_origins,omitempty"`
	AllowHeaders     []string        `yaml:"allow_headers,omitempty"`
	AllowMethods     []string        `yaml:"allow_methods,omitempty"`
	AllowWildcard    bool            `yaml:"allow_wildcard,omitempty"`
	AllowCredentials bool            `yaml:"allow_credentials,omitempty"`
	ExposeHeaders    []string        `yaml:"expose_headers,omitempty"`
	MaxAge           time.Duration   `yaml:"max_age,omitempty"`
	ReadTimeout      time.Duration   `yaml:"read_timeout,omitempty"`
	WriteTimeout     time.Duration   `yaml:"write_timeout,omitempty"`
	Websocket        WebsocketConfig `yaml:"websocket"`
}

// Websocket config
type WebsocketConfig struct {
	URL                  string        `yaml:"url,omitempty"`
	ReadDeadline         time.Duration `yaml:"read_deadline,omitempty"`
	WriteDeadline        time.Duration `yaml:"write_deadline,omitempty"`
	ReconnectionInterval time.Duration `yaml:"reconnection_interval,omitempty"`
	PingInterval         time.Duration `yaml:"ping_interval,omitempty"`
}

// InitConfig parse args and load config from YAML file
func InitConfig() Config {
	cfg := newDefaultConfig()

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
	return cfg
}

// newDefaultConfig init default Config
func newDefaultConfig() (cfg Config) {
	cfg.DebugMode = true
	cfg.AppMode = "dev"
	cfg.HTTP.Port = 8001
	cfg.Client.Port = 8000
	cfg.DB.User = "postgres"
	cfg.DB.Password = ""
	cfg.DB.Host = "localhost"
	cfg.DB.Port = 5432
	cfg.DB.DbName = "main"
	cfg.Kafka.Brokers = []string{"kafka:9092"}
	cfg.Kafka.Topic = "accounter"
	cfg.Kafka.Group = "accounter"
	return
}

// parseEnv parse environment params
func (c *Config) parseEnv() {
	if key, ok := utils.GetEnv[string]("secret-key"); ok {
		c.SecretKey = key
	}

	if port, ok := utils.GetEnv[uint]("server-port"); ok {
		c.HTTP.Port = port
	}

	if port, ok := utils.GetEnv[uint]("client-port"); ok {
		c.Client.Port = port
	}

	if user, ok := utils.GetEnv[string]("dbuser"); ok {
		c.DB.User = user
	}

	if password, ok := utils.GetEnv[string]("dbpassword"); ok {
		c.DB.Password = password
	}

	if host, ok := utils.GetEnv[string]("dbhost"); ok {
		c.DB.Host = host
	}

	if port, ok := utils.GetEnv[int]("dbport"); ok {
		c.DB.Port = port
	}

	if dbName, ok := utils.GetEnv[string]("dbname"); ok {
		c.DB.DbName = dbName
	}

	if brokers, ok := utils.GetEnv[string]("qbrokers"); ok {
		c.Kafka.Brokers = strings.Split(brokers, ",")
	}

	if topic, ok := utils.GetEnv[string]("qtopic"); ok {
		c.Kafka.Topic = topic
	}

	if group, ok := utils.GetEnv[string]("qgroup"); ok {
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
