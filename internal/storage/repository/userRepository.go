package repository

import (
	"database/sql"
	"fmt"
	"rent/internal/models"
)

type UserRepository struct {
	Db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{Db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (name, password, email, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at
	`

	err := r.Db.QueryRow(query, user.Name, user.Password, user.Email).Scan(&user.ID, &user.CreatedAt)

	return err
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	query := `
		SELECT id, name, password, email, created_at
		FROM users
		WHERE id = $1
	`
	var user models.User
	err := r.Db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Password, &user.Email, &user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, name, password, email, created_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.Db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Password, &user.Email, &user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAll() ([]*models.User, error) {
	query := `
		SELECT id, name, password, email, created_at
		FROM users
		ORDER BY id
	`

	rows, err := r.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, rows.Err()
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET name = $1, password = $2, email = $3
		WHERE id = $4
	`

	result, err := r.Db.Exec(query, user.Name, user.Password, user.Email, user.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user with id %d not found", user.ID)
	}

	return nil
}

func (r *UserRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.Db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}
