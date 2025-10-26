package main

import (
	"accounter/backend/core"
	"accounter/backend/server"
	"accounter/config"
	"accounter/frontend"
)

func main() {
	ctx, cancel := config.InitGracefulShutdownCtx()
	cfg := config.InitConfig()

	logger := config.NewLogger(cfg.DebugMode, cfg.AppMode, "logs")

	frontApp := frontend.NewApp(cfg, *logger)
	backendApp := core.NewApp(cfg, *logger).Init(ctx)

	server := server.NewServer(backendApp)

	backendApp.RegisterTask("Backend HTTP server", &server)
	backendApp.RegisterTask("Frontend HTTP server", &frontApp)

	defer func() {
		cancel()
		backendApp.Shutdown()
	}()

	backendApp.Run(ctx)
}
