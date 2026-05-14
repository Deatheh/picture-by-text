package service

import (
	"context"
	"picture-service/internal/ai"
	"picture-service/internal/config"
	"picture-service/internal/repository"
)

type Generation interface {
	StartGeneration(ctx context.Context, userID, text string) (string, error)
	GetTaskStatus(ctx context.Context, taskID string) (string, []string, error)
}

type Service struct {
	Generation
}

func NewService(repo *repository.Repository, cfg *config.Config, aiRepo *ai.AiRepository) *Service {
	return &Service{
		Generation: NewGenerationService(repo, cfg, aiRepo),
	}
}
