package projects

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID string, req CreateProjectRequest) (*ProjectResponse, error)
	List(ctx context.Context, userID string, page, limit int) (*ListProjectsResponse, error)
	GetDetails(ctx context.Context, id string) (*ProjectDetailsResponse, error)
	GetStats(ctx context.Context, id string) (*ProjectStatsResponse, error)
	Update(ctx context.Context, id, userID string, req UpdateProjectRequest) (*ProjectResponse, error)
	Delete(ctx context.Context, id, userID string) error
}

type projectService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &projectService{repo: repo}
}

func (s *projectService) Create(ctx context.Context, userID string, req CreateProjectRequest) (*ProjectResponse, error) {
	p := &Project{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	return s.mapToResponse(p), nil
}

func (s *projectService) List(ctx context.Context, userID string, page, limit int) (*ListProjectsResponse, error) {
	projects, total, err := s.repo.List(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	res := make([]ProjectResponse, 0, len(projects))
	for _, p := range projects {
		res = append(res, *s.mapToResponse(&p))
	}

	return &ListProjectsResponse{
		Projects: res,
		Pagination: PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

func (s *projectService) GetDetails(ctx context.Context, id string) (*ProjectDetailsResponse, error) {
	p, tasks, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = make([]Task, 0)
	}

	return &ProjectDetailsResponse{
		ProjectResponse: *s.mapToResponse(p),
		Tasks:           tasks,
	}, nil
}

func (s *projectService) GetStats(ctx context.Context, id string) (*ProjectStatsResponse, error) {
	statusCounts, assigneeCounts, err := s.repo.GetStats(ctx, id)
	if err != nil {
		return nil, err
	}

	if statusCounts == nil {
		statusCounts = make([]StatusCount, 0)
	}
	if assigneeCounts == nil {
		assigneeCounts = make([]AssigneeCount, 0)
	}

	return &ProjectStatsResponse{
		StatusCounts:   statusCounts,
		AssigneeCounts: assigneeCounts,
	}, nil
}

func (s *projectService) Update(ctx context.Context, id, userID string, req UpdateProjectRequest) (*ProjectResponse, error) {
	p, _, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if p.OwnerID != userID {
		return nil, errors.New("forbidden")
	}

	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	p.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}

	return s.mapToResponse(p), nil
}

func (s *projectService) Delete(ctx context.Context, id, userID string) error {
	p, _, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if p.OwnerID != userID {
		return errors.New("forbidden")
	}

	return s.repo.Delete(ctx, id)
}

func (s *projectService) mapToResponse(p *Project) *ProjectResponse {
	return &ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		OwnerID:     p.OwnerID,
		CreatedAt:   p.CreatedAt,
	}
}
