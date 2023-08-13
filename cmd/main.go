package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/omaciel/GoDoIt/database"
	"github.com/omaciel/GoDoIt/router"
)

func main() {
	// Create Sqlite Database
	database.InitDB()

	// database.ConnectDb()

	app := fiber.New()

	router.SetupRoutes(app)

	err := app.Listen(":3000")
	if err != nil {
		fmt.Println("Could not start the server.", err)
		return
	}
}
