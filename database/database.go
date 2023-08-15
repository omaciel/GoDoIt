package database

import (
	"fmt"
	"log"
	"os"

	sql "github.com/omaciel/GoDoIt/domain/sqlite"
	"github.com/omaciel/GoDoIt/domain/task"
	"github.com/omaciel/GoDoIt/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConnector interface {
	connectDatabase() (func(), error)
}

type Database struct {
	db interface{}
	dc DatabaseConnector
}

func newDatabase(dc DatabaseConnector) *Database {
	return &Database{dc: dc}
}

func NewSqliteDatabase() {
	db := &SqliteDatabase{dns: "file::memory:?cache=shared"}
	db.connectDatabase()
}

type SqliteDatabase struct {
	dns string
}

func (sl SqliteDatabase) connectDatabase() (func(), error) {
	db, err := gorm.Open(sqlite.Open(sl.dns), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		return nil, err
	}

	// Return a function for cleanup after the test is done.
	cleanUp := func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}

	DB = DbInstance{
		Db: db,
	}

	return cleanUp, nil
}

type DbInstance struct {
	Db *gorm.DB
}

var DB DbInstance

var Repo task.TaskRepository

func InitDB() {
	Repo, _ = sql.NewSqliteDBRepository()
}

// MockSetupTestDB sets up an in-memory SQLite database for testing purposes.
func MockSetupTestDB() (func(), error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		return nil, err
	}

	// Return a function for cleanup after the test is done.
	cleanup := func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}

	DB = DbInstance{
		Db: db,
	}
	return cleanup, nil
}

func ConnectDb() {
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
	}

	log.Println("Connected to the database.")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running database migrations.")
	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatal("Failed to migrate the database schema. \n", err)
	}

	DB = DbInstance{
		Db: db,
	}
}
