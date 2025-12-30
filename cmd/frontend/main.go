package main

import (
	"accounter/config"
	wasmgoap "accounter/internal/infrastructure/ui/wasm-goap"
	"accounter/pkg/logger"
)

func main() {
	// Init graceful shutdown context
	ctx, cancel := config.InitGracefulShutdownCtx()

	// Init application config
	cfg := config.InitConfig()

	// Create logger
	logger := logger.NewLogger(cfg.DebugMode, cfg.AppMode, "logs")

	// Create frontend application instance
	frontApp := wasmgoap.NewApp(cfg, logger)

	// Defer canceling and shutdown application
	defer cancel()

	// Run frontend application
	frontApp.Run(ctx)
}
