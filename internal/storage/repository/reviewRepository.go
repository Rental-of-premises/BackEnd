package repository

import (
	"database/sql"
	"fmt"
	"rent/internal/models"
)

type ReviewRepository struct {
	Db *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{Db: db}
}

func (r *ReviewRepository) Create(review *models.Review) error {
	query := `
		INSERT INTO reviews (user_id, apartment_id, comment, stars)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.Db.QueryRow(query,
		review.UserID,
		review.ApartmentID,
		review.Comment,
		review.Stars,
	).Scan(&review.ID, &review.CreatedAt)

	return err
}

func (r *ReviewRepository) GetByID(id int64) (*models.Review, error) {
	query := `
		SELECT id, user_id, apartment_id, comment, stars, created_at
		FROM reviews
		WHERE id = $1
	`

	var review models.Review
	err := r.Db.QueryRow(query, id).Scan(
		&review.ID,
		&review.UserID,
		&review.ApartmentID,
		&review.Comment,
		&review.Stars,
		&review.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &review, nil
}



func (r *ReviewRepository) GetByApartment(apartmentID int64, limit, offset int) ([]*models.Review, error) {
	query := `
		SELECT id, user_id, apartment_id, comment, stars, created_at
		FROM reviews
		WHERE apartment_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.Db.Query(query, apartmentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*models.Review
	for rows.Next() {
		var review models.Review
		err := rows.Scan(
			&review.ID,
			&review.UserID,
			&review.ApartmentID,
			&review.Comment,
			&review.Stars,
			&review.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}

	return reviews, rows.Err()
}

func (r *ReviewRepository) GetAverageStars(apartmentID int64) (float64, error) {
	query := `
		SELECT COALESCE(AVG(stars), 0) as average
		FROM reviews
		WHERE apartment_id = $1
	`

	var avg float64
	err := r.Db.QueryRow(query, apartmentID).Scan(&avg)
	return avg, err
}

func (r *ReviewRepository) Delete(id int64) error {
	query := `DELETE FROM reviews WHERE id = $1`

	result, err := r.Db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("review with id %d not found", id)
	}

	return nil
}
