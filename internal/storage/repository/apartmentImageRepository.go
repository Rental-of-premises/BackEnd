package repository

import (
	"database/sql"
	"log"
	"rent/internal/models"
)

type ApartmentImageRepository struct {
	Db *sql.DB
}

func NewApartmentImageRepository(db *sql.DB) *ApartmentImageRepository {
	return &ApartmentImageRepository{Db: db}
}

func (r *ApartmentImageRepository) Create(image *models.ApartmentImage) error {
	query := `
		INSERT INTO apartment_images (apartment_id, image_data, position)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.Db.QueryRow(query,
		image.ApartmentID,
		image.ImageData,
		image.Position,
	).Scan(&image.ID, &image.CreatedAt)
}

func (r *ApartmentImageRepository) GetByApartmentID(apartmentID int64) ([]*models.ApartmentImage, error) {
	query := `
		SELECT id, apartment_id, image_data, position, created_at
		FROM apartment_images
		WHERE apartment_id = $1
		ORDER BY position ASC
	`
	rows, err := r.Db.Query(query, apartmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*models.ApartmentImage
	for rows.Next() {
		var img models.ApartmentImage
		err := rows.Scan(
			&img.ID,
			&img.ApartmentID,
			&img.ImageData,
			&img.Position,
			&img.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, &img)
	}
	return images, rows.Err()
}

func (r *ApartmentImageRepository) GetByID(id int64) (*models.ApartmentImage, error) {
	query := `
		SELECT id, apartment_id, image_data, position, created_at
		FROM apartment_images
		WHERE id = $1
	`
	var img models.ApartmentImage
	err := r.Db.QueryRow(query, id).Scan(
		&img.ID,
		&img.ApartmentID,
		&img.ImageData,
		&img.Position,
		&img.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *ApartmentImageRepository) Delete(id int64) error {
	query := `DELETE FROM apartment_images WHERE id = $1`
	_, err := r.Db.Exec(query, id)
	return err
}

func (r *ApartmentImageRepository) DeleteByApartmentID(apartmentID int64) error {
	log.Printf("🗑️ DeleteByApartmentID: удаление записей для помещения %d", apartmentID)
	query := `DELETE FROM apartment_images WHERE apartment_id = $1`
	result, err := r.Db.Exec(query, apartmentID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	log.Printf("✅ Удалено %d записей из БД", rows)
	return nil
}

func (r *ApartmentImageRepository) DeleteByURL(imageURL string) error {
	query := `DELETE FROM apartment_images WHERE image_data = $1`
	_, err := r.Db.Exec(query, imageURL)
	return err
}

func (r *ApartmentImageRepository) UpdatePosition(id int64, position int) error {
	query := `UPDATE apartment_images SET position = $1 WHERE id = $2`
	_, err := r.Db.Exec(query, position, id)
	return err
}

func (r *ApartmentImageRepository) GetMainImage(apartmentID int64) (*models.ApartmentImage, error) {
	query := `
		SELECT id, apartment_id, image_data, position, created_at
		FROM apartment_images
		WHERE apartment_id = $1 AND position = 0
	`
	var img models.ApartmentImage
	err := r.Db.QueryRow(query, apartmentID).Scan(
		&img.ID,
		&img.ApartmentID,
		&img.ImageData,
		&img.Position,
		&img.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *ApartmentImageRepository) CountByApartmentID(apartmentID int64) (int, error) {
	query := `SELECT COUNT(*) FROM apartment_images WHERE apartment_id = $1`
	var count int
	err := r.Db.QueryRow(query, apartmentID).Scan(&count)
	return count, err
}
