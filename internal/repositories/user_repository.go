package repositories

import (
	"cusror_ai/internal/models"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`
	var user models.User
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	query := `INSERT INTO users (name, email,password) VALUES ($1, $2,$3) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(query, user.Name, user.Email,user.Password).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(user *models.User) (*models.User, error) {
	query := `UPDATE users SET name = $1, email = $2, updated_at = NOW() WHERE id = $3 RETURNING created_at, updated_at`
	err := r.db.QueryRow(query, user.Name, user.Email, user.ID).Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	row := r.db.QueryRow("SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1", email)
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}
