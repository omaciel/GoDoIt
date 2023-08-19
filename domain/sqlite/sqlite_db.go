package sqlite

import (
	"context"
	"log"

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

// NewSqliteDBRepository creates an in-memory SqliteDB datastore for Tasks
func NewSqliteDBRepository() (*SqliteDBRepository, error) {
	db, err := gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	if err != nil {
		log.Fatal("Failed to connect to the database. \n", err)
		return nil, err
	}

	err = db.AutoMigrate(&entity.Task{})
	if err != nil {
		log.Fatal("Failed to migrate the database schema. \n", err)
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
func (repo *SqliteDBRepository) Post(ctx context.Context, task *entity.Task) error {
	if task.ID == uuid.Nil {
		task.ID = uuid.New()
	}
	result := repo.Db.Create(&task)
	return result.Error
}

// Delete satisfies the Delete TaskRepository interface method
func (repo *SqliteDBRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := repo.Db.Where("id = ?", id).Delete(entity.Task{})
	if result.Error != nil {
		return entity.ErrCouldNotDeleteTask
	}
	return nil
}

// All satisfies the All TaskRepository interface
func (repo *SqliteDBRepository) All(ctx context.Context) ([]entity.Task, error) {
	var tasks []entity.Task = make([]entity.Task, 0)
	repo.Db.Find(&tasks)
	return tasks, nil
}

// Put satisfies the Put TaskRepository interface method
func (repo *SqliteDBRepository) Put(ctx context.Context, task *entity.Task) error {
	if result := repo.Db.Save(&task); result.Error != nil {
		return result.Error
	}
	return nil
}
