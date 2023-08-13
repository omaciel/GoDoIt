package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/omaciel/GoDoIt/entity"
)

// MemoryRepository fulfills the TaskRepository interface
type MemoryRepository struct {
	Records map[uuid.UUID]entity.Task
	sync.Mutex
}

// NewMemoryRepository creates an in-memory datastore for Tasks
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		Records: make(map[uuid.UUID]entity.Task),
	}
}

// Get satifies the Get TaskRepository interface method
func (mr *MemoryRepository) Get(ctx context.Context, id uuid.UUID) (entity.Task, error) {
	if task, ok := mr.Records[id]; ok {
		return task, nil
	}
	return entity.Task{}, entity.ErrTaskNotFound
}

// Post satifies the Post TaskRepository interface method
func (mr *MemoryRepository) Post(ctx context.Context, task entity.Task) error {
	if mr.Records == nil {
		mr.Lock()
		mr.Records = make(map[uuid.UUID]entity.Task)
		mr.Unlock()
	}

	// Does the Task already exist?
	if _, ok := mr.Records[task.ID]; ok {
		return entity.ErrTaskUniqueConstraint
	}
	mr.Lock()
	mr.Records[task.ID] = task
	mr.Unlock()
	return nil
}

// Delete satisfies the Delete TaskRepository interface method
func (mr *MemoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if Task exists first.
	if _, ok := mr.Records[id]; !ok {
		return entity.ErrTaskNotFound
	}

	// Delete the Task.
	delete(mr.Records, id)
	// Assure that Task could not be found.
	if _, ok := mr.Records[id]; ok {
		return entity.ErrCouldNotDeleteTask
	}
	return nil
}

// All satisfies the All TaskRepository interface method
func (mr *MemoryRepository) All(ctx context.Context) ([]entity.Task, error) {
	if mr.Records == nil {
		mr.Lock()
		mr.Records = make(map[uuid.UUID]entity.Task)
		mr.Unlock()
	}

	values := make([]entity.Task, 0)

	for _, value := range mr.Records {
		values = append(values, value)
	}

	return values, nil
}

// Put satisfies the Put TaskRepository interface method method
func (mr *MemoryRepository) Put(ctx context.Context, task entity.Task) error {
	if _, ok := mr.Records[task.ID]; !ok {
		return entity.ErrTaskNotFound
	}

	mr.Records[task.ID] = task
	return nil
}
