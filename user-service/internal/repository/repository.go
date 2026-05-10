package repository

import "user-service/internal/repository/db"

type Repository struct {
	DatabaseRepository db.DatabaseRepository
}
