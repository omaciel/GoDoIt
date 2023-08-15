package sqlite_test

import (
	"context"
	"testing"

	"github.com/omaciel/GoDoIt/domain/sqlite"
	"github.com/omaciel/GoDoIt/entity"
)

func TestSqliteDbRepository_Get(t *testing.T) {
	repo, err := sqlite.NewSqliteDBRepository()
	if err != nil {
		t.Fatalf("failed to start Sqlite database: %v", err)
	}

	defer func() {
		sqlDB, _ := repo.Db.DB()
		sqlDB.Close()
	}()

	task := entity.NewTask("Test Task")
	result := repo.Db.Create(&task)
	if result.Error != nil {
		t.Fatalf("failed to create a task in the Sqlite database: %v", result.Error)
	}

	record, err := repo.Get(context.Background(), task.ID)
	if err != nil {
		t.Fatalf("could not find record matching ID %s: %v", task.ID, err)
	}

	if record.Description != task.Description {
		t.Fatalf("expected description to be %s, got %s", task.Description, record.Description)
	}

	if record.Priority != task.Priority {
		t.Fatalf("expected priority to be %v, got %v", task.Priority, record.Priority)
	}

	if record.Completed != task.Completed {
		t.Fatalf("expected completed to be %v, got %v", task.Completed, record.Completed)
	}
}
