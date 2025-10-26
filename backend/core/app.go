package core

import (
	"accounter/backend/db"
	"accounter/config"
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

type App struct {
	Config   config.Config
	Logger   config.Logger
	DbClient *sqlx.DB

	tasks map[string]task
}

func NewApp(config config.Config, logger config.Logger) *App {
	return &App{Config: config, Logger: logger, tasks: make(map[string]task)}
}

func (a *App) Init(ctx context.Context) *App {
	a.DbClient = db.InitConnection(ctx, a.Config.DB.DSN)

	return a
}

func (a *App) Shutdown() {
	if a.DbClient != nil {
		a.DbClient.Close()
	}

	log.Println("Shutdown system")
}

type task interface {
	Run(ctx context.Context) error
}

func (a *App) RegisterTask(name string, task task) *App {
	a.tasks[name] = task

	return a
}

func (a *App) launchTasks(ctx context.Context) *App {
	for name, task := range a.tasks {
		go a.launchTask(ctx, name, task)
	}

	return a
}

func (a *App) launchTask(ctx context.Context, name string, task task) {
	log.Printf("Launch task: %s", name)

	if err := task.Run(ctx); err != nil {
		log.Printf("Error launch task %s: %s\n", name, err.Error())
	}
}

func (a *App) Run(ctx context.Context) {
	a.launchTasks(ctx)

	<-ctx.Done()
}
