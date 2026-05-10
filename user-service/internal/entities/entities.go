package entities

type User struct {
	Uuid     string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	RoleID   int    `db:"role_id"`
	IsActive bool   `db:"is_active"`
}
