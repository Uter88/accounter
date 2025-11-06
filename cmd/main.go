package main

import (
	"accounter/backend/core"
	"accounter/backend/server"
	"accounter/config"
	"accounter/frontend"
)

func main() {
	// Init graceful shutdown context
	ctx, cancel := config.InitGracefulShutdownCtx()

	// Init config
	cfg := config.InitConfig()

	// Create logger
	logger := config.NewLogger(cfg.DebugMode, cfg.AppMode, "logs")

	// Create frontend application instance
	frontApp := frontend.NewApp(cfg, *logger)

	// Create backend application instance and init it
	backendApp := core.NewApp(cfg, *logger).Init(ctx)

	// Create HTTP server instance
	server := server.NewServer(backendApp)

	// Register background tasks
	backendApp.RegisterTask("Backend HTTP server", &server)
	backendApp.RegisterTask("Frontend HTTP server", &frontApp)

	// Defer canceling and shutdown application
	defer func() {
		cancel()
		backendApp.Shutdown()
	}()

	// Run backend application
	backendApp.Run(ctx)
}
