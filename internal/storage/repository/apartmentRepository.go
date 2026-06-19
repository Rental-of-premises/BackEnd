package repository

import (
    "database/sql"
    "fmt"
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
        INSERT INTO apartments (name, seller_id, description, capacity, price_per_hour, is_active, address, metro)
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
        apartment.Address,
        apartment.Metro,
    ).Scan(&apartment.ID, &apartment.CreatedAt)

	if len(apartment.Amenities) > 0 {
        for _, amenity := range apartment.Amenities {
            _, err = r.Db.Exec(
                `INSERT INTO apartment_amenities (apartment_id, amenity_id) VALUES ($1, $2)`,
                apartment.ID, amenity.ID,
            )
            if err != nil {
                return err
            }
        }
    }

    return err
}

func (r *ApartmentRepository) GetByID(id int64) (*models.Apartment, error) {
    query := `
        SELECT id, seller_id, name, description, capacity, price_per_hour, is_active, address, metro, created_at
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
        &apartment.Address,
        &apartment.Metro,
        &apartment.CreatedAt,
    )

    if err != nil {
        return nil, err
    }
    amenitiesQuery := `
    SELECT a.id, a.name, a.icon
    FROM amenities a
    JOIN apartment_amenities aa ON a.id = aa.amenity_id
    WHERE aa.apartment_id = $1
    ORDER BY a.name
    `
    rows, err := r.Db.Query(amenitiesQuery, id)
    if err != nil {
        return &apartment, err
    }
    defer rows.Close()

    var amenities []*models.Amenity
    for rows.Next() {
        var am models.Amenity
        err := rows.Scan(&am.ID, &am.Name, &am.Icon)
        if err != nil {
            return &apartment, err
        }
        amenities = append(amenities, &am)
    }
    apartment.Amenities = amenities

    return &apartment, nil
}

func (r *ApartmentRepository) GetAll(filter *models.ApartmentFilter) ([]*models.Apartment, error) {
    query := `
        SELECT a.id
        FROM apartments a
    `

    var args []interface{}
    argCounter := 1
    conditions := []string{}

    if filter != nil && len(filter.Amenities) > 0 {
        query += " JOIN apartment_amenities aa ON a.id = aa.apartment_id"
        placeholders := []string{}
        for _, amenityID := range filter.Amenities {
            placeholders = append(placeholders, fmt.Sprintf("$%d", argCounter))
            args = append(args, amenityID)
            argCounter++
        }
        conditions = append(conditions, fmt.Sprintf("aa.amenity_id IN (%s)", strings.Join(placeholders, ", ")))
    }

    query += " WHERE 1=1"

    if filter != nil {
        if filter.IsActive != nil {
            conditions = append(conditions, fmt.Sprintf("a.is_active = $%d", argCounter))
            args = append(args, *filter.IsActive)
            argCounter++
        }

        if filter.SellerID != nil {
            conditions = append(conditions, fmt.Sprintf("a.seller_id = $%d", argCounter))
            args = append(args, *filter.SellerID)
            argCounter++
        }

        if filter.MinPrice != nil {
            conditions = append(conditions, fmt.Sprintf("a.price_per_hour >= $%d", argCounter))
            args = append(args, *filter.MinPrice)
            argCounter++
        }

        if filter.MaxPrice != nil {
            conditions = append(conditions, fmt.Sprintf("a.price_per_hour <= $%d", argCounter))
            args = append(args, *filter.MaxPrice)
            argCounter++
        }
    }

    for _, cond := range conditions {
        query += " AND " + cond
    }

    limit := 10
    offset := 0
    if filter != nil {
        if filter.Limit != nil {
            limit = *filter.Limit
        }
        if filter.Offset != nil {
            offset = *filter.Offset
        }
    }

    query += " ORDER BY a.id LIMIT $" + fmt.Sprintf("%d", argCounter) + " OFFSET $" + fmt.Sprintf("%d", argCounter+1)
    args = append(args, limit, offset)

    rows, err := r.Db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var apartments []*models.Apartment
    for rows.Next() {
        var apartmentID int64
        err := rows.Scan(
            &apartmentID,
        )
        if err != nil {
            return nil, err
        }

        apartment, err := r.GetByID(apartmentID)
        if err != nil {
            return nil, err
        }

        apartments = append(apartments, apartment)
    }

    return apartments, rows.Err()
}

func (r *ApartmentRepository) UpdatePartial(id int64, updates map[string]interface{}) error {
    if len(updates) == 0 {
        return nil
    }

    var amenityIDs []int64

    if amenitiesVal, ok := updates["amenities"]; ok {
        if arr, ok := amenitiesVal.([]interface{}); ok {
            amenityIDs = make([]int64, 0, len(arr))
            for _, v := range arr {
                if id, ok := v.(float64); ok {
                    amenityIDs = append(amenityIDs, int64(id))
                }
            }
        }
        delete(updates, "amenities")
    }

    if len(updates) > 0 {
        setParts := []string{}
        args := []interface{}{}
        i := 1

        for field, value := range updates {
            switch field {
            case "name":
                if v, ok := value.(string); ok {
                    value = v
                }
            case "description":
                if v, ok := value.(string); ok {
                    value = v
                }
            case "address":
                if v, ok := value.(string); ok {
                    value = v
                }
            case "metro":
                if v, ok := value.(string); ok {
                    value = v
                }
            case "price_per_hour":
                if v, ok := value.(float64); ok {
                    value = int64(v)
                }
            case "capacity":
                if v, ok := value.(float64); ok {
                    value = int16(v)
                }
            case "is_active":
                if v, ok := value.(bool); ok {
                    value = v
                }
            }

            setParts = append(setParts, fmt.Sprintf("%s = $%d", field, i))
            args = append(args, value)
            i++
        }

        if len(setParts) > 0 {
            args = append(args, id)
            query := fmt.Sprintf(
                "UPDATE apartments SET %s WHERE id = $%d",
                strings.Join(setParts, ", "),
                i,
            )

            _, err := r.Db.Exec(query, args...)
            if err != nil {
                return err
            }
        }
    }

    if len(amenityIDs) > 0 {
        _, err := r.Db.Exec(`DELETE FROM apartment_amenities WHERE apartment_id = $1`, id)
        if err != nil {
            return err
        }

        for _, amenityID := range amenityIDs {
            _, err = r.Db.Exec(
                `INSERT INTO apartment_amenities (apartment_id, amenity_id) VALUES ($1, $2)`,
                id, amenityID,
            )
            if err != nil {
                return err
            }
        }
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
        SELECT id, seller_id, name, description, capacity, price_per_hour, is_active, address, metro, created_at
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
            &apartment.Address,
            &apartment.Metro,
            &apartment.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        apartments = append(apartments, &apartment)
    }

    return apartments, rows.Err()
}


func (r *ApartmentRepository) GetAmenityByID(id int64) (*models.Amenity, error) {
    query := `
        SELECT id, name, icon
        FROM amenities
        WHERE id = $1
    `

    var amenity models.Amenity
    err := r.Db.QueryRow(query, id).Scan(
        &amenity.ID,
        &amenity.Name,
        &amenity.Icon,
    )

    if err != nil {
        return nil, err
    }

    return &amenity, nil
}


func (r *ApartmentRepository) GetAllAmenities() ([]*models.Amenity, error) {
    query := `
        SELECT a.id, a.name, a.icon
        FROM amenities a
    `

    rows, err := r.Db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var amenities []*models.Amenity
    for rows.Next() {
        var amenity models.Amenity
        err := rows.Scan(
            &amenity.ID, 
            &amenity.Name,
            &amenity.Icon,
	    )

        if err != nil {
            return nil, err
        }

        amenities = append(amenities, &amenity)
    }

    return amenities, nil
}