package postgres

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/omaciel/GoDoIt/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgresRepository fulfills the TaskRepository interface
type PostgresRepository struct {
	Db *gorm.DB
}

// NewPostgresRepository creates a Postgres datastore
func NewPostgresRepository() (*PostgresRepository, error) {
	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=America/New_York",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Fatal("Failed to connect to the database. \n", err)
		return nil, err
	}

	log.Println("Connected to the database.")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running database migrations.")
	err = db.AutoMigrate(&entity.Task{})
	if err != nil {
		log.Fatal("Failed to migrate the database schema. \n", err)
		return nil, err
	}

	return &PostgresRepository{
		Db: db,
	}, nil
}

// Get satifies the Get TaskRepository interface method
func (pr *PostgresRepository) Get(ctx context.Context, id uuid.UUID) (entity.Task, error) {
	var task entity.Task

	result := pr.Db.Where("id = ?", id).First(&task)
	if result.Error != nil {
		return task, result.Error
	}
	return task, nil
}

// Post satifies the Post TaskRepository interface method
func (pr *PostgresRepository) Post(ctx context.Context, task *entity.Task) error {
	if task.ID == uuid.Nil {
		task.ID = uuid.New()
	}
	result := pr.Db.Create(&task)
	return result.Error
}

// Delete satisfies the Delete TaskRepository interface method
func (pr *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

// All satisfies the All TaskRepository interface
func (pr *PostgresRepository) All(ctx context.Context) ([]entity.Task, error) {
	var tasks []entity.Task = make([]entity.Task, 0)
	pr.Db.Find(&tasks)
	return tasks, nil
}

// Put satisfies the Put TaskRepository interface method
func (pr *PostgresRepository) Put(ctx context.Context, task *entity.Task) error {
	return nil
}