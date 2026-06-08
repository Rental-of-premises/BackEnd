package models

import (
	"time"
)

type Review struct {
	ID          int64     `json:"id"`
	UserID      *int64    `json:"user_id,omitempty"`
	ApartmentID *int64    `json:"apartment_id,omitempty"`
	Comment     *string   `json:"comment,omitempty"`
	Stars       *int16    `json:"stars,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

