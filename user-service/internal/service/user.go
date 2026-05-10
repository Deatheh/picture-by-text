package service

import (
	"errors"
	"user-service/internal/entities"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository *repository.Repository
}

func (s *UserService) IsEmailExists(email string) (bool, error) {
	return s.repository.DatabaseRepository.ExistsByEmail(email)
}

func (s *UserService) Add(email, plainPassword string) (*entities.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	return s.repository.DatabaseRepository.Add(email, string(hashed))
}
