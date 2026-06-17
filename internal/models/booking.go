package models

import (
	"time"
)

type Booking struct {
	ID          int64     `json:"id"`
	UserID      int64    `json:"user_id"`
	ApartmentID int64    `json:"apartment_id"`
	CreatedAt   time.Time `json:"created_at"`
	Status      string    `json:"status"`
	TimeFrom    time.Time `json:"time_from"`
	TimeTo      time.Time `json:"time_to"`
}

type BookingFilter struct {
	Status   *string `json:"status"`
	SellerID *int64   `json:"seller_id"`
	ApartmentID *int64   `json:"apartment_id"`
	UserID *int64   `json:"user_id"`
    MinPrice *int64  `json:"min_price"`
    MaxPrice *int64  `json:"max_price"`
    Limit    *int   `json:"limit"`
	Offset   *int   `json:"offset"` 
}