package restapi

import (
	"accounter/config"
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// HTTP server
type Server struct {
	http.Server
	cfg    config.Config
	logger config.Logger
}

// Creates new Server instance
func NewServer(ctx context.Context, cfg config.Config, logger config.Logger) *Server {
	if cfg.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.ContextWithFallback = true

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.HTTP.AllowOrigins,
		AllowWildcard:    cfg.HTTP.AllowWildcard,
		AllowMethods:     cfg.HTTP.AllowMethods,
		AllowHeaders:     cfg.HTTP.AllowHeaders,
		ExposeHeaders:    cfg.HTTP.ExposeHeaders,
		AllowCredentials: cfg.HTTP.AllowCredentials,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: cfg.HTTP.MaxAge,
	}))

	r.Use(gin.Recovery())

	return &Server{
		Server: http.Server{
			Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.HTTP.Port),
			Handler:      r,
			BaseContext:  func(l net.Listener) context.Context { return ctx },
			ReadTimeout:  time.Minute * 20,
			WriteTimeout: time.Minute * 20,
		},
		cfg:    cfg,
		logger: logger,
	}
}

type Engine interface {
	RegisterRoutes(e *gin.Engine)
}

func (s *Server) RegisterEngine(e Engine) {
	switch tp := s.Handler.(type) {
	case *gin.Engine:
		e.RegisterRoutes(tp)
	}
}
