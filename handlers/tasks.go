package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/omaciel/GoDoIt/database"
	"github.com/omaciel/GoDoIt/models"
)

func ListTasks(c *fiber.Ctx) error {
	var tasks []models.Task
	database.DB.Db.Find(&tasks)

	return c.Status(200).JSON(tasks)
}

func CreateTask(c *fiber.Ctx) error {
	task := new(models.Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	database.DB.Db.Create(&task)

	return c.Status(200).JSON(task)

}