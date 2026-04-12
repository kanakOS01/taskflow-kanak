package auth

import (
	"context"
	"errors"
	"time"

	"taskflow/internal/users"
	"taskflow/internal/utils"

	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (string, *users.User, error)
	Login(ctx context.Context, req LoginRequest) (string, *users.User, error)
}

type authService struct {
	userRepo  users.Repository
	jwtSecret string
}

func NewService(userRepo users.Repository, jwtSecret string) Service {
	return &authService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) (string, *users.User, error) {
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return "", nil, err
	}

	user := &users.User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  hash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return "", nil, err
	}

	token, err := utils.GenerateToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (string, *users.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
