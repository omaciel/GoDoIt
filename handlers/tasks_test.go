package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/omaciel/GoDoIt/database"
	"github.com/omaciel/GoDoIt/domain/memory"
	"github.com/omaciel/GoDoIt/entity"
	"github.com/omaciel/GoDoIt/router"
	"github.com/stretchr/testify/assert"
)

const (
	API_PATH_WITH_ID          string = "/task/%s"
	GENERIC_TASK_NAME         string = "Test Task"
	HEADER_APPLICATION_FORMAT string = "application/json"
	HEADER_CONTENT_TYPE       string = "Content-Type"
	NO_ERROR_EXPECTED         string = "did not expect to receive an error"
)

// Tests for GetTask method
func TestGetTaskInvalidUUID(t *testing.T) {
	database.Repo = memory.NewMemoryRepository()
	taskUuid := "aaa"

	app := fiber.New()
	router.SetupTaskRoutes(app)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf(API_PATH_WITH_ID, taskUuid), nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "Should return HTTP 400 code")
}

func TestGetTaskContentDoesNotExist(t *testing.T) {
	database.Repo = memory.NewMemoryRepository()
	taskUuid := uuid.New()

	defer func() {
		err := database.Repo.Delete(context.Background(), taskUuid)
		if err != nil {
			return
		}
	}()

	app := fiber.New()
	router.SetupTaskRoutes(app)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf(API_PATH_WITH_ID, taskUuid), nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode, "Should return HTTP 204 code")
}

func TestGetTaskValidContent(t *testing.T) {
	database.Repo = memory.NewMemoryRepository()

	task := entity.NewTask(GENERIC_TASK_NAME)
	err := database.Repo.Post(context.Background(), task)
	assert.NoError(t, err, NO_ERROR_EXPECTED)

	app := fiber.New()
	router.SetupTaskRoutes(app)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/task/%s", task.ID), nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode, "Should return HTTP 200 code")
	assert.Equal(t, task.Description, "Test Task", "Should return correct description")
}

// Tests for ListTasks method
func TestListTasks(t *testing.T) {
	database.Repo = memory.NewMemoryRepository()

	task1 := entity.NewTask("Test Task 1").WithPriority(entity.PriorityHigh)
	task2 := entity.NewTask("Test Task 2").WithPriority(entity.PriorityMedium)
	err := database.Repo.Post(context.Background(), task1)
	assert.NoError(t, err, NO_ERROR_EXPECTED)
	err = database.Repo.Post(context.Background(), task2)
	assert.NoError(t, err, NO_ERROR_EXPECTED)

	defer func(uuids ...uuid.UUID) {
		for _, id := range uuids {
			err := database.Repo.Delete(context.Background(), id)
			if err != nil {
				return
			}
		}
	}(task1.ID, task2.ID)

	app := fiber.New()
	router.SetupTaskRoutes(app)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode, "Should return HTTP 200 code")

	var tasks []entity.Task
	err = json.NewDecoder(resp.Body).Decode(&tasks)
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

// Tests for PostTask method
func TestPostTaskInvalidJSON(t *testing.T) {
	database.Repo = memory.NewMemoryRepository()

	app := fiber.New()
	router.SetupTaskRoutes(app)

	req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set(HEADER_CONTENT_TYPE, HEADER_APPLICATION_FORMAT)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Check the response status code and body
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var errorMessage map[string]string
	err = json.NewDecoder(resp.Body).Decode(&errorMessage)
	assert.NoError(t, err)
	assert.Contains(t, errorMessage, "message")
}

func TestPostTaskInvalidTask(t *testing.T) {
	database.Repo = memory.NewMemoryRepository()

	// Add a new Task
	task := entity.NewTask(GENERIC_TASK_NAME)
	err := database.Repo.Post(context.Background(), task)
	assert.NoError(t, err, NO_ERROR_EXPECTED)

	taskJSON, _ := json.Marshal(task)

	app := fiber.New()
	router.SetupTaskRoutes(app)

	// Try to create the same Task with same ID
	req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewBuffer(taskJSON))
	req.Header.Set(HEADER_CONTENT_TYPE, HEADER_APPLICATION_FORMAT)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Check the response status code and body
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func TestPostTaskSuccess(t *testing.T) {
	database.Repo = memory.NewMemoryRepository()

	task := entity.NewTask("New Task").WithCompleted(false)
	defer func() {
		err := database.Repo.Delete(context.Background(), task.ID)
		if err != nil {
			return
		}
	}()

	app := fiber.New()
	router.SetupTaskRoutes(app)

	// Test case 1: Valid request body
	// task := models.Task{Description: "New Task", Completed: false}
	taskJSON, _ := json.Marshal(task)
	req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewBuffer(taskJSON))
	req.Header.Set(HEADER_CONTENT_TYPE, HEADER_APPLICATION_FORMAT)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Check the response status code and body
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var createdTask entity.Task
	err = json.NewDecoder(resp.Body).Decode(&createdTask)
	assert.NoError(t, err)
	assert.Equal(t, task.Description, createdTask.Description)
	assert.False(t, createdTask.Completed)
}

// Tests for UpdateTask method
func TestUpdateTask(t *testing.T) {
	database.Repo = memory.NewMemoryRepository()

	task := entity.NewTask(GENERIC_TASK_NAME).WithPriority(entity.PriorityMedium)
	err := database.Repo.Post(context.Background(), task)
	assert.ErrorIs(t, err, nil, NO_ERROR_EXPECTED)
	assert.Equal(t, task.Priority, entity.PriorityMedium, "expected Priority to be PriorityMedium")
	assert.False(t, task.Completed, "expected Completed to be false")
	defer func() {
		err := database.Repo.Delete(context.Background(), task.ID)
		if err != nil {
			return
		}
	}()

	app := fiber.New()
	router.SetupTaskRoutes(app)

	// Change Priority to PriorityHigh and Completed to true
	task.Priority = entity.PriorityHigh
	task.Completed = true

	updatedTaskJSON, _ := json.Marshal(task)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf(API_PATH_WITH_ID, task.ID), bytes.NewBuffer(updatedTaskJSON))
	req.Header.Set(HEADER_CONTENT_TYPE, HEADER_APPLICATION_FORMAT)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Check the response status code and body
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var updated entity.Task
	err = json.NewDecoder(resp.Body).Decode(&updated)
	assert.NoError(t, err)
	assert.Equal(t, task.Description, updated.Description)
	assert.Equal(t, entity.PriorityHigh, updated.Priority)
	assert.True(t, updated.Completed)
}

// Tests for DeleteTask method
func TestDeleteTaskInvalidUUID(t *testing.T) {
	database.Repo = memory.NewMemoryRepository()
	taskUuid := "aaa"

	app := fiber.New()
	router.SetupTaskRoutes(app)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf(API_PATH_WITH_ID, taskUuid), nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Check the response status code and body
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestDeleteTaskContentDoesNotExist(t *testing.T) {
	database.Repo = memory.NewMemoryRepository()
	taskUuid := uuid.New()

	app := fiber.New()
	router.SetupTaskRoutes(app)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf(API_PATH_WITH_ID, taskUuid), nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Check the response status code and body
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}
func TestDeleteTaskSuccess(t *testing.T) {
	database.Repo = memory.NewMemoryRepository()

	task := entity.NewTask(GENERIC_TASK_NAME)
	err := database.Repo.Post(context.Background(), task)
	assert.ErrorIs(t, err, nil, NO_ERROR_EXPECTED)
	defer func() {
		err := database.Repo.Delete(context.Background(), task.ID)
		if err != nil {
			return
		}
	}()

	app := fiber.New()
	router.SetupTaskRoutes(app)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf(API_PATH_WITH_ID, task.ID), nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Check the response status code and body
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
