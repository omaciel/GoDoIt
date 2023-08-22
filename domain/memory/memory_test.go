package memory_test

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/omaciel/GoDoIt/domain/memory"
	"github.com/omaciel/GoDoIt/entity"
	"github.com/stretchr/testify/assert"
)

func TestMemoryRepositoryGet(t *testing.T) {
	task := entity.NewTask("Write a unittest for MemoryRepository").
		WithPriority(entity.PriorityHigh).
		WithCompleted(true)

	id := task.ID

	repo := memory.MemoryRepository{
		Records: map[uuid.UUID]entity.Task{id: *task},
	}

	tests := []struct {
		name        string
		id          uuid.UUID
		expectedErr error
	}{
		{
			name:        "No task by ID",
			id:          uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"),
			expectedErr: entity.ErrTaskNotFound,
		},
		{
			name:        "Found task by ID",
			id:          id,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := repo.Get(context.Background(), tt.id)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *memory.MemoryRepository
	}{
		{
			name: "Create a new MemoryRepository with empty Records.",
			want: &memory.MemoryRepository{
				Records: make(map[uuid.UUID]entity.Task),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := memory.NewMemoryRepository()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
			assert.Equal(t, len(got.Records), 0)
		})
	}
}

func TestMemoryRepositoryPost(t *testing.T) {
	task0 := entity.NewTask("task 0")

	tests := []struct {
		name    string
		task    entity.Task
		wantErr error
	}{
		{
			"Add a new task",
			*entity.NewTask("task 1"),
			nil,
		},
		{
			"Add existing task",
			*task0,
			entity.ErrTaskUniqueConstraint,
		},
	}
	mr := memory.MemoryRepository{
		Records: map[uuid.UUID]entity.Task{
			task0.ID: *task0,
		},
		Mutex: sync.Mutex{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mr.Post(context.Background(), &tt.task)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestMemoryRepositoryDelete(t *testing.T) {
	task0 := entity.NewTask("task 0")
	tests := []struct {
		name    string
		uuid    uuid.UUID
		wantErr error
	}{
		{"Task UUID doesn't exist", uuid.New(), entity.ErrTaskNotFound},
		{"Task is deleted", task0.ID, nil},
	}
	mr := memory.MemoryRepository{
		Records: map[uuid.UUID]entity.Task{
			task0.ID: *task0,
		},
		Mutex: sync.Mutex{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mr.Delete(context.Background(), tt.uuid)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestMemoryRepositoryAll(t *testing.T) {
	task0 := entity.NewTask("task 0")
	task1 := entity.NewTask("task 1")
	task2 := entity.NewTask("task 2")
	tests := []struct {
		name         string
		mr           *memory.MemoryRepository
		expectedSize int
		wantErr      error
	}{
		{
			"Return list of 3 tasks",
			&memory.MemoryRepository{
				Records: map[uuid.UUID]entity.Task{
					task0.ID: *task0,
					task1.ID: *task1,
					task2.ID: *task2,
				},
				Mutex: sync.Mutex{},
			},
			3,
			nil,
		},
		{
			"Return list with zero tasks",
			memory.NewMemoryRepository(),
			0,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := tt.mr.All(context.Background())
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.expectedSize, len(items))
		})
	}
}

func TestMemoryRepositoryPut(t *testing.T) {
	task0 := entity.NewTask("task 0")

	tests := []struct {
		name    string
		task    entity.Task
		wantErr error
	}{
		{
			"Cannot update Not Found task",
			*entity.NewTask("task 1"),
			entity.ErrTaskNotFound,
		},
		{
			"Update existing task description",
			*task0,
			nil,
		},
	}
	mr := memory.MemoryRepository{
		Records: map[uuid.UUID]entity.Task{
			task0.ID: *task0,
		},
		Mutex: sync.Mutex{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mr.Put(context.Background(), &tt.task)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
