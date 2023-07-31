package database

import (
	"fmt"
	"log"
	"os"

	"github.com/omaciel/GoDoIt/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var DB DbInstance

func ConnectDb() {
	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=America/New_York",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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