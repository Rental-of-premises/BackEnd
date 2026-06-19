 package repository

// import (
//     "database/sql"
//     "rent/internal/models"
// )

// type AvatarRepository struct {
//     Db *sql.DB
// }

// func NewAvatarRepository(db *sql.DB) *AvatarRepository {
//     return &AvatarRepository{Db: db}
// }

// func (r *AvatarRepository) Create(image *models.Avatar) error {
//     query := `
//         INSERT INTO avatar (user_id, image_url)
//         VALUES ($1, $2)
//         RETURNING id
//     `
//     return r.Db.QueryRow(query, 
//         image.ApartmentID, 
//         image.ImageURL, 
//     ).Scan(&image.ID)
// }

// func (r *AvatarRepository) GetByUserID(userID int64) ([]*models.Avatar, error) {
//     query := `
//         SELECT id, apartment_id, image_url
//         FROM apartment_images
//         WHERE user_id = $1
//     `
//     rows, err := r.Db.Query(query, apartmentID)
//     if err != nil {
//         return nil, err
//     }
//     defer rows.Close()

//     var images []*models.Avatar
//     for rows.Next() {
//         var img models.Avatar
//         err := rows.Scan(
//             &img.ID,
//             &img.ApartmentID,
//             &img.ImageURL,
//             &img.Position,
//             &img.CreatedAt,
//         )
//         if err != nil {
//             return nil, err
//         }
//         images = append(images, &img)
//     }
//     return images, rows.Err()
// }

// func (r *AvatarRepository) GetByID(id int64) (*models.Avatar, error) {
//     query := `
//         SELECT id, apartment_id, image_url, position, created_at
//         FROM apartment_images
//         WHERE id = $1
//     `
//     var img models.Avatar
//     err := r.Db.QueryRow(query, id).Scan(
//         &img.ID,
//         &img.ApartmentID,
//         &img.ImageURL,
//         &img.Position,
//         &img.CreatedAt,
//     )
//     if err == sql.ErrNoRows {
//         return nil, nil
//     }
//     if err != nil {
//         return nil, err
//     }
//     return &img, nil
// }

// func (r *AvatarRepository) Delete(id int64) error {
//     query := `DELETE FROM apartment_images WHERE id = $1`
//     _, err := r.Db.Exec(query, id)
//     return err
// }

// func (r *AvatarRepository) DeleteByApartmentID(apartmentID int64) error {
//     query := `DELETE FROM apartment_images WHERE apartment_id = $1`
//     _, err := r.Db.Exec(query, apartmentID)
//     return err
// }

// func (r *AvatarRepository) DeleteByURL(imageURL string) error {
//     query := `DELETE FROM apartment_images WHERE image_url = $1`
//     _, err := r.Db.Exec(query, imageURL)
//     return err
// }

// func (r *AvatarRepository) UpdatePosition(id int64, position int) error {
//     query := `UPDATE apartment_images SET position = $1 WHERE id = $2`
//     _, err := r.Db.Exec(query, position, id)
//     return err
// }

// func (r *AvatarRepository) GetMainImage(apartmentID int64) (*models.Avatar, error) {
//     query := `
//         SELECT id, apartment_id, image_url, position, created_at
//         FROM apartment_images
//         WHERE apartment_id = $1 AND position = 0
//     `
//     var img models.Avatar
//     err := r.Db.QueryRow(query, apartmentID).Scan(
//         &img.ID,
//         &img.ApartmentID,
//         &img.ImageURL,
//         &img.Position,
//         &img.CreatedAt,
//     )
//     if err == sql.ErrNoRows {
//         return nil, nil
//     }
//     if err != nil {
//         return nil, err
//     }
//     return &img, nil
// }

// func (r *AvatarRepository) CountByApartmentID(apartmentID int64) (int, error) {
//     query := `SELECT COUNT(*) FROM apartment_images WHERE apartment_id = $1`
//     var count int
//     err := r.Db.QueryRow(query, apartmentID).Scan(&count)
//     return count, err
// }