package service

import (
	"context"
	"user-service/internal/config"
	"user-service/internal/repository"
)

type CacheService struct {
	repository *repository.Repository
	envConf    *config.Config
}

func (s *CacheService) SaveRefreshToken(ctx context.Context, userID, token string, ttlSeconds int) error {
	return s.repository.Cache.SaveRefreshToken(ctx, userID, token, ttlSeconds)
}

// GetRefreshToken возвращает сохранённый refresh_token для пользователя
func (s *CacheService) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	return s.repository.Cache.GetRefreshToken(ctx, userID)
}

// DeleteRefreshToken удаляет refresh_token (при логауте)
func (s *CacheService) DeleteRefreshToken(ctx context.Context, userID string) error {
	return s.repository.Cache.DeleteRefreshToken(ctx, userID)
}
