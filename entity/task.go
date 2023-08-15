package entity

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrInvalidPriorityLevel   = errors.New("invalid priority level")
	ErrInvalidTaskDescription = errors.New("the description cannot be empty")
	ErrTaskUniqueConstraint   = errors.New("unique constraint or index violation")
	ErrTaskNotFound           = errors.New("the task was not found in the repository")
	ErrCouldNotDeleteTask     = errors.New("could not delete the task")
)

// Priority represents how important a Task is for the user.
type Priority uint

func (p Priority) Validate() error {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return nil
	}
	return ErrInvalidPriorityLevel
}

const (
	// PriorityLow represents a non-urgent task.
	PriorityLow = Priority(iota + 1)

	// PriorityMedium represents an important task.
	PriorityMedium

	// PriorityHigh represents a very important task.
	PriorityHigh
)

type Task struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key;unique;type:uuid;column:id"`
	Description string    `json:"description" gorm:"text;not null;default:null"`
	Priority    Priority  `json:"priority" gorm:"default:3"`
	Completed   bool      `json:"completed" gorm:"default:false"`
}

// NewTask creates a new Task with sane default values
func NewTask(description string) *Task {
	return &Task{
		ID:          uuid.New(),
		Description: description,
		Priority:    PriorityLow,
		Completed:   false,
	}
}

// WithPriority returns a Task with the provided Priority set
func (t *Task) WithPriority(level Priority) *Task {
	t.Priority = level
	return t
}

// WithCompleted returns a Task with the provided Priority set
func (t *Task) WithCompleted(done bool) *Task {
	t.Completed = done
	return t
}

func (t *Task) Validate() error {
	if t.Description == "" {
		return ErrInvalidTaskDescription
	}

	if err := t.Priority.Validate(); err != nil {
		return fmt.Errorf("priority is invalid: %w", err)
	}

	return nil
}
