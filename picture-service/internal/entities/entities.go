package entities

type Task struct {
	ID     string `db:"id"`
	UserID string `db:"user_id"`
	Status string `db:"status"`
}

type Scene struct {
	ID       string `db:"id"`
	TaskID   string `db:"task_id"`
	Index    int    `db:"scene_index"`
	Text     string `db:"text"`
	ImageURL string `db:"image_url"`
}
