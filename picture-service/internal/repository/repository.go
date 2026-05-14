package repository

import (
	"picture-service/internal/repository/db"
	"picture-service/internal/repository/minio"
)

type Repository struct {
	DatabaseRepository *db.DatabaseRepository
	MinioRepository    *minio.MinioRepository
}
