package sqlite

import (
	"context"

	"github.com/google/uuid"
	"github.com/omaciel/GoDoIt/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SqliteDBRepository fulfills the TaskRepository interface
type SqliteDBRepository struct {
	Db *gorm.DB
}

// New creates an in-SqliteDB datastore for Tasks
func New() (*SqliteDBRepository, error) {
	db, err := gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&entity.Task{})
	if err != nil {
		return nil, err
	}
	return &SqliteDBRepository{
		Db: db,
	}, nil
}

// Get satifies the Get TaskRepository interface method
func (repo *SqliteDBRepository) Get(ctx context.Context, id uuid.UUID) (entity.Task, error) {
	var task entity.Task

	result := repo.Db.Where("id = ?", id).First(&task)
	if result.Error != nil {
		return task, result.Error
	}
	return task, nil
}

// Post satifies the Post TaskRepository interface method
func (repo *SqliteDBRepository) Post(ctx context.Context, task entity.Task) error {
	result := repo.Db.Create(&task)
	return result.Error
}

// Delete satisfies the Delete TaskRepository interface method
func (repo *SqliteDBRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

// All satisfies the All TaskRepository interface
func (repo *SqliteDBRepository) All(ctx context.Context) ([]entity.Task, error) {
	return make([]entity.Task, 0), nil
}

// Put satisfies the Put TaskRepository interface method
func (repo *SqliteDBRepository) Put(ctx context.Context, task entity.Task) error {
	return nil
}
