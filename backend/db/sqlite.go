package db

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// InitMySQLConnection connect to MySQL database
func InitConnection(ctx context.Context, addr string) *sqlx.DB {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	if conn, err := sqlx.ConnectContext(ctx, "sqlite3", "database.sql"); err != nil {
		log.Fatalln("Could not create MySQL map DB pool: ", err.Error())

	} else if err := conn.PingContext(ctx); err != nil {
		log.Fatalln("Error ping MySQL map DB: ", err.Error())

	} else {
		log.Println("Connecting to MySQL map DB is success")

		if err = migrate(conn); err != nil {
			log.Fatalln("Error migrate: ", err.Error())
		} else {
			log.Println("Success migration of tables")
		}

		return conn
	}

	return nil
}

func migrate(conn *sqlx.DB) error {
	schema := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		login TEXT,
		password TEXT,
		name TEXT,
		surname TEXT,
		patronymic TEXT,
		price_per_hour REAL
	);`

	// execute a query on the server
	if _, err := conn.Exec(schema); err != nil {
		return err
	}

	schema = `
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			task_id TEXT NULL,
			task_status TEXT,
			description TEXT,
			work_begin INTEGER,
			work_end INTEGER,
			month INTEGER
		)
	`

	if _, err := conn.Exec(schema); err != nil {
		return err
	}

	return nil
}
