package main

import (
	"accounter/adapters/adapter_sql"
	"accounter/config"
)

const SQLiteUserSchema = `CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY,
	login TEXT,
	password TEXT,
	name TEXT,
	surname TEXT,
	patronymic TEXT,
	price_per_hour REAL
);`

const SQLiteTaskSchema = `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY,
		user_id INTEGER,
		task_id TEXT NULL,
		task_status TEXT,
		description TEXT,
		work_begin INTEGER,
		work_end INTEGER,
		date INTEGER
	)
`

const PostgresUserSchema = `CREATE TABLE IF NOT EXISTS public.users (
    id BIGSERIAL NOT NULL,
    login CHARACTER VARYING(50) NOT NULL,
    password CHARACTER VARYING(50) NOT NULL,
    name CHARACTER VARYING(50) NOT NULL,
    surname CHARACTER VARYING(50) NOT NULL,
    patronymic CHARACTER VARYING(50) NOT NULL,
    price_per_hour NUMERIC(12, 2) NOT NULL,
    PRIMARY KEY (id)
);`

const PostgresTaskSchema = `
	CREATE TABLE IF NOT EXISTS public.tasks (
    	id BIGSERIAL NOT NULL,
		user_id BIGINT NOT NULL,
    	task_id CHARACTER VARYING(30) NOT NULL,
    	task_status CHARACTER VARYING(30) NOT NULL,
    	description TEXT NOT NULL,
    	work_begin BIGINT NOT NULL,
    	work_end BIGINT NOT NULL,
		date BIGINT NOT NULL,
    PRIMARY KEY (id)
	)
`

func main() {
	ctx, cancel := config.InitGracefulShutdownCtx()
	cfg := config.InitConfig()
	logger := config.NewLogger(cfg.DebugMode, cfg.AppMode, "logs")
	client := adapter_sql.NewSQLClient(cfg.DB.Driver, cfg.DB.DSN)

	if err := client.Connect(ctx); err != nil {
		logger.Fatalln("Fail to migrate:", err.Error())
	}

	defer func() {
		cancel()
		client.Disconnect()
	}()

	var (
		usersScheme, tasksScheme string
	)

	switch cfg.DB.Driver {
	case "sqlite3":
		usersScheme = SQLiteUserSchema
		tasksScheme = SQLiteTaskSchema

	case "postgres":
		usersScheme = PostgresUserSchema
		tasksScheme = PostgresTaskSchema
	}

	// execute a query on the server
	if _, err := client.DB().ExecContext(ctx, usersScheme); err != nil {
		logger.Fatalf("error migrate users table: %s", err.Error())
	}

	if _, err := client.DB().ExecContext(ctx, tasksScheme); err != nil {
		logger.Fatalf("error migrate tasks table: %s", err.Error())
	}

	logger.Info("Migrates is success!")
}
