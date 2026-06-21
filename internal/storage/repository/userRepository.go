package repository

import (
	"database/sql"
	"fmt"
	"rent/internal/models"
	"log"
	"strings"
)

type UserRepository struct {
	Db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{Db: db}
}

func (r *UserRepository) Create(user *models.User) error {
    query := `
        INSERT INTO users (name, password, email, is_active, email_token)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

    err := r.Db.QueryRow(query,
        user.Name,
        user.Password,
        user.Email,
        user.IsActive,
        user.EmailToken,
    ).Scan(&user.ID)

    if err != nil {
        return err
    }


    return nil
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
		SELECT id, name, password, email, is_active, email_token, created_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.Db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Password, &user.Email, &user.IsActive, &user.EmailToken, &user.CreatedAt,
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
    query := "UPDATE users SET"
    args := []interface{}{}
    argCounter := 1

    if user.Name != "" {
        query += " name = " + fmt.Sprintf("$%d", argCounter) + ","
        args = append(args, user.Name)
        argCounter++
    }

    if user.Password != "" {
        query += " password = " + fmt.Sprintf("$%d", argCounter) + ","
        args = append(args, user.Password)
        argCounter++
    }

    if user.Email != "" {
        query += " email = " + fmt.Sprintf("$%d", argCounter) + ","
        args = append(args, user.Email)
        argCounter++
    }

    query = strings.TrimSuffix(query, ",")

    query += " WHERE id = " + fmt.Sprintf("$%d", argCounter)
    args = append(args, user.ID)

    result, err := r.Db.Exec(query, args...)
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

func (r *UserRepository) UpdateEmailToken(userID int64, token string) error {
    log.Printf("📝 UpdateEmailToken: userID=%d, token=%s", userID, token)
    query := `UPDATE users SET email_token = $1 WHERE id = $2`
    result, err := r.Db.Exec(query, token, userID)
    if err != nil {
        log.Printf("❌ Ошибка UPDATE: %v", err)
        return err
    }
    rows, _ := result.RowsAffected()
    log.Printf("✅ Обновлено строк: %d", rows)
    return nil
}

func (r *UserRepository) ActivateUser(token string) (int64, error) {
    var userID int64
    query := `UPDATE users SET is_active = true, email_token = NULL WHERE email_token = $1 RETURNING id`
    err := r.Db.QueryRow(query, token).Scan(&userID)
    if err == sql.ErrNoRows {
        return 0, nil
    }
    return userID, err
}