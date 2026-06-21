 package repository

import (
    "database/sql"
    "rent/internal/models"
)

type AvatarRepository struct {
    Db *sql.DB
}

func NewAvatarRepository(db *sql.DB) *AvatarRepository {
    return &AvatarRepository{Db: db}
}

func (r *AvatarRepository) Create(image *models.Avatar) error {
    query := `
        INSERT INTO avatar (user_id, image_url)
        VALUES ($1, $2)
        RETURNING id, created_at
    `
    return r.Db.QueryRow(query, 
        image.UserID, 
        image.ImageURL, 
    ).Scan(&image.ID, &image.CreatedAt)
}

func (r *AvatarRepository) GetByUserID(userID int64) (*models.Avatar, error) {
    query := `
        SELECT id, user_id, image_url, created_at
        FROM avatar
        WHERE user_id = $1
    `
    rows, err := r.Db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var image models.Avatar
	err = rows.Scan(
		&image.ID,
		&image.UserID,
		&image.ImageURL,
        &image.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
    return &image, err
}

func (r *AvatarRepository) GetByID(id int64) (*models.Avatar, error) {
    query := `
        SELECT id, user_id, image_url, created_at
        FROM avatar
        WHERE id = $1
    `
    var img models.Avatar
    err := r.Db.QueryRow(query, id).Scan(
        &img.ID,
        &img.UserID,
        &img.ImageURL,
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

func (r *AvatarRepository) Delete(id int64) error {
    query := `DELETE FROM avatar WHERE id = $1`
    _, err := r.Db.Exec(query, id)
    return err
}

func (r *AvatarRepository) DeleteByUserID(userID int64) error {
    query := `DELETE FROM avatar WHERE user_id = $1`
    _, err := r.Db.Exec(query, userID)
    return err
}

func (r *AvatarRepository) DeleteByURL(imageURL string) error {
    query := `DELETE FROM avatar WHERE image_url = $1`
    _, err := r.Db.Exec(query, imageURL)
    return err
}
