package db

import (
	"fmt"
	"picture-service/internal/entities"

	"github.com/google/uuid"
)

func (r *DatabaseRepository) CreateTask(userID string) (*entities.Task, error) {
	id := uuid.New().String()

	task := &entities.Task{
		ID:     id,
		UserID: userID,
		Status: "pending",
	}

	query := `INSERT INTO tasks (id, user_id, status) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, task.ID, task.UserID, task.Status)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}

func (r *DatabaseRepository) UpdateTaskStatus(id, status string) error {
	query := `UPDATE tasks SET status=$1 WHERE id=$2`
	_, err := r.DB.Exec(query, status, id)
	return err
}

func (r *DatabaseRepository) GetTask(id string) (*entities.Task, error) {
	var task entities.Task
	query := `SELECT id, user_id, status FROM tasks WHERE id=$1`
	err := r.DB.Get(&task, query, id)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *DatabaseRepository) CreateScene(scene *entities.Scene) error {
	scene.ID = uuid.New().String()
	query := `INSERT INTO scenes (id, task_id, scene_index, text, image_url) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.Exec(query, scene.ID, scene.TaskID, scene.Index, scene.Text, scene.ImageURL)
	return err
}

func (r *DatabaseRepository) GetScenesByTaskID(taskID string) ([]entities.Scene, error) {
	var scenes []entities.Scene
	query := `SELECT id, task_id, scene_index, text, image_url FROM scenes WHERE task_id=$1 ORDER BY scene_index`
	err := r.DB.Select(&scenes, query, taskID)
	return scenes, err
}
