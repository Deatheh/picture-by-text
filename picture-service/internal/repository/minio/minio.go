package minio

import (
	"bytes"
	"context"
	"fmt"
	"picture-service/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioRepository struct {
	Client *minio.Client
	Bucket string
}

func NewMinioRepository(cfg *config.Config) (*MinioRepository, error) {
	client, err := minio.New(cfg.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.MinIO.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.MinIO.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &MinioRepository{
		Client: client,
		Bucket: cfg.MinIO.BucketName,
	}, nil
}

func (r *MinioRepository) UploadImage(ctx context.Context, taskID, sceneID string, imageData []byte) (string, error) {
	objectName := fmt.Sprintf("%s/%s.png", taskID, sceneID)
	contentType := "image/png"

	// Используем bytes.NewReader вместо nil
	reader := bytes.NewReader(imageData)

	_, err := r.Client.PutObject(ctx, r.Bucket, objectName, reader, int64(len(imageData)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	return fmt.Sprintf("/%s/%s", r.Bucket, objectName), nil
}

func (r *MinioRepository) GetImageURL(ctx context.Context, taskID, sceneID string) (string, error) {
	objectName := fmt.Sprintf("%s/%s.png", taskID, sceneID)

	// Проверяем, существует ли объект
	_, err := r.Client.StatObject(ctx, r.Bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("image not found: %w", err)
	}

	// Возвращаем путь (URL можно сгенерировать отдельно через PresignedGetObject)
	return fmt.Sprintf("/%s/%s", r.Bucket, objectName), nil
}
