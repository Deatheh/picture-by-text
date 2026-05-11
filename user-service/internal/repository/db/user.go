package db

import (
	"fmt"
	"user-service/internal/entities"

	"github.com/google/uuid"
)

func (dbr *DatabaseRepository) ExistsByEmail(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`
	err := dbr.DB.Get(&exists, query, email)
	if err != nil {
		return false, fmt.Errorf("check email existence: %w", err)
	}
	return exists, nil
}

func (dbr *DatabaseRepository) Add(email, passwordHash string) (*entities.User, error) {
	id := uuid.New().String()

	query := `
		INSERT INTO users (id, email, password, role_id, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email, password, role_id, is_active
	`

	var user entities.User
	err := dbr.DB.QueryRowx(query, id, email, passwordHash, 2, true).StructScan(&user)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}
	return &user, nil
}

func (dbr *DatabaseRepository) GetByEmail(email string) (*entities.User, error) {
	var user entities.User
	query := `SELECT id, email, password, role_id, is_active FROM users WHERE email=$1`
	err := dbr.DB.Get(&user, query, email)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &user, nil
}

func (dbr *DatabaseRepository) GetByID(id string) (*entities.User, error) {
	var user entities.User
	query := `SELECT id, email, password, role_id, is_active FROM users WHERE id=$1`
	err := dbr.DB.Get(&user, query, id)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}
