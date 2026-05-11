package service

import (
	"context"
	"errors"
	"user-service/internal/config"
	"user-service/internal/entities"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository *repository.Repository
	envConf    *config.Config
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

func (s *UserService) GetByEmail(email string) (*entities.User, error) {
	return s.repository.DatabaseRepository.GetByEmail(email)
}

func (s *UserService) CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func (s *UserService) GetByID(id string) (*entities.User, error) {
	return s.repository.DatabaseRepository.GetByID(id)
}

func (s *UserService) SaveRefreshToken(ctx context.Context, userID, token string, ttlSeconds int) error {
	return s.repository.Cache.SaveRefreshToken(ctx, userID, token, ttlSeconds)
}

func (s *UserService) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	return s.repository.Cache.GetRefreshToken(ctx, userID)
}
