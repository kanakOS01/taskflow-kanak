package auth

import (
	"errors"
	"time"

	"taskflow/internal/users"
	"taskflow/internal/utils"

	"github.com/google/uuid"
)

type Service interface {
	Register(req RegisterRequest) error
	Login(req LoginRequest) (string, error)
}

type authService struct {
	userRepo  users.Repository
	jwtSecret string
}

func NewService(userRepo users.Repository, jwtSecret string) Service {
	return &authService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *authService) Register(req RegisterRequest) error {
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &users.User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  hash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.userRepo.Create(user)
}

func (s *authService) Login(req LoginRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return utils.GenerateToken(user.ID, user.Email, s.jwtSecret)
}
