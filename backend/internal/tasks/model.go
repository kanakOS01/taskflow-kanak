package tasks

import (
	"fmt"
	"time"
)

// Status enum

type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

func (s *Status) UnmarshalJSON(data []byte) error {
	// strip surrounding quotes
	val := Status(data[1 : len(data)-1])
	switch val {
	case StatusTodo, StatusInProgress, StatusDone:
		*s = val
		return nil
	}
	return fmt.Errorf("invalid status %q: must be one of todo, in_progress, done", val)
}

// Priority enum

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

func (p *Priority) UnmarshalJSON(data []byte) error {
	val := Priority(data[1 : len(data)-1])
	switch val {
	case PriorityLow, PriorityMedium, PriorityHigh:
		*p = val
		return nil
	}
	return fmt.Errorf("invalid priority %q: must be one of low, medium, high", val)
}

// Task domain model

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	Priority    Priority   `json:"priority"`
	ProjectID   string     `json:"project_id"`
	AssigneeID  *string    `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

