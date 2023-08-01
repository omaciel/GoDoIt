package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/omaciel/GoDoIt/database"
	"github.com/omaciel/GoDoIt/models"
)

func ListTasks(c *fiber.Ctx) error {
	var tasks []models.Task
	database.DB.Db.Find(&tasks)

	return c.Status(fiber.StatusOK).JSON(tasks)
}

func CreateTask(c *fiber.Ctx) error {
	task := new(models.Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	result := database.DB.Db.Create(&task)
	if result.Error != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

func ShowTask(c *fiber.Ctx) error {
	task := models.Task{}
	id := c.Params("id")

	result := database.DB.Db.Where("id = ?", id).First(&task)
	if result.Error != nil || result.RowsAffected < 1 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": result.Error.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(task)
}

func UpdateTask(c *fiber.Ctx) error {
	task := new(models.Task)
	id := c.Params("id")

	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"message": err})
	}

	result := database.DB.Db.Model(&task).Where("id = ?", id).Updates(task)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

func DeleteTask(c *fiber.Ctx) error {
	task := models.Task{}
	id := c.Params("id")

	result := database.DB.Db.Where("id = ?", id).Delete(&task)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": result.Error.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}
