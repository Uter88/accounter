package main

import (
	"accounter/config"
	"accounter/internal/app"
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

	// Create AppContext instance and init it
	backendApp := app.NewAppContext(ctx, cfg, logger).Init(ctx)

	// Register background tasks
	backendApp.RegisterTask("Frontend HTTP server", &frontApp)

	// Defer canceling and shutdown application
	defer func() {
		cancel()
		backendApp.Shutdown()
	}()

	// Run backend application
	backendApp.Run(ctx)
}
