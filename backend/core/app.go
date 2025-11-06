package core

import (
	"accounter/adapters/adapter_sql"
	"accounter/config"
	"context"
	"log"
)

// Common application instance
type App struct {

	// Application config
	Config config.Config

	// Default logger
	Logger config.Logger

	// Database client
	DbClient adapter_sql.SQLClient

	// Background tasks map
	tasks map[string]task
}

// Creates new App
func NewApp(config config.Config, logger config.Logger) *App {
	return &App{
		Config:   config,
		Logger:   logger,
		tasks:    make(map[string]task),
		DbClient: adapter_sql.NewSQLClient(config.DB.Driver, config.DB.DSN),
	}
}

// Init application, connections, etc.
func (a *App) Init(ctx context.Context) *App {
	if err := a.DbClient.Connect(ctx); err != nil {
		a.Logger.Fatalln(err)
	}

	if err := a.DbClient.Migrate(ctx); err != nil {
		a.Logger.Fatalln(err)
	}

	return a
}

// Shutdown application
func (a *App) Shutdown() {
	if err := a.DbClient.Disconnect(); err != nil {
		a.Logger.Errorf("Error disconnect from db: %s", err.Error())
	}

	a.Logger.Info("Shutdown system")
}

// Background task
type task interface {
	Run(ctx context.Context) error
}

// RegisterTask add background task to local store
func (a *App) RegisterTask(name string, task task) *App {
	a.tasks[name] = task

	return a
}

// Run application
func (a *App) Run(ctx context.Context) {
	a.launchTasks(ctx)

	<-ctx.Done()
}

// Launch all background tasks in goroutines
func (a *App) launchTasks(ctx context.Context) *App {
	for name, task := range a.tasks {
		go a.launchTask(ctx, name, task)
	}

	return a
}

// Launch background task
func (a *App) launchTask(ctx context.Context, name string, task task) {
	log.Printf("Launch task: %s", name)

	if err := task.Run(ctx); err != nil {
		log.Printf("Error launch task %s: %s\n", name, err.Error())
	}
}
