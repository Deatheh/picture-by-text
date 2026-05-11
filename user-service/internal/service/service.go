package service

import (
	"context"
	"user-service/internal/config"
	"user-service/internal/entities"
	"user-service/internal/repository"
)

type User interface {
	Add(email, password string) (*entities.User, error)
	IsEmailExists(email string) (bool, error)
	GetByEmail(email string) (*entities.User, error)
	CheckPassword(hashedPassword, plainPassword string) bool
	GetByID(id string) (*entities.User, error)
	SaveRefreshToken(ctx context.Context, userID, token string, ttlSeconds int) error
	GetRefreshToken(ctx context.Context, userID string) (string, error)
}

type Service struct {
	User
}

func NewService(repo *repository.Repository, envConf *config.Config) *Service {
	return &Service{
		User: &UserService{repository: repo, envConf: envConf},
	}
}
