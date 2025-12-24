package app

import (
	"accounter/adapters/adapter_sql"
	"accounter/adapters/auth"
	"accounter/adapters/renderers"
	restapi "accounter/adapters/rest_api"
	v1 "accounter/adapters/rest_api/controllers/v1"
	"accounter/config"
	"accounter/internal/domain/task"
	"accounter/internal/domain/user"
	"context"
)

// Application context
type AppContext struct {
	server *restapi.Server

	config config.Config

	// Default logger
	logger config.Logger

	// Database client
	db adapter_sql.SQLClient

	// Background tasks map
	tasks map[string]bgTask
}

// Background task
type bgTask interface {
	Run(ctx context.Context) error
}

// RegisterTask add background task to local store
func (a *AppContext) RegisterTask(name string, task bgTask) *AppContext {
	a.tasks[name] = task

	return a
}

// Run application
func (a *AppContext) Run(ctx context.Context) {
	a.launchTasks(ctx)

	go a.server.ListenAndServe()

	<-ctx.Done()
}

// Launch all background tasks in goroutines
func (a *AppContext) launchTasks(ctx context.Context) *AppContext {
	for name, task := range a.tasks {
		go a.launchTask(ctx, name, task)
	}

	return a
}

// Launch background task
func (a *AppContext) launchTask(ctx context.Context, name string, task bgTask) {
	a.logger.Infof("Launch task: %s", name)

	if err := task.Run(ctx); err != nil {
		a.logger.Fatalf("Error launch task %s: %s\n", name, err.Error())
	}
}

// Init application, connections, etc.
func (a *AppContext) Init(ctx context.Context) *AppContext {
	return a.initConnections(ctx).initServer()
}

// Init connections with databases
func (a *AppContext) initConnections(ctx context.Context) *AppContext {
	if err := a.db.Connect(ctx); err != nil {
		a.logger.Fatalln(err)
	}

	return a
}

// Init HTTP server
func (a *AppContext) initServer() *AppContext {
	engine := v1.NewEngine(a.config, a.logger)

	userRepo := adapter_sql.NewUserRepository(a.db)
	taskRepo := adapter_sql.NewTaskRepository(a.db)

	engine.AuthService = auth.NewAuthService(userRepo)
	engine.UserService = user.NewUserService(userRepo)
	engine.TaskService = task.NewTaskService(taskRepo, renderers.NewTaskRenderer())

	a.server.RegisterEngine(engine)

	return a
}

// Shutdown application
func (a *AppContext) Shutdown() {
	if err := a.db.Disconnect(); err != nil {
		a.logger.Errorf("Error disconnect from db: %s", err.Error())
	} else {
		a.logger.Info("Success disconnected from db")
	}

	a.logger.Info("Shutdown system")
}

// NewAppContext creates new AppContext
func NewAppContext(ctx context.Context, cfg config.Config, logger config.Logger) *AppContext {
	return &AppContext{
		config: cfg,
		logger: logger,
		db:     adapter_sql.NewSQLClient(cfg.DB.Driver, cfg.DB.DSN),
		server: restapi.NewServer(ctx, cfg, logger),
		tasks:  make(map[string]bgTask),
	}
}
