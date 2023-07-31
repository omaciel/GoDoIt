package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/omaciel/GoDoIt/handlers"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.ListTasks)

	app.Post("/task", handlers.CreateTask)
}
