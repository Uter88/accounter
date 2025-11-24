package main

import (
	"accounter/config"
	"accounter/frontend"
)

func main() {
	// Init graceful shutdown context
	ctx, cancel := config.InitGracefulShutdownCtx()

	// Init application config
	cfg := config.InitConfig()

	// Create logger
	logger := config.NewLogger(cfg.DebugMode, cfg.AppMode, "logs")

	// Create frontend application instance
	frontApp := frontend.NewApp(cfg, *logger)

	// Defer canceling and shutdown application
	defer cancel()

	// Run frontend application
	frontApp.Run(ctx)
}
