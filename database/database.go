package database

import (
	"os"

	postgres "github.com/omaciel/GoDoIt/domain/postgres"
	sql "github.com/omaciel/GoDoIt/domain/sqlite"
	"github.com/omaciel/GoDoIt/domain/task"
)

var Repo task.TaskRepository

func InitDB() {
	dataLayer := os.Getenv("DATABASE")

	switch dataLayer {
	case "postgres":
		Repo, _ = postgres.NewPostgresRepository()
	default:
		Repo, _ = sql.NewSqliteDBRepository()
	}
}
