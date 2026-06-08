package repository

import (
	"database/sql"
	"fmt"
	"rent/internal/models"
)

type ApartmentRepository struct {
	Db *sql.DB
}

func NewApartmentRepository(db *sql.DB) *ApartmentRepository {
	return &ApartmentRepository{Db: db}
}

func (r *ApartmentRepository) Create(apartment *models.Apartment) error {
	query := `
		INSERT INTO apartments (seller_id, name, description, capacity, price_per_hour, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	err := r.Db.QueryRow(query,
		apartment.SellerID,
		apartment.Name,
		apartment.Description,
		apartment.Capacity,
		apartment.PricePerHour,
		apartment.IsActive,
		apartment.CreatedAt,
	).Scan(&apartment.ID, &apartment.CreatedAt)

	return err
}

func (r *ApartmentRepository) GetByID(id int64) (*models.Apartment, error) {
	query := `
		SELECT id, seller_id, name, description, capacity, price_per_hour, is_active, created_at
		FROM apartments
		WHERE id = $1
	 `

	var apartment models.Apartment
	err := r.Db.QueryRow(query, id).Scan(
		&apartment.ID,
		&apartment.SellerID,
		&apartment.Name,
		&apartment.Description,
		&apartment.Capacity,
		&apartment.PricePerHour,
		&apartment.IsActive,
		&apartment.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &apartment, nil
}

func (r *ApartmentRepository) GetAll(filter map[string]interface{}, limit, offset int) ([]*models.Apartment, error) {
	query := `
		SELECT id, seller_id, name, description, capacity, price_per_hour, is_active, created_at
		FROM apartments
		WHERE 1=1
	`
	var args []interface{}
	argCounter := 1

	if isActive, ok := filter["is_active"].(bool); ok {
		query += fmt.Sprintf(" AND is_active = $%d", argCounter)
		args = append(args, isActive)
		argCounter++
	}

	if sellerID, ok := filter["seller_id"].(int64); ok {
		query += fmt.Sprintf(" AND seller_id = $%d", argCounter)
		args = append(args, sellerID)
		argCounter++
	}

	if minPrice, ok := filter["min_price"].(int32); ok {
		query += fmt.Sprintf(" AND price_per_hour >= $%d", argCounter)
		args = append(args, minPrice)
		argCounter++
	}

	if maxPrice, ok := filter["max_price"].(int32); ok {
		query += fmt.Sprintf(" AND price_per_hour <= $%d", argCounter)
		args = append(args, maxPrice)
		argCounter++
	}

	query += " ORDER BY id LIMIT $" + fmt.Sprintf("%d", argCounter) + " OFFSET $" + fmt.Sprintf("%d", argCounter+1)
	args = append(args, limit, offset)

	rows, err := r.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apartments []*models.Apartment
	for rows.Next() {
		var apartment models.Apartment
		err := rows.Scan(
			&apartment.ID,
			&apartment.SellerID,
			&apartment.Name,
			&apartment.Description,
			&apartment.Capacity,
			&apartment.PricePerHour,
			&apartment.IsActive,
			&apartment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		apartments = append(apartments, &apartment)
	}

	return apartments, rows.Err()
}

func (r *ApartmentRepository) Update(apartment *models.Apartment) error {
	query := `
		UPDATE apartments
		SET name = $1, description = $2, capacity = $3, 
				price_per_hour = $4, is_active = $5
		WHERE id = $6
	`

	result, err := r.Db.Exec(query,
		apartment.Name,
		apartment.Description,
		apartment.Capacity,
		apartment.PricePerHour,
		apartment.IsActive,
		apartment.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("apartment with id %d not found", apartment.ID)
	}

	return nil
}

func (r *ApartmentRepository) Delete(id int64) error {
	query := `DELETE FROM apartments WHERE id = $1`

	result, err := r.Db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("apartment with id %d not found", id)
	}

	return nil
}

func (r *ApartmentRepository) GetBySeller(sellerID int64) ([]*models.Apartment, error) {
	query := `
		SELECT id, seller_id, name, description, capacity, price_per_hour, is_active, created_at
		FROM apartments
		WHERE seller_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.Db.Query(query, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apartments []*models.Apartment
	for rows.Next() {
		var apartment models.Apartment
		err := rows.Scan(
			&apartment.ID,
			&apartment.SellerID,
			&apartment.Name,
			&apartment.Description,
			&apartment.Capacity,
			&apartment.PricePerHour,
			&apartment.IsActive,
			&apartment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		apartments = append(apartments, &apartment)
	}

	return apartments, rows.Err()
}
