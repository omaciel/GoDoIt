package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/omaciel/GoDoIt/database"
	"github.com/omaciel/GoDoIt/entity"
)

func AllTasks(c *fiber.Ctx) error {
	var tasks []entity.Task
	tasks, _ = database.Repo.All(context.Background())

	return c.Status(fiber.StatusOK).JSON(tasks)
}

func PostTask(c *fiber.Ctx) error {
	task := new(entity.Task)

	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	err := database.Repo.Post(context.Background(), *task)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": err})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

func GetTask(c *fiber.Ctx) error {
	uuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err})
	}

	task, err := database.Repo.Get(context.Background(), uuid)
	if err != nil {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": err})
	}

	return c.Status(fiber.StatusOK).JSON(task)
}

func PutTask(c *fiber.Ctx) error {
	task := new(entity.Task)

	uuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err})
	}

	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"message": err})
	}

	if _, err = database.Repo.Get(context.Background(), uuid); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"message": err})
	}

	if err = database.Repo.Put(context.Background(), *task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

func DeleteTask(c *fiber.Ctx) error {
	uuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err})
	}

	// Check that Task exists first.
	_, err = database.Repo.Get(context.Background(), uuid)
	if err != nil {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": err})
	}

	// Check that Task exists first.
	_, err = database.Repo.Get(context.Background(), uuid)
	if err != nil {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": err})
	}

	// Delete the Task.
	err = database.Repo.Delete(context.Background(), uuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err})
	}

	return c.SendStatus(fiber.StatusOK)
}
