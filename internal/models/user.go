package models

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	IsActive   bool       `json:"is_active"`
    EmailToken *string    `json:"email_token,omitempty"`
}

// type Avatar struct {
//     ID          int64     `json:"id"`
//     UserID      int64     `json:"user_id"`
//     ImageURL    string    `json:"image_url"`
// }
