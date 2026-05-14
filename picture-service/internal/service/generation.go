package service

import (
	"context"
	"fmt"
	"picture-service/internal/ai"
	"picture-service/internal/config"
	"picture-service/internal/entities"
	"picture-service/internal/repository"

	"github.com/google/uuid"
)

type GenerationService struct {
	repo   *repository.Repository
	cfg    *config.Config
	aiRepo *ai.AiRepository
}

func NewGenerationService(repo *repository.Repository, cfg *config.Config, aiRepo *ai.AiRepository) *GenerationService {
	return &GenerationService{
		repo:   repo,
		cfg:    cfg,
		aiRepo: aiRepo,
	}
}

func (s *GenerationService) StartGeneration(ctx context.Context, userID, text string) (string, error) {
	// 1. Создаём задачу
	task, err := s.repo.DatabaseRepository.CreateTask(userID)
	if err != nil {
		return "", fmt.Errorf("failed to create task: %w", err)
	}

	// 2. Запускаем генерацию в горутине с переданным context
	go s.processGeneration(context.Background(), task.ID, text)

	return task.ID, nil
}

func (s *GenerationService) processGeneration(ctx context.Context, taskID, text string) {
	// Обновляем статус на processing
	s.repo.DatabaseRepository.UpdateTaskStatus(taskID, "processing")

	// Разбиваем текст на сцены
	scenes, err := s.aiRepo.SplitTextToScenes(ctx, text)
	if err != nil {
		s.repo.DatabaseRepository.UpdateTaskStatus(taskID, "failed")
		return
	}

	// Для каждой сцены генерируем изображение
	for i, sceneText := range scenes {
		// Генерируем промпт
		prompt, err := s.aiRepo.GeneratePrompt(ctx, sceneText)
		if err != nil {
			continue
		}

		// Генерируем изображение
		imageData, err := s.aiRepo.GenerateImage(ctx, prompt)
		if err != nil {
			continue
		}

		// Сохраняем в MinIO
		sceneID := uuid.New().String()
		imageURL, err := s.repo.MinioRepository.UploadImage(ctx, taskID, sceneID, imageData)
		if err != nil {
			continue
		}

		// Сохраняем сцену в БД
		scene := &entities.Scene{
			ID:       sceneID,
			TaskID:   taskID,
			Index:    i + 1,
			Text:     sceneText,
			ImageURL: imageURL,
		}
		s.repo.DatabaseRepository.CreateScene(scene)
	}

	// Обновляем статус на completed
	s.repo.DatabaseRepository.UpdateTaskStatus(taskID, "completed")
}

func (s *GenerationService) GetTaskStatus(ctx context.Context, taskID string) (string, []string, error) {
	task, err := s.repo.DatabaseRepository.GetTask(taskID)
	if err != nil {
		return "", nil, fmt.Errorf("task not found: %w", err)
	}

	scenes, err := s.repo.DatabaseRepository.GetScenesByTaskID(taskID)
	if err != nil {
		return task.Status, nil, nil
	}

	var imageURLs []string
	for _, scene := range scenes {
		imageURLs = append(imageURLs, scene.ImageURL)
	}

	return task.Status, imageURLs, nil
}
