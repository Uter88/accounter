package main

import (
	"accounter/config"
	"accounter/internal/infrastructure/adapter_sql"
	"accounter/pkg/logger"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "embed"
)

var tables = []string{"users", "tasks", "events"}

func main() {
	path := flag.String("path", "/migrations", "Path with scripts")
	flag.Parse()

	ctx, cancel := config.InitGracefulShutdownCtx()
	cfg := config.InitConfig()
	logger := logger.NewLogger(cfg.DebugMode)
	client := adapter_sql.NewSQLClient(cfg.DB)

	if err := client.Connect(ctx); err != nil {
		logger.Fatalln("Fail to migrate:", err.Error())
	}

	defer func() {
		cancel()
		client.Disconnect()
	}()

	schemes, err := loadSchemes(*path, cfg.DB.Driver)

	if err != nil {
		logger.Fatalln("Fail to load tables schemes:", err.Error())
	}

	err = client.BeginTx(ctx, func(ctx context.Context) error {

		for _, scheme := range schemes {
			if _, err := client.GetExecutor(ctx).ExecContext(ctx, scheme); err != nil {
				return fmt.Errorf("error migrate %s table: %s", scheme, err.Error())
			}
		}

		return nil
	})

	if err != nil {
		logger.Errorf("Migration error: %s", err.Error())
	} else {
		logger.Info("Migrates is success!")
	}
}

func loadSchemes(folder, driver string) (result []string, err error) {
	for _, table := range tables {
		path, _ := os.Getwd()
		path = strings.Replace(path, "/tools/migrator", folder, 1)
		path = fmt.Sprintf("%s/%s/%s.sql", path, driver, table)

		contents, err := os.ReadFile(path)

		if err != nil {
			return nil, err
		}

		result = append(result, string(contents))
	}

	return
}
