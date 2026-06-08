package repository

import (
	"database/sql"
	"fmt"
	"rent/internal/models"
	"time"
)

type BookingRepository struct {
	Db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{Db: db}
}

func (r *BookingRepository) Create(booking *models.Booking) error {
	query := `
        INSERT INTO booking (user_id, apartment_id, created_at, status, time_from, time_to)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at
    `

	err := r.Db.QueryRow(query,
		booking.UserID,
		booking.ApartmentID,
		booking.CreatedAt,
		booking.Status,
		booking.TimeFrom,
		booking.TimeTo,
	).Scan(&booking.ID, &booking.CreatedAt)

	return err
}

func (r *BookingRepository) GetByID(id int64) (*models.Booking, error) {
	query := `
        SELECT id, user_id, apartment_id, created_at, status, time_from, time_to
        FROM booking
        WHERE id = $1
    `

	var booking models.Booking
	err := r.Db.QueryRow(query, id).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.ApartmentID,
		&booking.CreatedAt,
		&booking.Status,
		&booking.TimeFrom,
		&booking.TimeTo,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *BookingRepository) GetByUser(userID int64, limit, offset int) ([]*models.Booking, error) {
	query := `
        SELECT id, user_id, apartment_id, created_at, status, time_from, time_to
        FROM booking
        WHERE user_id = $1
        ORDER BY time_from DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.Db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.ApartmentID,
			&booking.CreatedAt,
			&booking.Status,
			&booking.TimeFrom,
			&booking.TimeTo,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, &booking)
	}

	return bookings, rows.Err()
}

func (r *BookingRepository) GetByApartment(apartmentID int64, limit, offset int) ([]*models.Booking, error) {
	query := `
        SELECT id, user_id, apartment_id, created_at, status, time_from, time_to
        FROM booking
        WHERE apartment_id = $1
        ORDER BY time_from DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.Db.Query(query, apartmentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.ApartmentID,
			&booking.CreatedAt,
			&booking.Status,
			&booking.TimeFrom,
			&booking.TimeTo,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, &booking)
	}

	return bookings, rows.Err()
}

func (r *BookingRepository) CheckAvailability(apartmentID int64, timeFrom, timeTo time.Time) (bool, error) {
	query := `
        SELECT COUNT(*)
        FROM booking
        WHERE apartment_id = $1
        AND status NOT IN ('cancelled')
        AND (
            (time_from <= $2 AND time_to > $2) OR
            (time_from < $3 AND time_to >= $3) OR
            (time_from >= $2 AND time_to <= $3)
        )
    `

	var count int
	err := r.Db.QueryRow(query, apartmentID, timeFrom, timeTo).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *BookingRepository) UpdateStatus(id int64, status string) error {
	query := `UPDATE booking SET status = $1 WHERE id = $2`

	result, err := r.Db.Exec(query, status, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("booking with id %d not found", id)
	}

	return nil
}

func (r *BookingRepository) Cancel(id int64) error {
	return r.UpdateStatus(id, "")
}

func (r *BookingRepository) Delete(id int64) error {
	query := `DELETE FROM booking WHERE id = $1`

	result, err := r.Db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("booking with id %d not found", id)
	}

	return nil
}
