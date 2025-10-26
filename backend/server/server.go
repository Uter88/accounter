package server

import (
	"accounter/backend/core"
	v1 "accounter/backend/server/handlers/v1"
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	*core.App
}

func NewServer(app *core.App) Server {
	return Server{
		App: app,
	}
}

func (s *Server) Run(ctx context.Context) error {
	if s.Config.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.ContextWithFallback = true

	r.Use(cors.New(cors.Config{
		AllowOrigins:     s.Config.HTTP.AllowOrigins,
		AllowWildcard:    s.Config.HTTP.AllowWildcard,
		AllowMethods:     s.Config.HTTP.AllowMethods,
		AllowHeaders:     s.Config.HTTP.AllowHeaders,
		ExposeHeaders:    s.Config.HTTP.ExposeHeaders,
		AllowCredentials: s.Config.HTTP.AllowCredentials,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: s.Config.HTTP.MaxAge,
	}))

	r.Use(gin.Recovery())

	v1.NewEngine(s.App).RegisterRoutes(r)

	serv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", s.Config.HTTP.Port),
		Handler:      r,
		BaseContext:  func(l net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Minute * 20,
		WriteTimeout: time.Minute * 20,
	}

	s.Logger.Infof("Start backend HTTP server on %d port", s.Config.HTTP.Port)

	return serv.ListenAndServe()
}
