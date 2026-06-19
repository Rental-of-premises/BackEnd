package repository

import (
	"database/sql"
	"fmt"
	"log"
	"rent/internal/models"
	"strings"
)

type ApartmentRepository struct {
	Db *sql.DB
}

func NewApartmentRepository(db *sql.DB) *ApartmentRepository {
	return &ApartmentRepository{Db: db}
}

func (r *ApartmentRepository) Create(apartment *models.Apartment) error {
	query := `
		INSERT INTO apartments (name, seller_id, description, capacity, price_per_hour, is_active, metro, address)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	err := r.Db.QueryRow(query,
		apartment.Name,
		apartment.SellerID,
		apartment.Description,
		apartment.Capacity,
		apartment.PricePerHour,
		apartment.IsActive,
		apartment.Metro,
		apartment.Address,
	).Scan(&apartment.ID, &apartment.CreatedAt)

	return err
}

func (r *ApartmentRepository) GetByID(id int64) (*models.Apartment, error) {
	query := `
		SELECT id, seller_id, name, description, capacity, price_per_hour, is_active, created_at, metro, address
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
		&apartment.Metro,
		&apartment.Address,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &apartment, nil
}

func (r *ApartmentRepository) GetAll(filter *models.ApartmentFilter) ([]*models.Apartment, error) {
	query := `
		SELECT id, seller_id, name, description, capacity, price_per_hour, is_active, created_at, metro, address
		FROM apartments
		WHERE 1=1
	`

	limit := filter.Limit
	offset := filter.Offset
	var args []interface{}
	argCounter := 1

	if filter.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argCounter)
		args = append(args, *filter.IsActive)
		argCounter++
	}

	if filter.SellerID != nil {
		query += fmt.Sprintf(" AND seller_id = $%d", argCounter)
		args = append(args, *filter.SellerID)
		argCounter++
	}

	if filter.MinPrice != nil {
		query += fmt.Sprintf(" AND price_per_hour >= $%d", argCounter)
		args = append(args, *filter.MinPrice)
		argCounter++
	}

	if filter.MaxPrice != nil {
		query += fmt.Sprintf(" AND price_per_hour <= $%d", argCounter)
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
			&apartment.Metro,
			&apartment.Address,
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
			price_per_hour = $4, is_active = $5, metro = $6, address = $7
		WHERE id = $8
	`

	result, err := r.Db.Exec(query,
		apartment.Name,
		apartment.Description,
		apartment.Capacity,
		apartment.PricePerHour,
		apartment.IsActive,
		apartment.Metro,
		apartment.Address,
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

func (r *ApartmentRepository) UpdatePartial(id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	log.Printf("📝 UpdatePartial: id=%d, updates=%+v", id, updates)

	// Разрешенные поля для обновления (только те, что есть в таблице apartments)
	// ⚠️ image_url НЕТ в этом списке - он хранится в отдельной таблице apartment_images
	allowedFields := map[string]bool{
		"name":           true,
		"description":    true,
		"capacity":       true,
		"price_per_hour": true,
		"is_active":      true,
		"metro":          true,
		"address":        true,
		"amenities":      true,
		// "image_url" НЕ ДОБАВЛЯТЬ - этой колонки нет в таблице apartments!
	}

	setParts := []string{}
	args := []interface{}{}
	i := 1

	for field, value := range updates {
		if !allowedFields[field] {
			log.Printf("⚠️ Поле %s пропущено (не разрешено для обновления в apartments)", field)
			continue
		}

		// Обработка nil значений
		if value == nil {
			setParts = append(setParts, fmt.Sprintf("%s = NULL", field))
			continue
		}

		// Обработка пустых строк
		if str, ok := value.(string); ok && str == "" {
			setParts = append(setParts, fmt.Sprintf("%s = NULL", field))
			continue
		}

		// Приводим типы
		switch field {
		case "price_per_hour":
			switch v := value.(type) {
			case float64:
				value = int64(v)
			case int:
				value = int64(v)
			}
		case "capacity":
			switch v := value.(type) {
			case float64:
				value = int16(v)
			case int:
				value = int16(v)
			}
		case "is_active":
			switch v := value.(type) {
			case bool:
				value = v
			default:
				return fmt.Errorf("неверный тип для is_active: %T", value)
			}
		}

		setParts = append(setParts, fmt.Sprintf("%s = $%d", field, i))
		args = append(args, value)
		i++
	}

	if len(setParts) == 0 {
		log.Printf("⚠️ Нет полей для обновления после фильтрации")
		return nil
	}

	args = append(args, id)

	query := fmt.Sprintf(
		"UPDATE apartments SET %s WHERE id = $%d",
		strings.Join(setParts, ", "),
		i,
	)

	log.Printf("📝 SQL запрос: %s", query)
	log.Printf("📝 Аргументы: %+v", args)

	_, err := r.Db.Exec(query, args...)
	if err != nil {
		log.Printf("❌ Ошибка выполнения SQL: %v", err)
		return fmt.Errorf("ошибка выполнения запроса: %w", err)
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
		SELECT id, seller_id, name, description, capacity, price_per_hour, is_active, created_at, metro, address
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
			&apartment.Metro,
			&apartment.Address,
		)
		if err != nil {
			return nil, err
		}
		apartments = append(apartments, &apartment)
	}

	return apartments, rows.Err()
}
