package repository

import (
	"database/sql"
	"fmt"
	"rent/internal/models"
	"time"
	"errors"
)

type BookingRepository struct {
	Db *sql.DB
}


func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{Db: db}
}

func (r *BookingRepository) Create(booking *models.Booking) error {
	query := `
        INSERT INTO booking (user_id, apartment_id, status, time_from, time_to)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at
    `

	err := r.Db.QueryRow(query,
		booking.UserID,
		booking.ApartmentID,
		booking.Status,
		booking.TimeFrom,
		booking.TimeTo,
	).Scan(&booking.ID, &booking.CreatedAt)

	return err
}

func (r *BookingRepository) GetAll(filter *models.BookingFilter) ([]*models.Booking, error) {
	query := `
        SELECT b.id, b.user_id, b.apartment_id, b.status, b.time_from, b.time_to, b.created_at
        FROM booking b
        INNER JOIN apartments a ON b.apartment_id = a.id
        WHERE 1=1
    `
	if filter.Limit == nil || filter.Offset == nil {
		return nil, errors.New("missing Limit or Offset in Json for pagination")
	}
	limit := filter.Limit
	offset := filter.Offset

	var args []interface{}
	argCounter := 1

    if filter.Status != nil {
        query += fmt.Sprintf(" AND b.status = $%d", argCounter)
        args = append(args, *filter.Status)
        argCounter++
    }
    
    if filter.SellerID != nil {
        query += fmt.Sprintf(" AND a.seller_id = $%d", argCounter)
        args = append(args, *filter.SellerID)
        argCounter++
    }

    if filter.ApartmentID != nil {
        query += fmt.Sprintf(" AND a.id = $%d", argCounter)
        args = append(args, *filter.ApartmentID)
        argCounter++
    }
	
    if filter.UserID != nil {
        query += fmt.Sprintf(" AND b.user_id = $%d", argCounter)
        args = append(args, *filter.UserID)
        argCounter++
    }
    
    if filter.MinPrice != nil {
        query += fmt.Sprintf(" AND price_per_hour *  EXTRACT(EPOCH FROM (b.time_to - b.time_from)) / 3600 >= $%d", argCounter)
        args = append(args, *filter.MinPrice)
        argCounter++
    }
    
    if filter.MaxPrice != nil {
        query += fmt.Sprintf(" AND price_per_hour *  EXTRACT(EPOCH FROM (b.time_to - b.time_from)) / 3600 <= $%d", argCounter)
        args = append(args, *filter.MaxPrice)
        argCounter++
    }

	query += " ORDER BY id LIMIT $" + fmt.Sprintf("%d", argCounter) + " OFFSET $" + fmt.Sprintf("%d", argCounter+1)
	args = append(args, limit, offset)

	rows, err := r.Db.Query(query, args...)
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
			&booking.Status,
			&booking.TimeFrom,
			&booking.TimeTo,
			&booking.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, &booking)
	}

	return bookings, rows.Err()
}

func (r *BookingRepository)CheckReviewValidation(UserID, ApartmentID int64) (bool, error){
	limit := 9999999
	offset := 0
	status := "completed"
	filter := &models.BookingFilter{
		Limit : &limit,
		Offset : &offset,
		Status : &status,
		UserID : &UserID,
		ApartmentID : &ApartmentID,
	}
	bookings, err := r.GetAll(filter)
	return bookings != nil, err
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

func (r *BookingRepository) GetByUser(userID int64) ([]*models.Booking, error) {
	query := `
        SELECT id, user_id, apartment_id, created_at, status, time_from, time_to
        FROM booking
        WHERE user_id = $1
        ORDER BY time_from DESC
    `

	rows, err := r.Db.Query(query, userID)
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

func (r *BookingRepository) CheckAvailability(apartmentID int64, timeFrom, timeTo time.Time, userID int64) (bool, error) {
	query := `
        SELECT COUNT(*), a.seller_id
        FROM booking b
        JOIN apartments a ON b.apartment_id = a.id
        WHERE b.apartment_id = $1
        AND b.status NOT IN ('cancelled', 'rejected')
        AND (
            (b.time_from <= $2 AND b.time_to > $2) OR
            (b.time_from < $3 AND b.time_to >= $3) OR
            (b.time_from >= $2 AND b.time_to <= $3)
        )
        GROUP BY a.seller_id
    `

	var count int
	var sellerID int64
	err := r.Db.QueryRow(query, apartmentID, timeFrom, timeTo).Scan(&count, &sellerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}
	return count == 0 && sellerID != userID, nil
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
	return r.UpdateStatus(id, "cancelled")
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

func (r *BookingRepository) CompletePastBookings() error {
	query := `
		UPDATE booking 
		SET status = 'completed' 
		WHERE status = 'confirmed' 
		AND time_to < NOW()
	`
	result, err := r.Db.Exec(query)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rows > 0 {
		fmt.Printf("✅ Автоматически завершено %d бронирований\n", rows)
	}

	return nil
}
