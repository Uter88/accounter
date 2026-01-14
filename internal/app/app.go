package app

import (
	"accounter/config"
	"accounter/internal/domain/event"
	"accounter/internal/domain/task"
	"accounter/internal/domain/user"
	"accounter/internal/infrastructure/adapters/postgres"
	"accounter/internal/infrastructure/auth"
	comparator "accounter/internal/infrastructure/comparators"
	"accounter/internal/infrastructure/messaging/kafka"
	"accounter/internal/infrastructure/renderers"
	"accounter/internal/infrastructure/transport/rest"
	v1 "accounter/internal/infrastructure/transport/rest/endpoints/v1"
	"accounter/internal/infrastructure/transport/websocket"

	"accounter/pkg/logger"
	"context"
)

// Application context
type AppContext struct {
	// Application config
	config config.Config

	// HTTP server
	server *rest.Server

	// Default logger
	logger logger.Logger

	// Database client
	db *postgres.SQLClient

	// Message queue broker
	mq *kafka.KafkaBroker

	// Background tasks map
	tasks map[string]bgTask
}

// Background task
type bgTask interface {
	Run(ctx context.Context) error
	Name() string
}

// RegisterTask add background task to local store
func (a *AppContext) RegisterTask(task bgTask) *AppContext {
	a.tasks[task.Name()] = task

	return a
}

// Run application
func (a *AppContext) Run(ctx context.Context) {
	a.RegisterTask(a.mq)
	a.RegisterTask(a.server)

	a.launchTasks(ctx)

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
	return a.initConnections(ctx).initServices(ctx)
}

// Init connections with databases
func (a *AppContext) initConnections(ctx context.Context) *AppContext {
	if err := a.db.Connect(ctx); err != nil {
		a.logger.Fatalln(err)
	}

	return a
}

// Init application services
func (a *AppContext) initServices(ctx context.Context) *AppContext {

	// Init repositories
	userRepo := postgres.NewUserRepository(a.db)
	taskRepo := postgres.NewTaskRepository(a.db)
	eventRepo := postgres.NewEventRepository(a.db)

	// Init main services

	// Init Event service
	eventService := event.NewEventService(
		eventRepo,
		comparator.NewComparator().Fields(task.ComparationFields[:]...),
	)

	// Init authorization service
	authService := auth.NewAuthService(userRepo)

	// Init User service
	userService := user.NewUserService(userRepo)

	// Init Task service
	taskService := task.NewTaskService(
		taskRepo,
		renderers.NewTaskRenderer(),
		eventService,
	)

	// Init websocket service
	websocketService := websocket.NewWebsocketService(
		authService,
		a.config,
		a.logger.WithPerfix("WS"),
	)

	// Register Event publishers
	eventService.RegisterPublisher(a.mq)

	// Register Event subscribers
	a.mq.RegisterSubscribers(websocketService)

	// Init HTTP server engine
	params := v1.EngineParams{
		Config:           a.config,
		Logger:           a.logger.WithPerfix("HTTP"),
		AuthService:      authService,
		UserService:      userService,
		TaskService:      taskService,
		EventService:     eventService,
		WebsocketService: websocketService,
	}
	v1Engine := v1.NewEngine(params)

	// Init HTTP server
	ginServer := rest.NewGinServer(a.config, v1Engine)
	a.server = rest.NewHTTPServer(ctx, a.config, params.Logger, ginServer)

	return a
}

// Shutdown application
func (a *AppContext) Shutdown() {
	if err := a.db.Disconnect(); err != nil {
		a.logger.Errorf("Error disconnect from db: %s", err.Error())
	} else {
		a.logger.Info("Success disconnected from db")
	}

	if err := a.mq.Close(); err != nil {
		a.logger.Errorf("Error close mq reader: %s", err.Error())
	} else {
		a.logger.Info("Sucess close mq reader")
	}

	a.logger.Info("Shutdown system")
}

// NewAppContext creates new AppContext
func NewAppContext(ctx context.Context, cfg config.Config, logger logger.Logger) *AppContext {
	return &AppContext{
		config: cfg,
		logger: logger.WithPerfix("APP"),
		db:     postgres.NewSQLClient(cfg.DB),
		mq:     kafka.NewBroker(ctx, cfg, logger.WithPerfix("KAFKA")),
		tasks:  make(map[string]bgTask),
	}
}
