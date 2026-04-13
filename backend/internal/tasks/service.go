package tasks

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"taskflow/internal/projects"
	"taskflow/internal/users"
)

var ErrForbidden = errors.New("forbidden")
var ErrPastDueDate = errors.New("due date cannot be in the past")
var ErrInvalidAssignee = errors.New("assignee user not found")

type Service interface {
	Create(ctx context.Context, projectID, userID string, req CreateTaskRequest) (*TaskResponse, error)
	List(ctx context.Context, projectID, status, assignee string, page, limit int) (*ListTasksResponse, error)
	Update(ctx context.Context, id, userID string, req UpdateTaskRequest) (*TaskResponse, error)
	Delete(ctx context.Context, id, userID string) error
}

type taskService struct {
	repo        Repository
	projectRepo projects.Repository
	userRepo    users.Repository
}

func NewService(repo Repository, projectRepo projects.Repository, userRepo users.Repository) Service {
	return &taskService{repo: repo, projectRepo: projectRepo, userRepo: userRepo}
}

func (s *taskService) Create(ctx context.Context, projectID, userID string, req CreateTaskRequest) (*TaskResponse, error) {
	// Verify project exists and get owner
	ownerID, err := s.repo.GetProjectOwner(ctx, projectID)
	if err != nil {
		return nil, err // ErrProjectNotFound propagated
	}

	// Only the project owner can create tasks
	if ownerID != userID {
		return nil, ErrForbidden
	}

	// Validate due date is not in the past (compare by date only, ignore time)
	if req.DueDate != nil {
		today := time.Now().UTC().Truncate(24 * time.Hour)
		due := req.DueDate.UTC().Truncate(24 * time.Hour)
		if due.Before(today) {
			return nil, ErrPastDueDate
		}
	}

	// Validate assignee exists if provided
	if req.AssigneeID != nil {
		exists, err := s.userRepo.Exists(ctx, *req.AssigneeID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrInvalidAssignee
		}
	}

	status := StatusTodo
	if req.Status != "" {
		status = req.Status
	}

	priority := PriorityMedium
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
		CreatedBy:   userID,
		DueDate:     req.DueDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}

	return s.mapToResponse(t), nil
}

func (s *taskService) List(ctx context.Context, projectID, status, assignee string, page, limit int) (*ListTasksResponse, error) {
	// Verify project exists before listing tasks
	exists, err := s.projectRepo.Exists(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrProjectNotFound
	}

	tasks, total, err := s.repo.ListByProject(ctx, projectID, status, assignee, page, limit)
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

	return &ListTasksResponse{
		Tasks: res,
		Pagination: PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
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

	// Allow: project owner, task creator, or current assignee
	isAssignee := t.AssigneeID != nil && *t.AssigneeID == userID
	isCreator := t.CreatedBy == userID
	if ownerID != userID && !isAssignee && !isCreator {
		return nil, ErrForbidden
	}

	// Validate new due date is not in the past
	if req.DueDate != nil {
		today := time.Now().UTC().Truncate(24 * time.Hour)
		due := req.DueDate.UTC().Truncate(24 * time.Hour)
		if due.Before(today) {
			return nil, ErrPastDueDate
		}
	}

	// Validate new assignee exists
	if req.AssigneeID != nil {
		exists, err := s.userRepo.Exists(ctx, *req.AssigneeID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrInvalidAssignee
		}
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

	// Allow project owner or task creator to delete
	isCreator := t.CreatedBy == userID
	if ownerID != userID && !isCreator {
		return ErrForbidden
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
		CreatedBy:   t.CreatedBy,
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
