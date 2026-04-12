package projects

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(userID string, req CreateProjectRequest) (*ProjectResponse, error)
	List(userID string) ([]ProjectResponse, error)
	GetDetails(id string) (*ProjectDetailsResponse, error)
	Update(id, userID string, req UpdateProjectRequest) (*ProjectResponse, error)
	Delete(id, userID string) error
}

type projectService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &projectService{repo: repo}
}

func (s *projectService) Create(userID string, req CreateProjectRequest) (*ProjectResponse, error) {
	p := &Project{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(p); err != nil {
		return nil, err
	}

	return s.mapToResponse(p), nil
}

func (s *projectService) List(userID string) ([]ProjectResponse, error) {
	projects, err := s.repo.List(userID)
	if err != nil {
		return nil, err
	}

	var res []ProjectResponse
	for _, p := range projects {
		res = append(res, *s.mapToResponse(&p))
	}
	
	if res == nil {
		res = make([]ProjectResponse, 0)
	}

	return res, nil
}

func (s *projectService) GetDetails(id string) (*ProjectDetailsResponse, error) {
	p, tasks, err := s.repo.GetByID(id)
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

func (s *projectService) Update(id, userID string, req UpdateProjectRequest) (*ProjectResponse, error) {
	p, _, err := s.repo.GetByID(id)
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

	if err := s.repo.Update(p); err != nil {
		return nil, err
	}

	return s.mapToResponse(p), nil
}

func (s *projectService) Delete(id, userID string) error {
	p, _, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if p.OwnerID != userID {
		return errors.New("forbidden")
	}

	return s.repo.Delete(id)
}

func (s *projectService) mapToResponse(p *Project) *ProjectResponse {
	return &ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		OwnerID:     p.OwnerID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
