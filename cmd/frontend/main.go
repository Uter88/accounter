package main

import (
	"accounter/config"
	"accounter/frontend"
)

func main() {
	ctx, cancel := config.InitGracefulShutdownCtx()
	cfg := config.InitConfig()

	logger := config.NewLogger(cfg.DebugMode, cfg.AppMode, "logs")

	frontApp := frontend.NewApp(cfg, *logger)

	defer cancel()

	frontApp.Run(ctx)
}
