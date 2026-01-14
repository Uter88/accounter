package rest

import (
	"accounter/config"
	"accounter/pkg/logger"
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// HTTP server
type Server struct {
	http.Server
	cfg    config.Config
	logger logger.Logger
}

// Name server information
func (s *Server) Name() string {
	return fmt.Sprintf("HTTP server on 0.0.0.0:%d", s.cfg.HTTP.Port)
}

// Run start HTTP listener
func (s *Server) Run(ctx context.Context) error {
	return s.ListenAndServe()
}

// NewHTTPServer creates new Server instance
func NewHTTPServer(ctx context.Context, cfg config.Config, logger logger.Logger, handler http.Handler) *Server {
	return &Server{
		Server: http.Server{
			Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.HTTP.Port),
			Handler:      handler,
			BaseContext:  func(l net.Listener) context.Context { return ctx },
			ReadTimeout:  cfg.HTTP.ReadTimeout,
			WriteTimeout: cfg.HTTP.WriteTimeout,
		},
		cfg:    cfg,
		logger: logger.WithPerfix("HTTP"),
	}
}

// NewGinServer creates new gin.Engine server
func NewGinServer(cfg config.Config, engines ...Engine) http.Handler {
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

	r.NoRoute(func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(404, "Page not found")
	})

	for _, engine := range engines {
		engine.RegisterRoutes(r)
	}

	return r
}

// GIN server endpoints engine
type Engine interface {
	RegisterRoutes(e *gin.Engine)
}
