package service

import (
	"user-service/internal/config"
	"user-service/internal/entities"
	"user-service/internal/repository"
)

type User interface {
	Add(email, password string) (*entities.User, error)
	IsEmailExists(email string) (bool, error)
}

type Service struct {
	User
}

func NewService(repo *repository.Repository, envConf *config.Config) *Service {
	return &Service{
		User: &UserService{repository: repo},
	}
}
