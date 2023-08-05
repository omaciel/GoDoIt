package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/omaciel/GoDoIt/handlers"
)

func SetupTaskRoutes(app *fiber.App) {
	app.Get("/", handlers.ListTasks)

	app.Post("/task", handlers.CreateTask)
	app.Get("/task/:id", handlers.ShowTask)
	app.Patch("/task/:id", handlers.UpdateTask)
	app.Delete("/task/:id", handlers.DeleteTask)
}

func SetupRoutes(app *fiber.App) {
	SetupTaskRoutes(app)
}
