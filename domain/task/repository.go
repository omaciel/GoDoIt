package task

import (
	"context"

	"github.com/google/uuid"
	"github.com/omaciel/GoDoIt/entity"
)

type Repository struct {
}

func NewTaskRepository() *Repository {
	return &Repository{}
}

type TaskRepository interface {
	Post(ctx context.Context, task *entity.Task) error
	Get(ctx context.Context, id uuid.UUID) (entity.Task, error)
	Put(ctx context.Context, task *entity.Task) error
	All(ctx context.Context) ([]entity.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
