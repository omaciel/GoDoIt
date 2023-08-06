package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/omaciel/GoDoIt/database"
	"github.com/omaciel/GoDoIt/models"
	"github.com/omaciel/GoDoIt/router"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestShowTaskNoContent(t *testing.T) {
	cleanUp, _ := database.MockSetupTestDB()
	defer cleanUp()

	app := fiber.New()
	router.SetupTaskRoutes(app)

	req := httptest.NewRequest(http.MethodGet, "/task/-1", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode, "Should return HTTP 204 code")
}

func TestShowTaskContent(t *testing.T) {
	cleanUp, _ := database.MockSetupTestDB()
	defer cleanUp()

	task := models.Task{Description: "Test Task"}
	database.DB.Db.Create(&task)

	app := fiber.New()
	router.SetupTaskRoutes(app)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/task/%d", task.ID), nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode, "Should return HTTP 200 code")
	assert.Equal(t, task.Description, "Test Task", "Should return correct description")
}

func TestListTasks(t *testing.T) {
	cleanUp, _ := database.MockSetupTestDB()
	defer cleanUp()

	task1 := models.Task{Description: "Test Task 1", Priority: 1}
	task2 := models.Task{Description: "Test Task 2", Priority: 2}
	database.DB.Db.Create(&task1)
	database.DB.Db.Create(&task2)

	app := fiber.New()
	router.SetupTaskRoutes(app)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode, "Should return HTTP 200 code")

	var tasks []models.Task
	err := json.NewDecoder(resp.Body).Decode(&tasks)
	assert.NoError(t, err)

	assert.Len(t, tasks, 2)
	assert.Equal(t, task1.ID, tasks[0].ID)
	assert.Equal(t, task1.Description, tasks[0].Description)
	assert.Equal(t, task1.Priority, tasks[0].Priority)
	assert.Equal(t, task1.Completed, tasks[0].Completed)

	assert.Equal(t, task2.ID, tasks[1].ID)
	assert.Equal(t, task2.Description, tasks[1].Description)
	assert.Equal(t, task2.Priority, tasks[1].Priority)
	assert.Equal(t, task2.Completed, tasks[1].Completed)
}

func TestCreateTask(t *testing.T) {
	// Mock the database and create a Fiber context for testing
	cleanUp, _ := database.MockSetupTestDB()
	defer cleanUp()

	app := fiber.New()
	router.SetupTaskRoutes(app)

	// Test case 1: Valid request body
	task := models.Task{Description: "New Task", Completed: false}
	taskJSON, _ := json.Marshal(task)
	req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Check the response status code and body
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var createdTask models.Task
	err = json.NewDecoder(resp.Body).Decode(&createdTask)
	assert.NoError(t, err)
	assert.Equal(t, task.Description, createdTask.Description)
	assert.False(t, createdTask.Completed)
}

func TestCreateTaskInternalServerError(t *testing.T) {
	// Mock the database and create a Fiber context for testing
	cleanUp, _ := database.MockSetupTestDB()
	defer cleanUp()

	app := fiber.New()
	router.SetupTaskRoutes(app)

	req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Check the response status code and body
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var errorMessage map[string]string
	err = json.NewDecoder(resp.Body).Decode(&errorMessage)
	assert.NoError(t, err)
	assert.Contains(t, errorMessage, "message")
}

func TestUpdateTask(t *testing.T) {
	// Mock the database and create a Fiber context for testing
	cleanUp, _ := database.MockSetupTestDB()
	defer cleanUp()

	app := fiber.New()
	router.SetupTaskRoutes(app)

	task := models.Task{Description: "Test Task", Completed: false}
	database.DB.Db.Create(&task)

	// Test case 1: Valid request body
	updatedTask := models.Task{Description: "Updated Task", Completed: true}
	updatedTaskJSON, _ := json.Marshal(updatedTask)
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/task/%d", task.ID), bytes.NewBuffer(updatedTaskJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Check the response status code and body
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var updated models.Task
	err = json.NewDecoder(resp.Body).Decode(&updated)
	assert.NoError(t, err)
	assert.Equal(t, updatedTask.Description, updated.Description)
	assert.True(t, updated.Completed)

	// Test case 2: Task not found
	req = httptest.NewRequest(http.MethodPatch, "/task/999", bytes.NewBuffer(updatedTaskJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, -1)
	assert.NoError(t, err)

	// Check the response status code and body
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

func TestDeleteTaskCompleted(t *testing.T) {
	// Mock the database and create a Fiber context for testing
	cleanUp, _ := database.MockSetupTestDB()
	defer cleanUp()

	app := fiber.New()
	router.SetupTaskRoutes(app)

	task := models.Task{Description: "Test Task", Completed: false}
	database.DB.Db.Create(&task)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/task/%d", task.ID), nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Check the response status code
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Check that the task has been deleted
	var deletedTask models.Task
	err = database.DB.Db.Model(models.Task{}).Where("id = ?", task.ID).First(&deletedTask).Error
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestDeleteTask(t *testing.T) {
	// Mock the database and create a Fiber context for testing
	cleanUp, _ := database.MockSetupTestDB()
	defer cleanUp()

	app := fiber.New()
	router.SetupTaskRoutes(app)

	req := httptest.NewRequest(http.MethodDelete, "/task/999", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Check the response status code and body
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}
