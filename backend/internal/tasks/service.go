package tasks

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, projectID, userID string, req CreateTaskRequest) (*TaskResponse, error)
	List(ctx context.Context, projectID, status, assignee string) ([]TaskResponse, error)
	Update(ctx context.Context, id, userID string, req UpdateTaskRequest) (*TaskResponse, error)
	Delete(ctx context.Context, id, userID string) error
}

type taskService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &taskService{repo: repo}
}

func (s *taskService) Create(ctx context.Context, projectID, userID string, req CreateTaskRequest) (*TaskResponse, error) {
	_, err := s.repo.GetProjectOwner(ctx, projectID)
	if err != nil {
		return nil, err
	}

	status := "todo"
	if req.Status != "" {
		status = req.Status
	}

	priority := "medium"
	if req.Priority != "" {
		priority = req.Priority
	}

	t := &Task{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Priority:    priority,
		ProjectID:   projectID,
		AssigneeID:  req.AssigneeID,
		DueDate:     req.DueDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}

	return s.mapToResponse(t), nil
}

func (s *taskService) List(ctx context.Context, projectID, status, assignee string) ([]TaskResponse, error) {
	tasks, err := s.repo.ListByProject(ctx, projectID, status, assignee)
	if err != nil {
		return nil, err
	}

	var res []TaskResponse
	for _, t := range tasks {
		res = append(res, *s.mapToResponse(&t))
	}
	
	if res == nil {
		res = make([]TaskResponse, 0)
	}

	return res, nil
}

func (s *taskService) Update(ctx context.Context, id, userID string, req UpdateTaskRequest) (*TaskResponse, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	ownerID, err := s.repo.GetProjectOwner(ctx, t.ProjectID)
	if err != nil {
		return nil, err
	}

	// Update allows project owner or assignee
	isAssignee := t.AssigneeID != nil && *t.AssigneeID == userID
	if ownerID != userID && !isAssignee {
		return nil, errors.New("forbidden")
	}

	if req.Title != nil {
		t.Title = *req.Title
	}
	if req.Description != nil {
		t.Description = *req.Description
	}
	if req.Status != nil {
		t.Status = *req.Status
	}
	if req.Priority != nil {
		t.Priority = *req.Priority
	}
	if req.AssigneeID != nil {
		t.AssigneeID = req.AssigneeID
	}
	if req.DueDate != nil {
		t.DueDate = req.DueDate
	}
	t.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}

	return s.mapToResponse(t), nil
}

func (s *taskService) Delete(ctx context.Context, id, userID string) error {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	ownerID, err := s.repo.GetProjectOwner(ctx, t.ProjectID)
	if err != nil {
		return err
	}

	// Only project owner can delete, assigning task deletion purely to owner
	if ownerID != userID {
		return errors.New("forbidden")
	}

	return s.repo.Delete(ctx, id)
}

func (s *taskService) mapToResponse(t *Task) *TaskResponse {
	return &TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Priority:    t.Priority,
		ProjectID:   t.ProjectID,
		AssigneeID:  t.AssigneeID,
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
