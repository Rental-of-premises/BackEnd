package models

import (
	"time"
)

type Booking struct {
	ID          int64     `json:"id"`
	UserID      *int64    `json:"user_id,omitempty"`
	ApartmentID *int64    `json:"apartment_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	Status      string    `json:"status"`
	TimeFrom    time.Time `json:"time_from"`
	TimeTo      time.Time `json:"time_to"`
}

