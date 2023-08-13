package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/omaciel/GoDoIt/handlers"
)

func SetupTaskRoutes(app *fiber.App) {
	app.Get("/", handlers.AllTasks)

	app.Post("/task", handlers.PostTask)
	app.Get("/task/:uuid", handlers.GetTask)
	app.Patch("/task/:uuid", handlers.PutTask)
	app.Delete("/task/:uuid", handlers.DeleteTask)
}

func SetupRoutes(app *fiber.App) {
	SetupTaskRoutes(app)
}
