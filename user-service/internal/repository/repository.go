package repository

import (
	"user-service/internal/repository/cache"
	"user-service/internal/repository/db"
)

type Repository struct {
	DatabaseRepository *db.DatabaseRepository
	Cache              *cache.RedisClient
}
